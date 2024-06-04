package evm

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	// maxFeeSafetyFactor is the safety factor for the MaxFeePerGas calculation to ensure
	// that the transaction will remain marketable for six consecutive 100% full blocks
	MaxFeeSafetyFactor = 2
)

// GetLatestBaseFee returns the base fee of the latest block.
func GetLatestBaseFee(client *ethclient.Client) (*big.Int, error) {
	// check if EIP-1559 is supported by the chain
	// Note: it's better to get base fee from the pending block header but `ethclient` can't handle it
	// see: https://github.com/ethereum/go-ethereum/issues/25537
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return header.BaseFee, nil
}

// EstimateGasPriceLondon estimates the gas for a transaction using the London fork rules.
func EstimateGasPriceLondon(client *ethclient.Client, baseFee *big.Int) (*big.Int, *big.Int, *big.Int, error) {
	// define the gas options for both legacy and EIP-1559
	var gasPrice *big.Int
	var maxPriorityFeePerGas *big.Int
	var maxFeePerGas *big.Int

	// check if EIP-1559 is supported by the chain
	isEIP1559Supported := baseFee != nil

	// estimate gas EIP-1559 or legacy if EIP-1559 not supported
	var err error
	if isEIP1559Supported {
		maxPriorityFeePerGas, err = client.SuggestGasTipCap(context.TODO())
		if err != nil {
			return nil, nil, nil, err
		}

		// use a simple heuristic to calculate the recommended MasFeePerGas for a given BaseFee
		// see: https://www.blocknative.com/blog/eip-1559-fees
		maxFeePerGas = new(big.Int).Add(
			maxPriorityFeePerGas,
			new(big.Int).Mul(baseFee, big.NewInt(MaxFeeSafetyFactor)),
		)
	} else {
		gasPrice, err = client.SuggestGasPrice(context.TODO())
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return gasPrice, maxPriorityFeePerGas, maxFeePerGas, nil
}
