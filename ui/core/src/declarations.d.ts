declare module "@metamask/detect-provider" {
  import { MetaMaskInpageProvider } from "@metamask/inpage-provider";
  export = function detectProvider(): Promise<MetaMaskInpageProvider | null> {};
}

declare module "toformat";

declare module "ganache-time-traveler" {
  export function advanceTime(time: any): Promise<any>;
  export function advanceBlock(): Promse<any>;
  export function advanceBlockAndSetTime(): Promse<any>;
  export function advanceTimeAndBlock(time: any): Promse<any>;
  export function takeSnapshot(): Promse<number>;
  export function revertToSnapshot(id: number): Promse<any>;
}
