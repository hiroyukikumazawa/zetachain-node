package main

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/app"
	zetae2econfig "github.com/zeta-chain/node/cmd/zetae2e/config"
	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/txserver"
	"github.com/zeta-chain/node/e2e/utils"
)

// NewSetupEVMCmd sets up an EVM chain for
func NewSetupEVMCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup-evm <config>",
		Short: "Setup an external evm chain",
		RunE:  runSetupEVM,
		Args:  cobra.ExactArgs(1),
	}
	return cmd
}

func runSetupEVM(_ *cobra.Command, args []string) error {
	// read the config file
	conf, err := config.ReadConfig(args[0])
	if err != nil {
		return err
	}
	logger := runner.NewLogger(false, color.FgHiYellow, "")
	app.SetConfig()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	zetaTxServer, err := txserver.NewZetaTxServer(
		conf.RPCs.ZetaCoreRPC,
		[]string{utils.AdminPolicyName},
		[]string{
			conf.PolicyAccounts.AdminPolicyAccount.RawPrivateKey.String(),
		},
		conf.ZetaChainID,
	)
	if err != nil {
		return fmt.Errorf("new zeta tx server: %w", err)
	}

	// initialize deployer runner with config
	r, err := zetae2econfig.RunnerFromConfig(
		ctx,
		"e2e",
		cancel,
		conf,
		conf.DefaultAccount,
		logger,
		runner.WithZetaTxServer(zetaTxServer),
	)
	if err != nil {
		return err
	}

	r.SetupEVM(false, false)
	r.SetupEVMV2()

	logger.Print("* EVM setup done")

	return nil
}
