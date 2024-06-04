package evm_test

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/zetaclient/chains/evm"

	"testing"
)

var (
	ENV_TEST_PK_EVM   = "TEST_PK_EVM"
	URLEthMainnet     = "https://rpc.ankr.com/eth"
	URLEthSepolia     = "https://rpc.ankr.com/eth_sepolia"
	URLBscMainnet     = "https://rpc.ankr.com/bsc"
	URLPolygonMainnet = "https://rpc.ankr.com/polygon"
)

// LoadTestPrivateKeyHex loads a test private key from an environment variable
func LoadTestPrivateKeyHex(t *testing.T) (*ecdsa.PrivateKey, ethcommon.Address) {
	// get private key from environment
	testPrivKeyStr := os.Getenv("TEST_PK_EVM")
	testPrivKeyECDSA, err := crypto.HexToECDSA(testPrivKeyStr)
	require.NoError(t, err)

	// derive 0x address
	testPubKey := testPrivKeyECDSA.Public()
	testPubKeyECDSA, ok := testPubKey.(*ecdsa.PublicKey)
	require.True(t, ok)
	address := crypto.PubkeyToAddress(*testPubKeyECDSA)
	fmt.Printf("test address: %s\n", address.Hex())

	return testPrivKeyECDSA, address
}

func TestRPCLive(t *testing.T) {
	//LiveTestEstimateGasPriceLondon(t)
	//LiveTestEIP1559TxFee(t)
}

func LiveTestEstimateGasPriceLondon(t *testing.T) {
	evmClient, err := ethclient.Dial(URLEthMainnet)
	if err != nil {
		t.Error(err)
	}

	// get latest header
	header, err := evmClient.HeaderByNumber(context.TODO(), nil)
	require.NoError(t, err)

	// estimate gas EIP1559
	gasPrice, maxPriorityFeePerGas, maxFeePerGas, err := evm.EstimateGasPriceLondon(evmClient, header.BaseFee)
	require.NoError(t, err)

	fmt.Printf("GasPrice: %s\n", gasPrice)
	fmt.Printf("BaseFee: %s\n", header.BaseFee)
	fmt.Printf("MaxPriorityFeePerGas: %s\n", maxPriorityFeePerGas)
	fmt.Printf("MaxFeePerGas: %s\n", maxFeePerGas)

	// check result
	if header.BaseFee == nil {
		require.NotNil(t, gasPrice)
		require.Nil(t, maxPriorityFeePerGas)
		require.Nil(t, maxFeePerGas)
	} else {
		require.Nil(t, gasPrice)
		require.NotNil(t, maxPriorityFeePerGas)
		require.NotNil(t, maxFeePerGas)

		// check if maxFeePerGas is calculated correctly
		expectedMaxFeePerGas := new(big.Int).Add(
			maxPriorityFeePerGas,
			new(big.Int).Mul(header.BaseFee, big.NewInt(evm.MaxFeeSafetyFactor)),
		)
		require.Equal(t, expectedMaxFeePerGas, maxFeePerGas)
	}
}

func LiveTestEIP1559TxFee(t *testing.T) {
	evmClient, err := ethclient.Dial(URLEthMainnet)
	if err != nil {
		t.Error(err)
	}

	// given a specific EIP1559 transaction
	// https://etherscan.io/tx/0xe70f28b7ca19a5f844dff0b1637a1d3827ce1b353e8735a8cacadf766468ef3f
	blockNumber := uint64(19987674)
	hash := ethcommon.HexToHash("0xe70f28b7ca19a5f844dff0b1637a1d3827ce1b353e8735a8cacadf766468ef3f")

	// get base fee for the block
	header, err := evmClient.HeaderByNumber(context.TODO(), new(big.Int).SetUint64(blockNumber))
	require.NoError(t, err)
	baseFee := header.BaseFee

	// get EIP-1559 transaction
	tx, isPending, err := evmClient.TransactionByHash(context.Background(), hash)
	require.NoError(t, err)
	require.False(t, isPending)

	// get EIP-1559 transaction receipt
	receipt, err := evmClient.TransactionReceipt(context.Background(), hash)
	require.NoError(t, err)
	require.EqualValues(t, tx.Type(), ethtypes.DynamicFeeTxType)

	// print transaction details
	fmt.Printf("BaseFee: %s\n", baseFee)
	fmt.Printf("GasTipCap: %s\n", tx.GasTipCap())
	fmt.Printf("GasFeeCap: %s\n", tx.GasFeeCap())

	// expected transaction fee
	// As of EIP-1559, the overall fee a transaction creator pays is calculated as: ( (base fee + priority fee) x units of gas used).
	gasPriceToPay := new(big.Int).Add(baseFee, tx.GasTipCap())
	gasFeeToPay := new(big.Int).Mul(gasPriceToPay, new(big.Int).SetUint64(receipt.GasUsed))
	gasFeeCap := new(big.Int).Mul(tx.GasFeeCap(), new(big.Int).SetUint64(receipt.GasUsed))

	// check actual gas fee (in etherscan) against expected numbers
	etherScanPaidFee := new(big.Int).SetUint64(115709596800000)
	require.True(t, tx.GasPrice().Cmp(tx.GasFeeCap()) == 0)
	require.True(t, gasFeeCap.Cmp(gasFeeToPay) >= 0)
	require.Equal(t, etherScanPaidFee, gasFeeToPay)
}

