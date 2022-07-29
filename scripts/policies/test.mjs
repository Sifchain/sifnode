#!/usr/bin/env zx

$.verbose = false;

const {
  createQueryClient,
  SifSigningStargateClient,
} = require("@sifchain/stargate");
const { DirectSecp256k1HdWallet } = require("@cosmjs/proto-signing");
const { Decimal } = require("@cosmjs/math");

const queryClients = await createQueryClient(process.env.SIFNODE_NODE);

const tokenEntries = await queryClients.tokenRegistry
  .entries({})
  .then((x) => x.registry?.entries);

const rowan = tokenEntries?.find((x) => x.baseDenom === "rowan");

const mnemonic = `${process.env.ADMIN_MNEMONIC}`;
const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic, {
  prefix: "sif",
});
const [sendingAccount] = await wallet.getAccounts();
const receivingAccount = "sif19vprdtfha0xsls0qlwqj2sas32nqqtf4f0ks3m";

const signingClient = await SifSigningStargateClient.connectWithSigner(
  process.env.SIFNODE_NODE,
  wallet
);

const fee = {
  amount: [
    {
      denom: "rowan",
      amount: "100000000000000000", // 0.1 ROWAN
    },
  ],
  gas: "180000000", // 180k
};

const msgSend = {
  typeUrl: "/cosmos.bank.v1beta1.MsgSend",
  value: {
    fromAddress: sendingAccount.address,
    toAddress: receivingAccount,
    amount: [
      {
        denom: rowan.denom,
        amount: Decimal.fromUserInput(
          "1000000000000000000",
          rowan.decimals.toNumber()
        ).toString(),
      },
    ],
  },
};

await signingClient.signAndBroadcast(
  sendingAccount.address,
  [...Array(1000).keys()].map(() => msgSend),
  fee
);
