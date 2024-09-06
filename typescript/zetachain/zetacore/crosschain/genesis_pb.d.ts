// @generated by protoc-gen-es v1.3.0 with parameter "target=dts"
// @generated from file zetachain/zetacore/crosschain/genesis.proto (package zetachain.zetacore.crosschain, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";
import type { OutboundTracker } from "./outbound_tracker_pb.js";
import type { GasPrice } from "./gas_price_pb.js";
import type { CrossChainTx, ZetaAccounting } from "./cross_chain_tx_pb.js";
import type { LastBlockHeight } from "./last_block_height_pb.js";
import type { InboundHashToCctx } from "./inbound_hash_to_cctx_pb.js";
import type { InboundTracker } from "./inbound_tracker_pb.js";
import type { RateLimiterFlags } from "./rate_limiter_flags_pb.js";

/**
 * GenesisState defines the crosschain module's genesis state.
 *
 * @generated from message zetachain.zetacore.crosschain.GenesisState
 */
export declare class GenesisState extends Message<GenesisState> {
  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.OutboundTracker outboundTrackerList = 2;
   */
  outboundTrackerList: OutboundTracker[];

  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.GasPrice gasPriceList = 5;
   */
  gasPriceList: GasPrice[];

  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.CrossChainTx CrossChainTxs = 7;
   */
  CrossChainTxs: CrossChainTx[];

  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.LastBlockHeight lastBlockHeightList = 8;
   */
  lastBlockHeightList: LastBlockHeight[];

  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.InboundHashToCctx inboundHashToCctxList = 9;
   */
  inboundHashToCctxList: InboundHashToCctx[];

  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.InboundTracker inbound_tracker_list = 11;
   */
  inboundTrackerList: InboundTracker[];

  /**
   * @generated from field: zetachain.zetacore.crosschain.ZetaAccounting zeta_accounting = 12;
   */
  zetaAccounting?: ZetaAccounting;

  /**
   * @generated from field: repeated string FinalizedInbounds = 16;
   */
  FinalizedInbounds: string[];

  /**
   * @generated from field: zetachain.zetacore.crosschain.RateLimiterFlags rate_limiter_flags = 17;
   */
  rateLimiterFlags?: RateLimiterFlags;

  constructor(data?: PartialMessage<GenesisState>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.crosschain.GenesisState";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GenesisState;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GenesisState;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GenesisState;

  static equals(a: GenesisState | PlainMessage<GenesisState> | undefined, b: GenesisState | PlainMessage<GenesisState> | undefined): boolean;
}
