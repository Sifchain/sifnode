type WindowWithPossibleKeplr = typeof window & {
  keplr?: any;
  getOfflineSigner?: any
};

// Todo
type provider = any

// Detect mossible keplr provider from browser
export default (): provider => {
  const win = window as WindowWithPossibleKeplr;

  if (!win) return null;

  if (win.keplr && win.getOfflineSigner) {
    // assign offline signer (they use __proto__ for some reason), so this is not as pretty as i'd like)
    win.keplr.__proto__.getOfflineSigner = win.getOfflineSigner
    return win.keplr
  }

  return null;
};