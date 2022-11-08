/**
 * List of colors to be used in the `print` function
 */
const colors = {
  // simple font colors
  black: "\x1b[30m",
  red: "\x1b[31m",
  green: "\x1b[32m",
  yellow: "\x1b[33m",
  blue: "\x1b[34m",
  magenta: "\x1b[35m",
  cyan: "\x1b[36m",
  white: "\x1b[37m",

  // highlights
  h_black: "\x1b[40m\x1b[37m",
  h_red: "\x1b[41m\x1b[37m",
  h_green: "\x1b[42m\x1b[30m",
  h_yellow: "\x1b[43m\x1b[30m",
  h_blue: "\x1b[44m\x1b[37m",
  h_magenta: "\x1b[45m\x1b[37m",
  h_cyan: "\x1b[46m\x1b[30m",
  h_white: "\x1b[47m\x1b[30m",

  // aliases
  highlight: "\x1b[47m\x1b[30m", // white bg and black font

  // mandatory close
  close: "\x1b[0m",
};

/**
 * Prints a colored message on your console/terminal
 * @param {string} color Can be one of the above colors
 * @param {string} message Whatever string
 * @param {bool} breakLine Should it break line after the message?
 */
function print(color, message, breakLine) {
  const lb = breakLine ? "\n" : "";
  console.log(`${colors[color]}${message}${colors.close}${lb}`);
}

/**
 * Will return false for a symbol that has spaces and/or special characters in it
 * @param {string} symbol
 * @returns {bool} does the symbol match the RegExp?
 */
function isValidSymbol(symbol) {
  const regexp = new RegExp("^[a-zA-Z0-9]+$");
  return regexp.test(symbol);
}

/**
 * Receives an object with the following properties, all of which are optional:
 * @param {string} prefix The actual name of the file, something like 'whitelist'
 * @param {string} extension The extension of the file, such as 'json'
 * @param {string} directory The target directory of the file, something like 'data'
 * @returns {string} The generated filename, something like 'data/whitelist_14_sep_2021.json'
 */
function generateTodayFilename({ prefix, extension, directory }) {
  // setup month names
  const monthNames = [
    "Jan",
    "Feb",
    "Mar",
    "Apr",
    "May",
    "Jun",
    "Jul",
    "Aug",
    "Sep",
    "Oct",
    "Nov",
    "Dec",
  ];

  // get current date (we do it manually so that it's not dependant on user's locale)
  const today = new Date();
  const day = String(today.getDate()).padStart(2, "0");
  const month = monthNames[today.getMonth()];
  const year = today.getFullYear();
  let finalDate;

  directory = directory ? `${directory}/` : "";
  prefix = prefix ? `${prefix}_` : "";
  finalDate = `${day}_${month}_${year}`;
  extension = extension ? extension : "json";

  // transform it in a string with the following format:
  // 'myDirectory/whitelist_mainnet_update_14_sep_2021.json' where
  // 'myDirectory' is `directory`
  // 'whitelist_mainnet_update' is `prefix`
  // '14_sep_2021' is today's date
  // and 'json' is `extension`
  const filename = `${directory}${prefix}${finalDate}.${extension}`;
  return filename;
}

/**
 * Busts cache
 * @param {string} url The url to be cacheBusted
 * @returns {string} The same URL with something like '?cacheBuster=95508245028' appended to it
 */
function cacheBuster(url) {
  const rand = Math.floor(Math.random() * (9999999999 - 2) + 1);
  const cacheBuster = `?cacheBuster=${rand}`;
  const finalUrl = `${url}${cacheBuster}`;
  return finalUrl;
}

/**
 * Removes duplicates from arrays
 * @param {array} list Your array
 * @returns {array} an array containing no duplicates
 */
function removeDuplicates(list) {
  const uniqueSet = new Set(list);
  return [...uniqueSet];
}

/**
 * Compares two arrays and returns true if they have the same length and the exact same elements, even if out of order
 * @param {array} List 1
 * @param {array} List 2
 * @returns {bool} Do they have the same elements and length?
 */
function hasSameElementsAndLength(a, b) {
  if (a.length !== b.length) return false;
  const uniqueValues = new Set([...a, ...b]);
  for (const v of uniqueValues) {
    const aCount = a.filter((e) => e === v).length;
    const bCount = b.filter((e) => e === v).length;
    if (aCount !== bCount) return false;
  }
  return true;
}

/**
 * Generates a valid Peggy1 Denom
 * @param {string} symbol The symbol that should be converted to a V1 denom
 * @returns {string} The denom, something like 'ceth'
 */
function generateV1Denom(symbol) {
  const denom = "c" + symbol.toLowerCase();
  return denom;
}

/**
 * Model of an object that the Sifnode team cares about
 */
const SIFNODE_MODEL = {
  decimals: "",
  denom: "",
  base_denom: "",
  path: "",
  ibc_channel_id: "",
  ibc_counterparty_channel_id: "",
  display_name: "",
  display_symbol: "",
  network: "",
  address: "",
  external_symbol: "",
  transfer_limit: "",
  permissions: ["CLP"],
  unit_denom: "",
  ibc_counterparty_denom: "",
  ibc_counterparty_chain_id: "",
};

function extractPrivateKeys(envString) {
  let finalList = [];
  if (envString.indexOf(",") == -1) {
    // there is only one key here
    finalList.push(envString);
  } else {
    finalList = envString.split(",");
  }

  return finalList;
}

module.exports = {
  print,
  isValidSymbol,
  generateTodayFilename,
  cacheBuster,
  removeDuplicates,
  hasSameElementsAndLength,
  generateV1Denom,
  SIFNODE_MODEL,
  extractPrivateKeys,
};
