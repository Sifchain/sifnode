export class Asset {
  constructor(
    public decimals: number,
    public symbol: string,
    public name: string
  ) {}
  static create(decimals: number, symbol: string, name: string): Asset {
    return new Asset(decimals, symbol, name);
  }
}

export const createAsset = Asset.create;
