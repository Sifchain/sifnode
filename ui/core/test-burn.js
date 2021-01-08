const { decodeSignature } = require("@cosmjs/launchpad");
const { toUtf8 } = require("@cosmjs/encoding");
const { Secp256k1, Sha256, Secp256k1Signature } = require("@cosmjs/crypto");

// const output = {
//   signedTx: {
//     msg: [
//       {
//         type: "ethbridge/MsgBurn",
//         value: {
//           cosmos_sender: "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
//           amount: "2000000000000000000",
//           symbol: "ceth",
//           ethereum_chain_id: "5777",
//           ethereum_receiver: "0x627306090abaB3A6e1400e9345bC60c78a8BEf57",
//         },
//       },
//     ],
//     fee: { amount: [], gas: "200000" },
//     memo: "",
//     signatures: [
//       {
//         pub_key: {
//           type: "tendermint/PubKeySecp256k1",
//           value: "A0mB4PyE5XeS3sNpFXIX536INyNoJHkMu1DEQ8FgH8Mq",
//         },
//         signature:
//           "JDBFYdDTZIa/h1gR59aUdbNr5SS2QQhNj+8RIDQBeCdKZ7K03BqfICiYKPCW8lRAv/qgm9aKsaDbjIVNrOxXog==",
//       },
//     ],
//   },
//   txHash: {
//     height: 0,
//     transactionHash:
//       "862299115853568FB5C583D2B60BB8C49DA85F1CFC04A37E780B967AB53BA2BA",
//     code: 4,
//     rawLog:
//       "unauthorized: signature verification failed; verify correct account sequence and chain-id",
//   },
// };

async function run() {
  const docToSign =
    '{"account_number":"4","chain_id":"sifchain","fee":{"amount":[],"gas":"200000"},"memo":"","msgs":[{"type":"ethbridge/MsgBurn","value":{"amount":"2000000000000000000","cosmos_sender":"sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5","ethereum_chain_id":"5777","ethereum_receiver":"0x627306090abaB3A6e1400e9345bC60c78a8BEf57","symbol":"ceth"}}],"sequence":"3"}';

  const decoded = decodeSignature({
    pub_key: {
      type: "tendermint/PubKeySecp256k1",
      value: "A0mB4PyE5XeS3sNpFXIX536INyNoJHkMu1DEQ8FgH8Mq",
    },
    signature:
      "JDBFYdDTZIa/h1gR59aUdbNr5SS2QQhNj+8RIDQBeCdKZ7K03BqfICiYKPCW8lRAv/qgm9aKsaDbjIVNrOxXog==",
  });

  const signature = Secp256k1Signature.fromFixedLength(decoded.signature);
  const content = new Sha256(toUtf8(docToSign));
  const verified = Secp256k1.verifySignature(
    signature,
    content,
    decoded.pubkey
  );
  if (verified) {
    console.log("Verified");
  }
}

run();
