// based on draft swagger spec https://raw.githubusercontent.com/Sifchain/sifnode/c1bb5a268da8b519d0fc90f81fa194d31c0f82b3/api/openapi/swagger.yml?token=AAJSXWM6CDXYAEETSC6BJ2S7Q2JLS

import { AuthAccountInfo } from "../entities/Auth";

export const authService = {
  async getAccount(address: string): Promise<AuthAccountInfo> {},
};
