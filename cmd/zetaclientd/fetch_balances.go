package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
)

type Output struct {
	NoOfDaysAgo             string `csv:"NoOfDaysAgo"`
	ZEVMBSCUSDT             string `csv:"ZEVMBSCUSDT"`
	Erc20CustodyUSDTBalance string `csv:"Erc20CustodyUSDTBalance"`
	ZetaBlock               string `csv:"ZetaBlock"`
	BSCblock                string `csv:"BSCblock"`
	SurplusAmount           string `csv:"SurplusAmount"`
}

func CheckBalanceCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "check-balance",
		Short: "Check the balance of an address",
		RunE:  checkBalance,
	}

}

func checkBalance(_ *cobra.Command, args []string) error {

	//bscClient, err := ethclient.Dial("https://bsc-mainnet.g.allthatnode.com/full/evm/ab05ec6995304adebfef618c2ae126fb")
	//if err != nil {
	//	panic(err)
	//}

	zevmClient, err := ethclient.Dial("https://zetachain-mainnet.g.allthatnode.com/archive/evm")
	if err != nil {
		panic(err)
	}

	//currentBlockZeta := int64(3982316)
	//currentBlockBSC := int64(40551771)
	//blocksInAnHourZeta := GetBlocksInAnHour(6)
	//blocksInAnHourBSC := GetBlocksInAnHour(3)
	//daysAgo := int64(0)
	//daysTocheck := int64(170)
	//checkingInterval := int64(1)
	//output := make([]Output, daysTocheck)

	//for i := int64(0); i < daysTocheck; i++ {
	//	fmt.Println("Fetching balances for days ago: ", daysAgo)
	//	//Fetch the total supply of zrc20 token for BSC.USDT
	//	zrc20BSCUSDTAddress := "0x91d4F0D54090Df2D81e834c3c8CE71C6c865e79F"
	//	zrc20BSCUSDT, err := zrc20.NewZRC20(common.HexToAddress(zrc20BSCUSDTAddress), zevmClient)
	//	if err != nil {
	//		panic(err)
	//	}
	//	zetaBlock := GetBlockNumberForDay(daysAgo, blocksInAnHourZeta, currentBlockZeta)
	//	//zetaBlock = zetaBlock.Sub(zetaBlock, big.NewInt(12500))
	//	totalSupply, err := zrc20BSCUSDT.TotalSupply(&bind.CallOpts{BlockNumber: zetaBlock})
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println("Total supply of BSC.USDT: ", totalSupply.String())
	//
	//	// fetch zrc20USDT balance of  substanceX proxy
	//	substanceXProxy := "0x64663c58D42BA8b5Bb79aD924621e5742e2232D8"
	//	balanceX, err := zrc20BSCUSDT.BalanceOf(&bind.CallOpts{BlockNumber: zetaBlock}, common.HexToAddress(substanceXProxy))
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println("SubstanceX proxy balance of BSC.USDT: ", len(balanceX.String()))
	//
	//	// Fetch the USDT balance for the ERC20 custody contract on bsc
	//	bscCustodyContract := "0x00000fF8fA992424957F97688015814e707A0115"
	//	usdtAdress := "0x55d398326f99059fF775485246999027B3197955"
	//	usdt, err := erc20.NewERC20(common.HexToAddress(usdtAdress), bscClient)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	balance, err := usdt.BalanceOf(&bind.CallOpts{BlockNumber: GetBlockNumberForDay(daysAgo, blocksInAnHourBSC, currentBlockBSC)}, common.HexToAddress(bscCustodyContract))
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	fmt.Println("USDT balance of BSC custody contract: ", balance.String())
	//
	//	surplusAmount := new(big.Int).Sub(balance, totalSupply)
	//	fmt.Printf("Surplus amount: %s\n", new(big.Int).Div(surplusAmount, big.NewInt(1e18)).String())
	//	fmt.Println("------------------------------------------------------")
	//	time.Sleep(2 * time.Second)
	//	output[i] = Output{
	//		NoOfDaysAgo:             strconv.FormatInt(daysAgo, 10),
	//		ZEVMBSCUSDT:             new(big.Int).Div(totalSupply, big.NewInt(1e18)).String(),
	//		Erc20CustodyUSDTBalance: new(big.Int).Div(balance, big.NewInt(1e18)).String(),
	//		ZetaBlock:               strconv.FormatInt(zetaBlock.Int64(), 10),
	//		BSCblock:                strconv.FormatInt(GetBlockNumberForDay(daysAgo, blocksInAnHourBSC, currentBlockBSC).Int64(), 10),
	//		SurplusAmount:           new(big.Int).Div(surplusAmount, big.NewInt(1e18)).String(),
	//	}
	//	daysAgo = daysAgo + checkingInterval
	//}
	//
	//file, err := os.Create("output.csv")
	//if err != nil {
	//	log.Fatalf("Failed to open file: %v", err)
	//}
	//defer file.Close()
	//
	//writer := csv.NewWriter(file)
	//defer writer.Flush()
	//
	//if err := writer.Write([]string{"NoOfDaysAgo", "Erc20CustodyUSDTBalance", "ZEVMBSCUSDT", "ZetaBlock", "BscBlock", "SurplusAmount"}); err != nil {
	//	log.Fatalf("Cannot write header: %v", err)
	//}
	//
	//for _, record := range output {
	//	if err := writer.Write([]string{record.NoOfDaysAgo, record.Erc20CustodyUSDTBalance, record.ZEVMBSCUSDT, record.ZetaBlock, record.BSCblock, record.SurplusAmount}); err != nil {
	//		log.Fatalf("Cannot write record: %v", err)
	//	}
	//}
	startingblock := big.NewInt(3762177)

	for i := int64(0); i < 4; i++ {
		zrc20BSCUSDTAddress := "0x91d4F0D54090Df2D81e834c3c8CE71C6c865e79F"
		zrc20BSCUSDT, err := zrc20.NewZRC20(common.HexToAddress(zrc20BSCUSDTAddress), zevmClient)
		if err != nil {
			panic(err)
		}
		//zetaBlock = zetaBlock.Sub(zetaBlock, big.NewInt(12500))
		totalSupply, err := zrc20BSCUSDT.TotalSupply(&bind.CallOpts{BlockNumber: startingblock})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Total supply of BSC.USDT: %s , block : %s \n", new(big.Int).Div(totalSupply, big.NewInt(1e18)).String(), startingblock.String())
		startingblock = startingblock.Add(startingblock, big.NewInt(1))
		time.Sleep(1 * time.Second)
	}

	return nil
}

func GetBlocksInAnHour(blockTime int64) int64 {
	return 60 / blockTime
}

func GetBlockNumberForDay(daysAgo int64, blocksInAnHour int64, currentBlock int64) *big.Int {
	blocksTotal := blocksInAnHour * 60 * 24 * daysAgo
	inBlock := currentBlock - blocksTotal
	return big.NewInt(inBlock)
}
