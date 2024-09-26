package app

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"testing"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/store"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	simcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/ethermint/app"
	evmante "github.com/zeta-chain/ethermint/app/ante"
	"github.com/zeta-chain/node/app/ante"
)

const (
	SimAppChainID  = "simulation_777-1"
	SimBlockMaxGas = 815000000
	TestAppChainID = "zetachain_777-1"
)

// NewSimApp disable feemarket on native tx, otherwise the cosmos-sdk simulation tests will fail.
func NewSimApp(logger log.Logger, db dbm.DB, baseAppOptions ...func(*baseapp.BaseApp)) (*App, error) {
	encCdc := MakeEncodingConfig()
	app := New(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		app.DefaultNodeHome,
		5,
		encCdc,
		simtestutil.EmptyAppOptions{},
	)
	// disable feemarket on native tx
	options := ante.HandlerOptions{
		AccountKeeper:   app.AccountKeeper,
		BankKeeper:      app.BankKeeper,
		EvmKeeper:       app.EvmKeeper,
		FeeMarketKeeper: app.FeeMarketKeeper,
		SignModeHandler: encCdc.TxConfig.SignModeHandler(),
		SigGasConsumer:  evmante.DefaultSigVerificationGasConsumer,
		MaxTxGasWanted:  0,
		ObserverKeeper:  app.ObserverKeeper,
	}

	anteHandler, err := ante.NewAnteHandler(options)
	if err != nil {
		panic(err)
	}

	app.SetAnteHandler(anteHandler)
	if err := app.LoadLatestVersion(); err != nil {
		return nil, err
	}
	return app, nil
}

// interBlockCacheOpt returns a BaseApp option function that sets the persistent
// inter-block write-through cache.
func interBlockCacheOpt() func(*baseapp.BaseApp) {
	return baseapp.SetInterBlockCache(store.NewCommitKVStoreCacheManager())
}

func TestAppStateDeterminism(t *testing.T) {
	if !simcli.FlagEnabledValue {
		t.Skip("skipping application simulation")
	}

	config := simcli.NewConfigFromFlags()
	config.InitialBlockHeight = 1
	config.ExportParamsPath = ""
	config.OnOperation = false
	config.AllInvariants = false
	config.ChainID = SimAppChainID
	config.BlockMaxGas = SimBlockMaxGas

	numSeeds := 3
	numTimesToRunPerSeed := 5
	appHashList := make([]json.RawMessage, numTimesToRunPerSeed)

	for i := 0; i < numSeeds; i++ {
		config.Seed = rand.Int63()

		for j := 0; j < numTimesToRunPerSeed; j++ {
			var logger log.Logger
			if simcli.FlagVerboseValue {
				logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
			} else {
				logger = log.NewNopLogger()
			}

			db := dbm.NewMemDB()
			app, err := NewSimApp(logger, db, interBlockCacheOpt(), baseapp.SetChainID(SimAppChainID))
			require.NoError(t, err)

			fmt.Printf(
				"running non-determinism simulation; seed %d: %d/%d, attempt: %d/%d\n",
				config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
			)

			_, _, err = simulation.SimulateFromSeed(
				t,
				os.Stdout,
				app.BaseApp,
				StateFn(app),
				simtypes.RandomAccounts, // Replace with own random account function if using keys other than secp256k1
				simtestutil.SimulationOperations(app, app.AppCodec(), config),
				app.ModuleAccountAddrs(),
				config,
				app.AppCodec(),
			)
			require.NoError(t, err)

			if config.Commit {
				simtestutil.PrintStats(db)
			}

			appHash := app.LastCommitID().Hash
			appHashList[j] = appHash

			if j != 0 {
				require.Equal(
					t, string(appHashList[0]), string(appHashList[j]),
					"non-determinism in seed %d: %d/%d, attempt: %d/%d\n", config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
				)
			}
		}
	}
}
