import web3 from "web3";

export function fixSignature(signature: string) {
  // in geth its always 27/28, in ganache its 0/1. Change to 27/28 to prevent
  // signature malleability if version is 0/1
  // see https://github.com/ethereum/go-ethereum/blob/v1.8.23/internal/ethapi/api.go#L465
  let v = parseInt(signature.slice(130, 132), 16);
  if (v < 27) {
    v += 27;
  }
  const vHex = v.toString(16);
  return signature.slice(0, 130) + vHex;
}

export function toEthSignedMessageHash(messageHex: string) {
  const messageBuffer = Buffer.from(messageHex.substring(2), "hex");
  const prefix = Buffer.from(`\u0019Ethereum Signed Message:\n${messageBuffer.length}`);
  return web3.utils.sha3(Buffer.concat([prefix, messageBuffer]).toString());
}

/**
 * Used to colorize logs without using libs
 * @dev Start your string with the color of choice and end it with .close
 * @dev Example: console.log(`${colors.green}Your message here${colors.close}`);
 */
const colors = {
  green: "\x1b[32m",
  red: "\x1b[41m\x1b[37m",
  white: "\x1b[37m",
  highlight: "\x1b[47m\x1b[30m",
  yellow: "\x1b[33m",
  blue: "\x1b[34m",
  magenta: "\x1b[35m",
  cyan: "\x1b[36m",
  close: "\x1b[0m",
};

export type Colors = keyof typeof colors;


/**
 * Colorizes and prints logs without using libs
 * @dev Example: colorLog('green', message);
 */
export function colorLog(colorName: Colors, message: string) {
  console.log(`${colors[colorName]}${message}${colors.close}`);
}