/*

Here is a comparision of a legacy transaction and an EIP-1559 transaction:


Transaction fee:

A legacy transaction fee is calculated as: `gasPrice` x `gasUsed`

An EIP-1559 transaction fee is calculated as: (`baseFee` + `maxPriorityFeePerGas`) x `gasUsed`, see: https://support.metamask.io/transactions-and-gas/gas-fees/how-to-estimate-the-gas-fee/#:~:text=As%20of%20EIP%2D1559%2C%20the,x%20units%20of%20gas%20used).


From the perspective of CCTX fee charging: `gasPrice` === `baseFee + maxPriorityFeePerGas`


To keep the CCTX fee charging model unchanged, zetaclient can simply post (`baseFee + maxPriorityFeePerGas`) as a replacement of `gasPrice` for EIP-1559 enabled chains.


Transaction building:
A legacy transaction requires: `gasPrice`, `gasLimit`
An EIP-1559 transaction requires: `maxPriorityFeePerGas`, `maxFeePerGas`, `gasLimit`

Sending a EIP-1559 transaction without explicit `maxPriorityFeePerGas` or `maxFeePerGas` will result in error (in live tests).

Where to get deterministic numbers of `maxPriorityFeePerGas` and `maxFeePerGas` when signing a TSS transaction?
Remember the equation: `gasPrice` === `baseFee + maxPriorityFeePerGas`
If `maxPriorityFeePerGas` becomes a field in CCTX outbound parameter, the zetaclient can simply calculate `baseFee` as `gasPrice - maxPriorityFeePerGas`.

The `maxFeePerGas` is the upper limit of gas price that the sender is willing to pay. The estimation of `maxFeePerGas` can be done by:
`maxFeePerGas = maxPriorityFeePerGas + baseFee * 2`
See this article for more details of gas fee estimation: https://www.blocknative.com/blog/eip-1559-fees
Also see ChainSafe's gas estimator for its bridge: https://github.com/ChainSafe/chainbridge-core/blob/main/chains/evm/calls/evmgaspricer/london.go#L65


The current model of a CCTX outbound parameter only contains `gasPrice` and `gasLimit`. To support EIP-1559, a new field `maxPriorityFeePerGas` needs to be added.



Transaction replacement:

To replace an EVM-chain pending outbound transaction stuck in mempool:

A legacy transaction can be replaced by a new transaction with
- same `nonce`
- a higher `gasPrice`


An EIP-1559 transaction can be replaced by a new transaction with
- same `nonce`
- a higher `maxPriorityFeePerGas`, at least 10% higher
- a higher `maxFeePerGas`, at least 10% higher

See Alchemy doc: https://docs.alchemy.com/docs/retrying-an-eip-1559-transaction
See go-ethereum: https://github.com/ethereum/go-ethereum/blob/87246f3cbaf10f83f56bc4d45f0f3e36e83e71e9/core/txpool/legacypool/list.go#L323
and https://github.com/ethereum/go-ethereum/blob/87246f3cbaf10f83f56bc4d45f0f3e36e83e71e9/core/txpool/legacypool/legacypool.go#L146


When the gas stability pool kicks in, it has to increase both `gasPrice` and `maxPriorityFeePerGas` by some percentage.
Again, remember the equation: `gasPrice` === `baseFee + maxPriorityFeePerGas`
The new `baseFee` will be calculated as `gasPrice - maxPriorityFeePerGas`, and the new `maxFeePerGas` will be calculated as `maxPriorityFeePerGas + baseFee * 2`.

After calculation, zetaclient will be able to broadcast updated transaction with new `maxPriorityFeePerGas` and `maxFeePerGas`.

*/

