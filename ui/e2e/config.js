export const DEX_TARGET = "localhost:5000";

export const KEPLR_CONFIG = {
  id: "dmkamcknogkgcdfhhbddcghachkejeap",
  ver: "0.8.1_0",
  get path() {
    return `./extensions/${this.id}/${this.ver}`;
  },
  options: {
    address: "sif1m625hcmnkc84cgmef6upzzyfu6mxd4jkpnfwwl",
    name: "juniper",
    mnemonic:
      "clump genre baby drum canvas uncover firm liberty verb moment access draft erupt fog alter gadget elder elephant divide biology choice sentence oppose avoid",
  },
};

export const MM_CONFIG = {
  id: "nkbihfbeogaeaoehlefnkodbefgpgknn",
  ver: "9.1.1_0",
  get path() {
    return `./extensions/${this.id}/${this.ver}`;
  },
  network: {
    name: "mm-e2e",
    port: "7545",
    chainId: "1337",
  },
  options: {
    address: "0x627306090abaB3A6e1400e9345bC60c78a8BEf57",
    mnemonic:
      "candy maple cake sugar pudding cream honey rich smooth crumble sweet treat",
    password: "coolguy21",
  },
};