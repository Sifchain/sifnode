module.exports.loadEnv = function() {

    if (process.env.CONSENSUS_THRESHOLD.length === 0) {
      return console.error(
        "Must provide consensus threshold as environment variable."
      );
    }

    if (process.env.OPERATOR.length === 0) {
      return console.error(
        "Must provide operator address as environment variable."
      );
    }

    if (process.env.OWNER.length === 0) {
      return console.error(
        "Must provide owner address as environment variable."
      );
    }
    let owner = process.env.OWNER;
    let pauser = process.env.PAUSER;
    let consensusThreshold = process.env.CONSENSUS_THRESHOLD;
    let operator = process.env.OPERATOR;
    let initialValidators = process.env.INITIAL_VALIDATOR_ADDRESSES.split(",");
    let initialPowers = process.env.INITIAL_VALIDATOR_POWERS.split(",");
  
    if (!initialPowers.length || !initialValidators.length) {
      return console.error(
        "Must provide validator and power."
      );
    }
  
    if (initialPowers.length !== initialValidators.length) {
      return console.error(
        "Each initial validator must have a corresponding power specified."
      );
    }
    initialPowers = initialPowers.map(e => {return parseInt(e)});
  
    return {consensusThreshold, operator, initialValidators, initialPowers, owner, pauser};
}