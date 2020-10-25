import { Asset } from "../Asset";
import { Token } from "../Token";

export function isToken(value?: Asset): value is Token {
  return value ? Object.keys(value).includes("address") : false;
}
