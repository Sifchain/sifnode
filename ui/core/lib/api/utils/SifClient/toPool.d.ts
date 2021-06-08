import { Asset, Pool } from "../../../entities";
import { RawPool } from "./x/clp";
export declare const toPool: (nativeAsset: Asset) => (poolData: RawPool) => Pool | null;
