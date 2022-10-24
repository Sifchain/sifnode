const Blocklist = artifacts.require("Blocklist");
const { expect } = require('chai');

require("chai")
  .use(require("chai-as-promised"))
  .should();


// The values below are from the old implementation (as they should be).
// The idea is to compare the costs of the old impl with the new one.
// The old impl didn't have the capability of returning the full list of blocklisted addresses.
// The new one does have that and it's more expensive because of that.
// Set `use` to true if you want to compare costs.
const gasProfiling = {
  use: false,
  addFirstUser: 45963,
  addAnotherUser: 45963,
  removeOneUser: 15949,
  removeLastUser: 15949,
  current: {}
}

contract("Blocklist", function (accounts) {
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
      
      // try to remove userOne from the blocklist
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

      // check if two events have been emitted correctly
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
      printDiff('Add a user to the blocklist', gasProfiling.addFirstUser, gas);
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
      gasProfiling.current.addAnotherUser = gas;
      printDiff('Add another user to the blocklist', gasProfiling.addAnotherUser, gas);
    });

    it("adding a third user to the blocklist should cost the same as adding the second user", async function () {
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

      // add userThree to the blocklist
      const tx = await state.blocklist.addToBlocklist(
        state.accounts.userThree,
        { from: state.accounts.owner }
      );
      
      const gas = Number(tx.receipt.gasUsed);
      printDiff('Add a third user to the blocklist', gasProfiling.current.addAnotherUser, gas);
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
      printDiff('Remove a user from the blocklist', gasProfiling.removeOneUser, gas);
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

    it("should allow us to compare the cost of adding three addresses to the blocklist VS batch-adding them", async function () {
      let isolatedAddCost = 0;
      let isolatedRemoveCost = 0;
      
      // add userOne to the blocklist
      let tx = await addSingleAddress(state.accounts.userOne, state);
      isolatedAddCost += Number(tx.receipt.gasUsed);

      // add userTwo to the blocklist
      tx = await addSingleAddress(state.accounts.userTwo, state);
      isolatedAddCost += Number(tx.receipt.gasUsed);

      // add userThree to the blocklist
      tx = await addSingleAddress(state.accounts.userThree, state);
      isolatedAddCost += Number(tx.receipt.gasUsed);

      // remove userOne from the blocklist
      tx = await removeSingleAddress(state.accounts.userOne, state);
      isolatedRemoveCost += Number(tx.receipt.gasUsed);

      // remove userTwo from the blocklist
      tx = await removeSingleAddress(state.accounts.userTwo, state);
      isolatedRemoveCost += Number(tx.receipt.gasUsed);

      // remove userThree from the blocklist
      tx = await removeSingleAddress(state.accounts.userThree, state);
      isolatedRemoveCost += Number(tx.receipt.gasUsed);

      // Add three users in a batch:
      const addressList = [state.accounts.userOne, state.accounts.userTwo, state.accounts.userThree];
      tx = await state.blocklist.batchAddToBlocklist(
        addressList,
        { from: state.accounts.owner }
      );
      const batchAddCost = Number(tx.receipt.gasUsed);

      // Remove three users in a batch:
      tx = await state.blocklist.batchRemoveFromBlocklist(
        addressList,
        { from: state.accounts.owner }
      );
      const batchRemoveCost = Number(tx.receipt.gasUsed);

      printDiff('Add 3 users separately VS batch add them', isolatedAddCost, batchAddCost);
      printDiff('Remove 3 users separately VS batch remove them', isolatedRemoveCost, batchRemoveCost);
    });
  });

  describe("Convoluted flows", function () {
    beforeEach(async function () {
      state.blocklist = await Blocklist.new();
    });

    it("should allow us to add and remove users to and from the blocklist sequentially and consistently", async function () {
      // Batch add 3 users to the blocklist
      const addressList = [state.accounts.userOne, state.accounts.userTwo, state.accounts.userThree];
      await expect(state.blocklist.batchAddToBlocklist(
        addressList,
        { from: state.accounts.owner }
      )).to.be.fulfilled;

      // Remove the second user from the blocklist
      await expect(state.blocklist.removeFromBlocklist(
        state.accounts.userTwo,
        { from: state.accounts.owner }
      )).to.be.fulfilled;

      // Check if data is consistent
      let isUserOneBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne);
      let isUserTwoBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userTwo);
      let isUserThreeBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userThree);
      expect(isUserOneBlocklisted).to.be.true;
      expect(isUserTwoBlocklisted).to.be.false;
      expect(isUserThreeBlocklisted).to.be.true;
      let fullList = await state.blocklist.getFullList();
      expect(fullList.length).to.be.equal(2);
      expect(fullList[0]).to.be.equal(state.accounts.userOne);
      expect(fullList[1]).to.be.equal(state.accounts.userThree);

      // Remove the first user from the blocklist
      await expect(state.blocklist.removeFromBlocklist(
        state.accounts.userOne,
        { from: state.accounts.owner }
      )).to.be.fulfilled;

      // Check if data is consistent
      isUserOneBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne);
      isUserTwoBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userTwo);
      isUserThreeBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userThree);
      expect(isUserOneBlocklisted).to.be.false;
      expect(isUserTwoBlocklisted).to.be.false;
      expect(isUserThreeBlocklisted).to.be.true;
      fullList = await state.blocklist.getFullList();
      expect(fullList.length).to.be.equal(1);
      expect(fullList[0]).to.be.equal(state.accounts.userThree);

      // Add the second user to the blocklist
      await expect(state.blocklist.addToBlocklist(
        state.accounts.userTwo,
        { from: state.accounts.owner }
      )).to.be.fulfilled;

      // Check if data is consistent
      isUserOneBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne);
      isUserTwoBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userTwo);
      isUserThreeBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userThree);
      expect(isUserOneBlocklisted).to.be.false;
      expect(isUserTwoBlocklisted).to.be.true;
      expect(isUserThreeBlocklisted).to.be.true;
      fullList = await state.blocklist.getFullList();
      expect(fullList.length).to.be.equal(2);
      expect(fullList[0]).to.be.equal(state.accounts.userThree);
      expect(fullList[1]).to.be.equal(state.accounts.userTwo);

      // Add the first user to the blocklist
      await expect(state.blocklist.addToBlocklist(
        state.accounts.userOne,
        { from: state.accounts.owner }
      )).to.be.fulfilled;

      // Check if data is consistent
      isUserOneBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne);
      isUserTwoBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userTwo);
      isUserThreeBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userThree);
      expect(isUserOneBlocklisted).to.be.true;
      expect(isUserTwoBlocklisted).to.be.true;
      expect(isUserThreeBlocklisted).to.be.true;
      fullList = await state.blocklist.getFullList();
      expect(fullList.length).to.be.equal(3);
      expect(fullList[0]).to.be.equal(state.accounts.userThree);
      expect(fullList[1]).to.be.equal(state.accounts.userTwo);
      expect(fullList[2]).to.be.equal(state.accounts.userOne);

      // Batch remove all users from the blocklist
      await expect(state.blocklist.batchRemoveFromBlocklist(
        addressList,
        { from: state.accounts.owner }
      )).to.be.fulfilled;

      // Check if data is consistent
      isUserOneBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne);
      isUserTwoBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userTwo);
      isUserThreeBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userThree);
      expect(isUserOneBlocklisted).to.be.false;
      expect(isUserTwoBlocklisted).to.be.false;
      expect(isUserThreeBlocklisted).to.be.false;
      fullList = await state.blocklist.getFullList();
      expect(fullList.length).to.be.equal(0);

      // Batch add 3 users to the blocklist again
      await expect(state.blocklist.batchAddToBlocklist(
        addressList,
        { from: state.accounts.owner }
      )).to.be.fulfilled;

      // Check if data is consistent
      isUserOneBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne);
      isUserTwoBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userTwo);
      isUserThreeBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userThree);
      expect(isUserOneBlocklisted).to.be.true;
      expect(isUserTwoBlocklisted).to.be.true;
      expect(isUserThreeBlocklisted).to.be.true;
      fullList = await state.blocklist.getFullList();
      expect(fullList.length).to.be.equal(3);
      expect(fullList[0]).to.be.equal(state.accounts.userOne);
      expect(fullList[1]).to.be.equal(state.accounts.userTwo);
      expect(fullList[2]).to.be.equal(state.accounts.userThree);

      // Batch remove users 1 and 3 from the blocklist
      await expect(state.blocklist.batchRemoveFromBlocklist(
        [state.accounts.userOne, state.accounts.userThree],
        { from: state.accounts.owner }
      )).to.be.fulfilled;

      // Check if data is consistent
      isUserOneBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne);
      isUserTwoBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userTwo);
      isUserThreeBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userThree);
      expect(isUserOneBlocklisted).to.be.false;
      expect(isUserTwoBlocklisted).to.be.true;
      expect(isUserThreeBlocklisted).to.be.false;
      fullList = await state.blocklist.getFullList();
      expect(fullList.length).to.be.equal(1);
      expect(fullList[0]).to.be.equal(state.accounts.userTwo);
    });
  });
});

async function addSingleAddress(address, state) {
  const tx = await state.blocklist.addToBlocklist(
    address,
    { from: state.accounts.owner }
  );

  return tx;
}

async function removeSingleAddress(address, state) {
  const tx = await state.blocklist.removeFromBlocklist(
    address,
    { from: state.accounts.owner }
  );

  return tx;
}

function printDiff(title, original, current) {
  if(!gasProfiling.use) return;

  const diff = current - original;
  const pct = Math.abs(((1 - current / original) * 100)).toFixed(2);

  console.log(`______________________________`);
  console.log(`${title}:`);
  console.log(`-> From ${original} to ${current} GAS | Diff: ${diff} (${pct}%)`);
}