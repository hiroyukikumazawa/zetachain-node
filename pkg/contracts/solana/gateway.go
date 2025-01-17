// Package solana privides structures and constants that are used when interacting with the gateway program on Solana chain.
package solana

import (
	"github.com/gagliardetto/solana-go"
	"github.com/pkg/errors"
	idlgateway "github.com/zeta-chain/protocol-contracts-solana/go-idl/generated"
)

const (
	// PDASeed is the seed for the Solana gateway program derived address
	PDASeed = "meta"

	// RentPayerPDASeed is the seed for the Solana gateway program derived address
	RentPayerPDASeed = "rent-payer"

	// AccountsNumberOfDeposit is the number of accounts required for Solana gateway deposit instruction
	// [signer, pda, system_program]
	accountsNumDeposit = 3

	// AccountsNumberOfDeposit is the number of accounts required for Solana gateway deposit spl instruction
	// [signer, pda, whitelist_entry, mint_account, token_program, from, to]
	accountsNumberDepositSPL = 7
)

var (
	// DiscriminatorInitialize returns the discriminator for Solana gateway 'initialize' instruction
	DiscriminatorInitialize = idlgateway.IDLGateway.GetDiscriminator("initialize")

	// DiscriminatorInitializeRentPayer returns the discriminator for Solana gateway 'initialize_rent_payer' instruction
	DiscriminatorInitializeRentPayer = idlgateway.IDLGateway.GetDiscriminator("initialize_rent_payer")

	// DiscriminatorDeposit returns the discriminator for Solana gateway 'deposit' instruction
	DiscriminatorDeposit = idlgateway.IDLGateway.GetDiscriminator("deposit")

	// DiscriminatorDepositSPL returns the discriminator for Solana gateway 'deposit_spl_token' instruction
	DiscriminatorDepositSPL = idlgateway.IDLGateway.GetDiscriminator("deposit_spl_token")

	// DiscriminatorWithdraw returns the discriminator for Solana gateway 'withdraw' instruction
	DiscriminatorWithdraw = idlgateway.IDLGateway.GetDiscriminator("withdraw")

	// DiscriminatorWithdrawSPL returns the discriminator for Solana gateway 'withdraw_spl_token' instruction
	DiscriminatorWithdrawSPL = idlgateway.IDLGateway.GetDiscriminator("withdraw_spl_token")

	// DiscriminatorWhitelist returns the discriminator for Solana gateway 'whitelist_spl_mint' instruction
	DiscriminatorWhitelistSplMint = idlgateway.IDLGateway.GetDiscriminator("whitelist_spl_mint")
)

// ParseGatewayWithPDA parses the gateway id and program derived address from the given string
func ParseGatewayWithPDA(gatewayAddress string) (solana.PublicKey, solana.PublicKey, error) {
	var gatewayID, pda solana.PublicKey

	// decode gateway address
	gatewayID, err := solana.PublicKeyFromBase58(gatewayAddress)
	if err != nil {
		return gatewayID, pda, errors.Wrap(err, "unable to decode address")
	}

	// compute gateway PDA
	seed := []byte(PDASeed)
	pda, _, err = solana.FindProgramAddress([][]byte{seed}, gatewayID)

	return gatewayID, pda, err
}

// ParseRentPayerPDA parses the rent payer program derived address from the given string
func RentPayerPDA(gateway solana.PublicKey) (solana.PublicKey, error) {
	var rentPayerPda solana.PublicKey
	seed := []byte(RentPayerPDASeed)
	rentPayerPda, _, err := solana.FindProgramAddress([][]byte{seed}, gateway)

	return rentPayerPda, err
}
