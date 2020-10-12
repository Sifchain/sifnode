// based on draft swagger spec https://raw.githubusercontent.com/Sifchain/sifnode/c1bb5a268da8b519d0fc90f81fa194d31c0f82b3/api/openapi/swagger.yml?token=AAJSXWM6CDXYAEETSC6BJ2S7Q2JLS

import { SlashingParameters, ValidatorSignInfo } from "../entities/Slashing";
import { Transaction } from "../entities/Transaction";

export const slashingService = {
  async getSigningInfos(
    page: number,
    limit: number
  ): Promise<ValidatorSignInfo[]> {
    return [];
  },

  async unjailJailedValidtor(
    validatorAddr: string,
    body: Transaction
  ): Promise<any> {},

  async getParameters(): Promise<SlashingParameters> {},
};
