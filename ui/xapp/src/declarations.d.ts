declare module "@metamask/detect-provider" {
  import { MetaMaskInpageProvider } from "@metamask/inpage-provider";
  export = function detectProvider(): Promise<MetaMaskInpageProvider | null> {};
}

declare module "toformat";
