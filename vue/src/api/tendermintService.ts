// based on draft swagger spec https://raw.githubusercontent.com/Sifchain/sifnode/c1bb5a268da8b519d0fc90f81fa194d31c0f82b3/api/openapi/swagger.yml?token=AAJSXWM6CDXYAEETSC6BJ2S7Q2JLS

import {
  TendermintBlock,
  TendermintNodeInfo,
  TendermintValidatorSet,
} from "../entities/Tendermint";

export const tendermintService = {
  // /node_info
  async getNodeInfo(): Promise<TendermintNodeInfo> {},

  // /syncing
  async getSyncing(): Promise<boolean> {
    return false;
  },

  // GET /blocks/lastest
  async getBlockLatest(): Promise<TendermintBlock> {},

  // GET /blocks/{height}
  async getBlockAtHeight(height: number): Promise<TendermintBlock> {},

  // /validatorsets/latest
  async getValidatorsetLatest(): Promise<TendermintValidatorSet> {},

  // /validatorsets/{height}
  async getValidatorsetAtHeight(
    height: number
  ): Promise<TendermintValidatorSet> {},
};
