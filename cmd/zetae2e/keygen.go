package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeta-chain/node/e2e/config"
)

func NewKeygenCmd() *cobra.Command {
	var InitCmd = &cobra.Command{
		Use:   "keygen <file>",
		Short: "generate new keys in a config file",
		RunE:  keygen,
		Args:  cobra.ExactArgs(1),
	}

	return InitCmd
}

func keygen(_ *cobra.Command, args []string) error {
	conf, err := config.ReadConfig(args[0])
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	err = conf.GenerateKeys()
	if err != nil {
		return fmt.Errorf("generate keys: %w", err)
	}

	err = config.WriteConfig(args[0], conf)
	if err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}
