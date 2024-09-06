package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/pkg/address"
)

const TypeMsgUpdateSystemContract = "update_system_contract"

var _ sdk.Msg = &MsgUpdateSystemContract{}

func NewMsgUpdateSystemContract(creator string, systemContractAddr string) *MsgUpdateSystemContract {
	return &MsgUpdateSystemContract{
		Creator:                  creator,
		NewSystemContractAddress: systemContractAddr,
	}
}

func (msg *MsgUpdateSystemContract) Route() string {
	return RouterKey
}

func (msg *MsgUpdateSystemContract) Type() string {
	return TypeMsgUpdateSystemContract
}

func (msg *MsgUpdateSystemContract) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateSystemContract) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateSystemContract) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	// check if the system contract address is valid
	err = address.ValidateEVMAddress(msg.NewSystemContractAddress)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid system contract address (%s): %s", msg.NewSystemContractAddress, err)
	}

	return nil
}
