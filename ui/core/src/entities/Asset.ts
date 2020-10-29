import { ChainId } from "./ChainId";

export class Asset {
  constructor(
    public decimals: number,
    public symbol: string,
    public name: string,
    public chainId: ChainId
  ) {}
  static create(
    symbol: string,
    decimals: number,
    name: string,
    chainId: ChainId
  ): Asset {
    return new Asset(decimals, symbol, name, chainId);
  }
}

export const createAsset = Asset.create;
