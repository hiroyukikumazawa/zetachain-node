package app

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	evmenc "github.com/evmos/ethermint/encoding"
)

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() EncodingConfig {
	//encodingConfig := params.MakeEncodingConfig()
	encodingConfig := evmenc.MakeConfig(ModuleBasics)
	//std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	//std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	//ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	//ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	//return encodingConfig
	return EncodingConfig{
		InterfaceRegistry: encodingConfig.InterfaceRegistry,
		Codec:             encodingConfig.Codec,
		TxConfig:          encodingConfig.TxConfig,
		Amino:             encodingConfig.Amino,
	}
}

type EncodingConfig struct {
	InterfaceRegistry types.InterfaceRegistry
	Codec             codec.Codec
	TxConfig          client.TxConfig
	Amino             *codec.LegacyAmino
}
