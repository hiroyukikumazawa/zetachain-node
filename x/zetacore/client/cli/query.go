package cli

import (
	"fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group zetacore queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdListZetaConversionRate())
	cmd.AddCommand(CmdShowZetaConversionRate())
	// this line is used by starport scaffolding # 1

	cmd.AddCommand(CmdShowKeygen())

	cmd.AddCommand(CmdListTSSVoter())
	cmd.AddCommand(CmdShowTSSVoter())

	cmd.AddCommand(CmdListTSS())
	cmd.AddCommand(CmdShowTSS())

	cmd.AddCommand(CmdTxListRich())

	cmd.AddCommand(CmdListInTx())
	cmd.AddCommand(CmdShowInTx())

	cmd.AddCommand(CmdShowTxList())

	cmd.AddCommand(CmdListGasBalance())
	cmd.AddCommand(CmdShowGasBalance())

	cmd.AddCommand(CmdListGasPrice())
	cmd.AddCommand(CmdShowGasPrice())

	cmd.AddCommand(CmdListChainNonces())
	cmd.AddCommand(CmdShowChainNonces())

	cmd.AddCommand(CmdListLastBlockHeight())
	cmd.AddCommand(CmdShowLastBlockHeight())

	cmd.AddCommand(CmdListReceive())
	cmd.AddCommand(CmdShowReceive())

	cmd.AddCommand(CmdListSend())
	cmd.AddCommand(CmdShowSend())

	cmd.AddCommand(CmdListNodeAccount())
	cmd.AddCommand(CmdShowNodeAccount())

	cmd.AddCommand(CmdLastMetaHeight())

	return cmd
}
