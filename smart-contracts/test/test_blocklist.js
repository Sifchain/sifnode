const Blocklist = artifacts.require("Blocklist");
const { expect } = require('chai');

require("chai")
  .use(require("chai-as-promised"))
  .should();

const gasProfiling = {
  use: false,
  addFirstUser: 45963,
  addAnotherUser: 45963,
  removeOneUser: 15949,
  removeLastUser: 15949,
}

contract.only("Blocklist", function (accounts) {
  const state = {
    accounts: {
      owner: accounts[0],
      userOne: accounts[1],
      userTwo: accounts[2],
      userThree: accounts[3],
    },
    blocklist: null
  };

  describe("Blocklist deployment and basics", function () {
    beforeEach(async function () {
      state.blocklist = await Blocklist.new();
    });

    it("should deploy the Blocklist, correctly setting the owner", async function () {
      const owner = await state.blocklist.owner();
      expect(owner).to.be.equal(state.accounts.owner);
    });

    it("should allow any user to query the blocklist", async function () {
      const isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne);
      expect(isBlocklisted).to.be.false;
    });
  });

  describe("Administration", function () {
    beforeEach(async function () {
      state.blocklist = await Blocklist.new();
    });

    it("should allow the owner to add an address to the blocklist", async function () {
      // add userOne to the blocklist
      const { logs } = await state.blocklist.addToBlocklist(
        state.accounts.userOne,
        { from: state.accounts.owner }
      );
      
      // check if an event has been emitted correctly
      const event = logs.find(e => e.event === "addedToBlocklist");
      expect(event.args.account).to.be.equal(state.accounts.userOne);
      expect(event.args.by).to.be.equal(state.accounts.owner);
      
      // check if the user is now blocklisted
      const isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne);
      expect(isBlocklisted).to.be.true;
    });

    it("should NOT allow the owner to add an address to the blocklist if it's already there", async function () {
      // add userOne to the blocklist
      const { logs } = await state.blocklist.addToBlocklist(
        state.accounts.userOne,
        { from: state.accounts.owner }
      );
      
      // check if an event has been emitted correctly
      const event = logs.find(e => e.event === "addedToBlocklist");
      expect(event.args.account).to.be.equal(state.accounts.userOne);
      expect(event.args.by).to.be.equal(state.accounts.owner);
      
      // check if the user is now blocklisted
      const isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne);
      expect(isBlocklisted).to.be.true;

      // now, try again with the same address, and fail
      await expect(state.blocklist.addToBlocklist(
        state.accounts.userOne,
        { from: state.accounts.owner }
      )).to.be.rejectedWith('Already in blocklist');
    });

    it("should NOT allow a user to add an address to the blocklist", async function () {
      await expect(state.blocklist.addToBlocklist(
        state.accounts.userTwo,
        { from: state.accounts.userOne }
      )).to.be.rejectedWith('Ownable: caller is not the owner.');
    });

    it("should allow the owner to batch add addresses to the blocklist", async function () {
      // add three users to the blocklist
      const addressList = [state.accounts.userOne, state.accounts.userTwo, state.accounts.userThree];
      const { logs } = await state.blocklist.batchAddToBlocklist(
        addressList,
        { from: state.accounts.owner }
      );
      
      // check if three events have been emitted correctly
      for (let i = 0; i < logs.length; i++) {
        const log = logs[i];
        const address = addressList[i];
        expect(log.args.account).to.be.equal(address);
      }
      
      // check if the users are now blocklisted
      let isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne);
      expect(isBlocklisted).to.be.true;

      isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userTwo);
      expect(isBlocklisted).to.be.true;

      isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userThree);
      expect(isBlocklisted).to.be.true;
    });

    it("should NOT allow the owner to batch add addressed to the blocklist if one of them is already there", async function () {
      // add userTwo to the blocklist
      await state.blocklist.addToBlocklist(
        state.accounts.userTwo,
        { from: state.accounts.owner }
      );
      
      // add three users to the blocklist
      const addressList = [state.accounts.userOne, state.accounts.userTwo, state.accounts.userThree];
      await expect(state.blocklist.batchAddToBlocklist(
        addressList,
        { from: state.accounts.owner }
      )).to.be.rejectedWith('Already in blocklist');

      // check if the users are now blocklisted (only userTwo should be)
      let isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne);
      expect(isBlocklisted).to.be.false;

      isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userTwo);
      expect(isBlocklisted).to.be.true;

      isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userThree);
      expect(isBlocklisted).to.be.false;
    });

    it("should NOT allow a user to batch add addresses to the blocklist", async function () {
      const addressList = [state.accounts.userOne, state.accounts.userTwo, state.accounts.userThree];
      await expect(state.blocklist.batchAddToBlocklist(
        addressList,
        { from: state.accounts.userOne }
      )).to.be.rejectedWith('Ownable: caller is not the owner.');
    });

    it("should allow the owner to remove an address from the blocklist", async function () {
      // add userOne to the blocklist
      await expect(state.blocklist.addToBlocklist(
        state.accounts.userOne,
        { from: state.accounts.owner }
      )).to.be.fulfilled;
      
      // remove userOne from the blocklist
      const { logs } = await state.blocklist.removeFromBlocklist(
        state.accounts.userOne,
        { from: state.accounts.owner }
      );
      
      // check if an event has been emitted correctly
      const event = logs.find(e => e.event === "removedFromBlocklist");
      expect(event.args.account).to.be.equal(state.accounts.userOne);
      expect(event.args.by).to.be.equal(state.accounts.owner);
      
      // check if the user is not blocklisted anymore
      const isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne);
      expect(isBlocklisted).to.be.false;
    });

    it("should NOT the owner to remove an address from the blocklist if it's not there", async function () {
      // remove userOne from the blocklist
      await expect(state.blocklist.removeFromBlocklist(
        state.accounts.userOne,
        { from: state.accounts.owner }
      )).to.be.rejectedWith('Not in blocklist');
    });

    it("should NOT allow a user to remove an address from the blocklist", async function () {
      // add userOne to the blocklist
      await expect(state.blocklist.addToBlocklist(
        state.accounts.userOne,
        { from: state.accounts.owner }
      )).to.be.fulfilled;
      
      // try yo remove userOne from the blocklist
      await expect(state.blocklist.removeFromBlocklist(
        state.accounts.userOne,
        { from: state.accounts.userTwo }
      )).to.be.rejectedWith('Ownable: caller is not the owner.');
      
      // check if the user is still blocklisted
      const isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne);
      expect(isBlocklisted).to.be.true;
    });

    it("should allow the owner to batch remove addresses from the blocklist", async function () {
      // add three users to the blocklist
      const addressList = [state.accounts.userOne, state.accounts.userTwo, state.accounts.userThree];
      await state.blocklist.batchAddToBlocklist(
        addressList,
        { from: state.accounts.owner }
      );
      
      // Remove users one and two from the blocklist
      const smallAddressList = [state.accounts.userOne, state.accounts.userTwo];
      const { logs } = await state.blocklist.batchRemoveFromBlocklist(
        smallAddressList,
        { from: state.accounts.owner }
      );

      // check if three events have been emitted correctly
      for (let i = 0; i < logs.length; i++) {
        const log = logs[i];
        const address = smallAddressList[i];
        expect(log.args.account).to.be.equal(address);
      }

      // check if the users have been removed from the blocklist
      let isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne);
      expect(isBlocklisted).to.be.false;

      isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userTwo);
      expect(isBlocklisted).to.be.false;

      // check if user 3 is still blocklisted
      isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userThree);
      expect(isBlocklisted).to.be.true;
    });

    it("should NOT allow the owner to batch remove addresses from the blocklist if one of them isn't there", async function () {
      // add two users to the blocklist
      const smallAddressList = [state.accounts.userOne, state.accounts.userTwo];
      await state.blocklist.batchAddToBlocklist(
        smallAddressList,
        { from: state.accounts.owner }
      );
      
      // Try to remove three users from the blocklist
      const addressList = [state.accounts.userOne, state.accounts.userTwo, state.accounts.userThree];
      await expect(state.blocklist.batchRemoveFromBlocklist(
        addressList,
        { from: state.accounts.owner }
      )).to.be.rejectedWith('Not in blocklist');

      // check if the users are still in the blocklist
      let isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne);
      expect(isBlocklisted).to.be.true;

      isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userTwo);
      expect(isBlocklisted).to.be.true;

      // check if user 3 is still not blocklisted
      isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userThree);
      expect(isBlocklisted).to.be.false;
    });

    it("should NOT allow a user to batch remove addresses from the blocklist", async function () {
      // add three users to the blocklist
      const addressList = [state.accounts.userOne, state.accounts.userTwo, state.accounts.userThree];
      await state.blocklist.batchAddToBlocklist(
        addressList,
        { from: state.accounts.owner }
      );
      
      // Remove users one and two from the blocklist
      await expect(state.blocklist.batchRemoveFromBlocklist(
        addressList,
        { from: state.accounts.userOne }
      )).to.be.rejectedWith('Ownable: caller is not the owner.');
      
      // check if the users are still in the blocklist
      let isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne);
      expect(isBlocklisted).to.be.true;

      isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userTwo);
      expect(isBlocklisted).to.be.true;

      // check if user 3 is still blocklisted
      isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userThree);
      expect(isBlocklisted).to.be.true;
    });
  });

  describe("Gas costs", function () {
    beforeEach(async function () {
      state.blocklist = await Blocklist.new();
    });

    it("should allow us to measure the cost of adding an address to the blocklist", async function () {
      // add userOne to the blocklist
      const tx = await state.blocklist.addToBlocklist(
        state.accounts.userOne,
        { from: state.accounts.owner }
      );
      
      const gas = Number(tx.receipt.gasUsed);
      printDiff('Remove the last user from the blocklist', gasProfiling.addFirstUser, gas);
    });

    it("should allow us to measure the cost of adding another address to the blocklist", async function () {
      // add userOne to the blocklist
      await state.blocklist.addToBlocklist(
        state.accounts.userOne,
        { from: state.accounts.owner }
      );

      // add userTwo to the blocklist
      const tx = await state.blocklist.addToBlocklist(
        state.accounts.userTwo,
        { from: state.accounts.owner }
      );
      
      const gas = Number(tx.receipt.gasUsed);
      printDiff('Remove the last user from the blocklist', gasProfiling.addAnotherUser, gas);
    });

    it("should allow us to measure the cost of removing an address from the blocklist", async function () {
      // add userOne to the blocklist
      await state.blocklist.addToBlocklist(
        state.accounts.userOne,
        { from: state.accounts.owner }
      );

      // add userTwo to the blocklist
      await state.blocklist.addToBlocklist(
        state.accounts.userTwo,
        { from: state.accounts.owner }
      );

      // remove userOne from the blocklist
      const tx = await state.blocklist.removeFromBlocklist(
        state.accounts.userOne,
        { from: state.accounts.owner }
      );
      
      const gas = Number(tx.receipt.gasUsed);
      printDiff('Remove the last user from the blocklist', gasProfiling.removeOneUser, gas);
    });

    it("should allow us to measure the cost of removing the last address from the blocklist", async function () {
      // add userOne to the blocklist
      await state.blocklist.addToBlocklist(
        state.accounts.userOne,
        { from: state.accounts.owner }
      );

      // remove userOne from the blocklist
      const tx = await state.blocklist.removeFromBlocklist(
        state.accounts.userOne,
        { from: state.accounts.owner }
      );
    
      const gas = Number(tx.receipt.gasUsed);
      printDiff('Remove the last user from the blocklist', gasProfiling.removeLastUser, gas);
    });
  });
});

function printDiff(title, original, current) {
  const diff = current - original;
  const pct = Math.abs(((1 - current / original) * 100)).toFixed(2);

  console.log(`______________________________`);
  console.log(`${title}:`);
  console.log(`-> From ${original} to ${current} GAS | Diff: ${diff} (${pct}%)`);
}