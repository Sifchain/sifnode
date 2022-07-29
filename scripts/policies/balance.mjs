#!/usr/bin/env zx

$.verbose = false;

const {
  createQueryClient,
  SifSigningStargateClient,
} = require("@sifchain/stargate");
const { DirectSecp256k1HdWallet } = require("@cosmjs/proto-signing");
const { Decimal } = require("@cosmjs/math");

const queryClients = await createQueryClient(process.env.SIFNODE_NODE);

// const response = await queryClients.clp.getPools({});

// const response = await queryClients.bank.balance(
//   process.env.ADMIN_ADDRESS,
//   "rowan"
// );

const response = await queryClients.bank.allBalances(
  "sif1v89glcgkk6n98yf7xfr5pnrfttuceagsy5jw97"
);

console.log(response);

// response.pools.forEach((pool) => console.log(pool));
