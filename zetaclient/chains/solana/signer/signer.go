package signer

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	clientcontext "github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/outboundprocessor"
)

// Signer deals with signing BTC transactions and implements the ChainSigner interface
type Signer struct {
	*base.Signer

	// client is the RPC client to interact with the Bitcoin chain
	client *rpc.Client
}

// NewSigner creates a new Bitcoin signer
func NewSigner(
	chain chains.Chain,
	appContext *clientcontext.AppContext,
	tss interfaces.TSSSigner,
	ts *metrics.TelemetryServer,
	logger base.Logger,
	// client *rpc.Client,
	// cfg config.BTCConfig
) (*Signer, error) {
	// create base signer
	baseSigner := base.NewSigner(chain, appContext, tss, ts, logger)
	// FIXME: config RPC
	client := rpc.New("http://solana:8899")

	return &Signer{
		Signer: baseSigner,
		client: client,
	}, nil
}

var _ interfaces.ChainSigner = &Signer{}

func (s *Signer) TryProcessOutbound(cctx *types.CrossChainTx, outboundProc *outboundprocessor.Processor, outboundID string, observer interfaces.ChainObserver, zetacoreClient interfaces.ZetacoreClient, height uint64) {
	defer func() {
		outboundProc.EndTryProcess(outboundID)
		if err := recover(); err != nil {
			s.Logger().Std.Error().Msgf("Solana TryProcessOutbound: %s, caught panic error: %v", cctx.Index, err)
		}
	}()

	logger := s.Logger().Std.With().
		Str("OutboundID", outboundID).
		Str("SendHash", cctx.Index).
		Logger()

	params := cctx.GetCurrentOutboundParam()
	coinType := cctx.InboundParams.CoinType
	if coinType == coin.CoinType_Zeta || coinType == coin.CoinType_ERC20 {
		logger.Error().Msgf("Solana TryProcessOutbound: can only send SOL to a Solana network")
		return
	}
	logger.Info().Msgf("Solana TryProcessOutbound: %s, value %d to %s", cctx.Index, params.Amount.BigInt(), params.Receiver)

	outboundTssNonce := params.TssNonce
	// get size limit and gas price
	//fee := 5000 // FIXME: this is the fixed fee (for signatures), explore priority fee for compute units

	//to, err := chains.DecodeBtcAddress(params.Receiver, params.ReceiverChainId)
	// NOTE: withrawal event hook must validate the receiver address format
	to := solana.MustPublicKeyFromBase58(params.Receiver)
	amount := params.Amount.Uint64()

	{ // TODO: refactor this piece out to a separate (withdraw) function
		// FIXME: config this; right now it's the same privkey used by local e2e test_solana_*.go
		privkey := solana.MustPrivateKeyFromBase58("4yqSQxDeTBvn86BuxcN5jmZb2gaobFXrBqu8kiE9rZxNkVMe3LfXmFigRsU4sRp7vk4vVP1ZCFiejDKiXBNWvs2C")
		type WithdrawInstructionParams struct {
			Discriminator [8]byte
			Amount        uint64
			Signature     [64]byte
			RecoveryID    uint8
			MessageHash   [32]byte
			Nonce         uint64
		}
		chain, chainParams, ok := s.AppContext().GetSolanaChainParams()
		if !ok {
			s.Logger().Std.Error().Msg("cannot get chain params")
			return
		}
		programId := solana.MustPublicKeyFromBase58(chainParams.GatewayAddress)
		seed := []byte("meta")
		pdaComputed, bump, err := solana.FindProgramAddress([][]byte{seed}, programId)
		if err != nil {
			panic(err)
		}
		fmt.Printf("computed pda: %s, bump %d\n", pdaComputed, bump)
		type PdaInfo struct {
			Discriminator [8]byte
			Nonce         uint64
			TssAddress    [20]byte
			Authority     [32]byte
			ChainID       uint64
		}
		pdaInfo, err := s.client.GetAccountInfo(context.TODO(), pdaComputed)
		if err != nil {
			panic(err)
		}
		fmt.Printf("pdainfo: %v\n", pdaInfo.Bytes())
		var pda PdaInfo
		err = borsh.Deserialize(&pda, pdaInfo.Bytes())
		if err != nil {
			panic(err)
		}
		fmt.Printf("pda parsed: %+v\n", pda)

		recent, err := s.client.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
		if err != nil {
			panic(err)
		}
		fmt.Println("recent blockhash:", recent.Value.Blockhash)
		var inst solana.GenericInstruction

		var message []byte
		bytes := make([]byte, 8)
		chainId := uint64(chain.ChainId)
		nonce := outboundTssNonce
		binary.BigEndian.PutUint64(bytes, chainId)
		message = append(message, bytes...)
		binary.BigEndian.PutUint64(bytes, nonce)
		message = append(message, bytes...)
		binary.BigEndian.PutUint64(bytes, amount)
		message = append(message, bytes...)
		message = append(message, to.Bytes()...)
		messageHash := crypto.Keccak256Hash(message)
		fmt.Printf("solana msghash: chainid %d, nonce %d, amount %d, to %s, hash %s", chainId, nonce, amount, to.String(), messageHash.String())
		// this sig will be 65 bytes; R || S || V, where V is 0 or 1
		signature, err := s.TSS().Sign(messageHash.Bytes(), height, nonce, chain.ChainId, "")
		if err != nil {
			s.Logger().Std.Error().Err(err).Msg("cannot sign message")
			panic(err)
		}
		s.Logger().Std.Info().
			Msgf("Key-sign success: %d => %s, nonce %d", cctx.InboundParams.SenderChainId, chain.ChainName, outboundTssNonce)

		s.Logger().Std.Info().Msgf("recovery id %d", signature[64])
		var sig [64]byte
		copy(sig[:], signature[:64])

		inst.DataBytes, err = borsh.Serialize(WithdrawInstructionParams{
			Discriminator: [8]byte{183, 18, 70, 156, 148, 109, 161, 34},
			Amount:        amount,
			Signature:     sig,
			RecoveryID:    signature[64],
			MessageHash:   messageHash,
			Nonce:         nonce,
		})
		var accountSlice []*solana.AccountMeta
		accountSlice = append(accountSlice, solana.Meta(privkey.PublicKey()).WRITE().SIGNER())
		accountSlice = append(accountSlice, solana.Meta(pdaComputed).WRITE())
		accountSlice = append(accountSlice, solana.Meta(to).WRITE())
		accountSlice = append(accountSlice, solana.Meta(programId))
		inst.ProgID = programId
		inst.AccountValues = accountSlice
		tx, err := solana.NewTransaction(
			[]solana.Instruction{&inst},
			recent.Value.Blockhash,
			solana.TransactionPayer(privkey.PublicKey()),
		)
		if err != nil {
			panic(err)
		}
		_, err = tx.Sign(
			func(key solana.PublicKey) *solana.PrivateKey {
				if privkey.PublicKey().Equals(key) {
					return &privkey
				}
				return nil
			},
		)
		if err != nil {
			panic(fmt.Errorf("unable to sign transaction: %w", err))
		}
		spew.Dump(tx)
		// FIXME: simulate before broadcast!
		txsig, err := s.client.SendTransactionWithOpts(
			context.TODO(),
			tx,
			rpc.TransactionOpts{
				// default PreflightCommitment is "finalized" which is too conservative
				// and results in repeated broadcast attempts that only 1 will succeed
				// Setting a "processed" level will simulate tx against more recent state
				// thus fails faster after a tx is already broadcasted and processed in a block.
				// This reduces the number of "failed" txs due to repeated broadcast attempts.
				PreflightCommitment: rpc.CommitmentProcessed,
			},
		)
		if err != nil {
			s.Logger().Std.Warn().Err(err).Msg("broadcast error")
		} else {
			s.Logger().Std.Info().Msgf("broadcast success! tx sig %s; waiting for confirmation...", txsig)
			// launch a go routine with timeout to check for tx confirmation;
			// repeatedly query until timeout or the transaction is included in a block, either with success or failure
			go func() {
				txsig := txsig // capture the value
				nonce := nonce
				t1 := time.Now()
				for {
					if time.Since(t1) > 2*time.Minute {
						return
					}
					out, err := s.client.GetConfirmedTransactionWithOpts(context.TODO(), txsig, &rpc.GetTransactionOpts{
						// I'd like to use "CommitmentProcessed" but it seems not supported in RPC: see https://solana.com/docs/rpc/http/gettransaction
						Commitment: rpc.CommitmentConfirmed,
					})
					if err == nil {
						if out.Meta.Err == nil { // successfully included in a block; report and exit goroutine
							txhash, err := zetacoreClient.AddOutboundTracker(s.Chain().ChainId, nonce, txsig.String(), nil, "", -1)
							if err != nil {
								s.Logger().Std.Error().Err(err).Msgf("unable to add to tracker: tx %s", txsig)
							} else {
								s.Logger().Std.Info().Msgf("added txsig %s to outbound tracker; zeta txhash %s", txsig, txhash)
							}
							return
						} else { // it's included by failed (likely competing txs succeeded). exit goroutine.
							s.Logger().Std.Warn().Msgf("tx %s failed: %v", txsig, out.Meta.Err)
							return
						}
					}
					time.Sleep(10 * time.Second)
				}
			}()
		}
	}
}

func (s *Signer) SetZetaConnectorAddress(address ethcommon.Address) {
	panic("implement me")
}

func (s *Signer) SetERC20CustodyAddress(address ethcommon.Address) {
	panic("SetERC20CustodyAddress should not be called on Solana signer")
}

func (s *Signer) GetZetaConnectorAddress() ethcommon.Address {
	panic("GetZetaConnectorAddress should not be called on Solana signer")
}

func (s *Signer) GetERC20CustodyAddress() ethcommon.Address {
	panic("GetERC20CustodyAddress should not be called on Solana signer")
}
