declare module "@metamask/detect-provider" {
  type DetectProviderReturn = Promise<MetaMaskInpageProvider | null>;
  import { MetaMaskInpageProvider } from "@metamask/inpage-provider";
  function detectProvider(): DetectProviderReturn;
  export = detectProvider;
}

declare module "toformat";

declare module "ganache-time-traveler" {
  export function advanceTime(time: unknown): Promise<unknown>;
  export function advanceBlock(): Promse<unknown>;
  export function advanceBlockAndSetTime(): Promse<unknown>;
  export function advanceTimeAndBlock(time: unknown): Promse<unknown>;
  export function takeSnapshot(): Promse<number>;
  export function revertToSnapshot(id: number): Promse<unknown>;
}