func TestEIP1559TxOnlyMaxFeePerGas(t *testing.T) {
	// create Sepolia client
	client, err := ethclient.Dial("https://rpc.ankr.com/eth_sepolia")
	//client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/m79YhhEuwED9ZOnpS1qFWKN3i2l9tv-d")
	if err != nil {
		t.Error(err)
	}

	histry, err := client.FeeHistory(context.Background(), 10, big.NewInt(6033812), nil)
	require.NoError(t, err)
	fmt.Printf("history: %v\n", histry)

	// sender and receiver addresses
	senderPrivateKey, senderAddress := LoadTestPrivateKeyHex(t)
	receiverAddress := ethcommon.HexToAddress("0x671Fb64365c7656C0D955aDcBcae8e3F62fF6A1B")

	// get chain ID
	chainID, err := client.ChainID(context.TODO())
	require.NoError(t, err)

	// get pending nonce
	nonce, err := client.PendingNonceAt(context.Background(), senderAddress)
	require.NoError(t, err)

	// get base fee from latest header
	header, err := client.HeaderByNumber(context.TODO(), nil)
	require.NoError(t, err)

	// estimate gas EIP-1559
	gasPrice, maxPriorityFeePerGas, maxFeePerGas, err := evm.EstimateGasPriceLondon(client, header.BaseFee)
	require.NoError(t, err)
	require.Nil(t, gasPrice)

	// print gas prices
	fmt.Printf("BaseFee  : %s\n", header.BaseFee)
	fmt.Printf("GasTipCap: %s\n", maxPriorityFeePerGas)
	fmt.Printf("GasFeeCap: %s\n", maxFeePerGas)

	// create transaction
	rawTx := ethtypes.NewTx(&ethtypes.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		To:        &receiverAddress,
		Value:     big.NewInt(12e13), // 0.00012 ETH
		GasTipCap: maxPriorityFeePerGas,
		GasFeeCap: maxFeePerGas,
		Gas:       21000,
		Data:      nil,
	})

	// sign transaction
	signedTx, err := ethtypes.SignTx(rawTx, ethtypes.LatestSignerForChainID(chainID), senderPrivateKey)
	require.NoError(t, err)
	fmt.Printf("default maxPriorityFeePerGas: %s\n", signedTx.GasTipCap())

	// send transaction
	err = client.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)

	// print transaction hash
	fmt.Printf("transaction sent: %s\n", signedTx.Hash().Hex())

	// wait for receipt
	receipt := utils.MustWaitForTxReceipt(context.Background(), client, signedTx, nil, 30*time.Second)
	if receipt.Status == ethtypes.ReceiptStatusFailed {
		t.Error("deposit failed")
	}

	// retrieve EIP-1559 transaction
	tx, isPending, err := client.TransactionByHash(context.Background(), signedTx.Hash())
	require.NoError(t, err)
	require.False(t, isPending)

	// get base fee of included block
	header, err = client.HeaderByNumber(context.TODO(), receipt.BlockNumber)
	require.NoError(t, err)
	fmt.Printf("BaseFee: %s\n", header.BaseFee)

	// print transaction details
	gasTipCap := tx.GasTipCap()
	fmt.Printf("paid gasTipCap: %s\n", gasTipCap)
	fmt.Printf("paid gasPrice: %s\n", tx.GasPrice())

	// expected transaction fee: ((base fee + priority fee) x gas used)
	gasPriceToPay := new(big.Int).Add(header.BaseFee, gasTipCap)
	gasFeeToPay := new(big.Int).Mul(gasPriceToPay, new(big.Int).SetUint64(receipt.GasUsed))
	fmt.Printf("expected gas fee: %s\n", gasFeeToPay)

	// check actual gas fee (in etherscan) against expected numbers (eye-balling)
}
