// Prints out all the ganache accounts

module.exports = async (cb) => {
    const accounts = await web3.eth.getAccounts();
    const result = {
        accounts,
    }
    console.log(JSON.stringify(result));
    return cb();
};
