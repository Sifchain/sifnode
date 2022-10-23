const web3 = require("web3");
const BigNumber = web3.BigNumber;
const { ethers } = require("hardhat");
const { use, expect } = require("chai");
const { solidity } = require("ethereum-waffle");

require("chai").use(require("chai-as-promised")).use(require("chai-bignumber")(BigNumber)).should();

use(solidity);

// The values below are from the old implementation (as they should be).
// The idea is to compare the costs of the old impl with the new one.
// The old impl didn't have the capability of returning the full list of blocklisted addresses.
// The new one does have that and it's more expensive because of that.
// Set `use` to true if you want to compare costs.
const gasProfiling = {
  use: true,
  addFirstUser: 45963,
  addAnotherUser: 45963,
  removeOneUser: 15949,
  removeLastUser: 15949,
  current: {},
};

describe("Blocklist", function (accounts) {
  const state = {
    accounts: {
      owner: null,
      userOne: null,
      userTwo: null,
      userThree: null,
    },
    blocklistFactory: null,
    blocklist: null,
  };

  before(async function () {
    accounts = await ethers.getSigners();

    state.blocklistFactory = await ethers.getContractFactory("Blocklist");

    state.accounts.owner = accounts[0];
    state.accounts.userOne = accounts[1];
    state.accounts.userTwo = accounts[2];
    state.accounts.userThree = accounts[3];
  });

  beforeEach(async function () {
    state.blocklist = await state.blocklistFactory.deploy();
    await state.blocklist.deployed();
  });

  describe("Blocklist deployment and basics", function () {
    it("should deploy the Blocklist, correctly setting the owner", async function () {
      const owner = await state.blocklist.owner();
      expect(owner).to.be.equal(state.accounts.owner.address);
    });

    it("should allow any user to query the blocklist", async function () {
      const isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne.address);
      expect(isBlocklisted).to.be.false;
    });
  });

  describe("Administration", function () {
    it("should allow the owner to add an address to the blocklist", async function () {
      // add userOne to the blocklist
      await addSingleAddress(state.accounts.userOne.address, state);

      // check if the user is now blocklisted
      const isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne.address);
      expect(isBlocklisted).to.be.true;
    });

    it("should NOT allow the owner to add an address to the blocklist if it's already there", async function () {
      // add userOne to the blocklist
      await addSingleAddress(state.accounts.userOne.address, state);

      // check if the user is now blocklisted
      const isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne.address);
      expect(isBlocklisted).to.be.true;

      // now, try again with the same address, and fail
      await expect(
        state.blocklist.connect(state.accounts.owner).addToBlocklist(state.accounts.userOne.address)
      ).to.be.rejectedWith("Already in blocklist");
    });

    it("should NOT allow a user to add an address to the blocklist", async function () {
      await expect(
        state.blocklist
          .connect(state.accounts.userOne)
          .addToBlocklist(state.accounts.userTwo.address)
      ).to.be.rejectedWith("Ownable: caller is not the owner");
    });

    it("should allow the owner to batch add addresses to the blocklist", async function () {
      // add three users to the blocklist
      const addressList = [
        state.accounts.userOne.address,
        state.accounts.userTwo.address,
        state.accounts.userThree.address,
      ];

      await batchAddAddresses(addressList, state);

      // check if the users are now blocklisted
      let isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne.address);
      expect(isBlocklisted).to.be.true;

      isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userTwo.address);
      expect(isBlocklisted).to.be.true;

      isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userThree.address);
      expect(isBlocklisted).to.be.true;
    });

    it("should NOT allow the owner to batch add addressed to the blocklist if one of them is already there", async function () {
      // add userTwo to the blocklist
      await addSingleAddress(state.accounts.userTwo.address, state);

      // add three users to the blocklist
      const addressList = [
        state.accounts.userOne.address,
        state.accounts.userTwo.address,
        state.accounts.userThree.address,
      ];
      await expect(
        state.blocklist.connect(state.accounts.owner).batchAddToBlocklist(addressList)
      ).to.be.rejectedWith("Already in blocklist");

      // check if the users are now blocklisted (only userTwo should be)
      let isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne.address);
      expect(isBlocklisted).to.be.false;

      isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userTwo.address);
      expect(isBlocklisted).to.be.true;

      isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userThree.address);
      expect(isBlocklisted).to.be.false;
    });

    it("should NOT allow a user to batch add addresses to the blocklist", async function () {
      const addressList = [
        state.accounts.userOne.address,
        state.accounts.userTwo.address,
        state.accounts.userThree.address,
      ];
      await expect(
        state.blocklist.connect(state.accounts.userOne).batchAddToBlocklist(addressList)
      ).to.be.rejectedWith("Ownable: caller is not the owner");
    });

    it("should allow the owner to remove an address from the blocklist", async function () {
      // add userOne to the blocklist
      await addSingleAddress(state.accounts.userOne.address, state);

      // remove userOne from the blocklist
      await removeSingleAddress(state.accounts.userOne.address, state);

      // check if the user is not blocklisted anymore
      const isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne.address);
      expect(isBlocklisted).to.be.false;
    });

    it("should NOT the owner to remove an address from the blocklist if it's not there", async function () {
      // remove userOne from the blocklist
      await expect(
        state.blocklist
          .connect(state.accounts.owner)
          .removeFromBlocklist(state.accounts.userOne.address)
      ).to.be.rejectedWith("Not in blocklist");
    });

    it("should NOT allow a user to remove an address from the blocklist", async function () {
      // add userOne to the blocklist
      await addSingleAddress(state.accounts.userOne.address, state);

      // try to remove userOne from the blocklist
      await expect(
        state.blocklist
          .connect(state.accounts.userTwo)
          .removeFromBlocklist(state.accounts.userOne.address)
      ).to.be.rejectedWith("Ownable: caller is not the owner");

      // check if the user is still blocklisted
      const isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne.address);
      expect(isBlocklisted).to.be.true;
    });

    it("should allow the owner to batch remove addresses from the blocklist", async function () {
      // add three users to the blocklist
      const addressList = [
        state.accounts.userOne.address,
        state.accounts.userTwo.address,
        state.accounts.userThree.address,
      ];
      await batchAddAddresses(addressList, state);

      // Remove users one and two from the blocklist
      const smallAddressList = [state.accounts.userOne.address, state.accounts.userTwo.address];
      await batchRemoveAddresses(smallAddressList, state);

      // check if the users have been removed from the blocklist
      let isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne.address);
      expect(isBlocklisted).to.be.false;

      isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userTwo.address);
      expect(isBlocklisted).to.be.false;

      // check if user 3 is still blocklisted
      isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userThree.address);
      expect(isBlocklisted).to.be.true;
    });

    it("should NOT allow the owner to batch remove addresses from the blocklist if one of them isn't there", async function () {
      // add two users to the blocklist
      const smallAddressList = [state.accounts.userOne.address, state.accounts.userTwo.address];
      await batchAddAddresses(smallAddressList, state);

      // Try to remove three users from the blocklist
      const addressList = [
        state.accounts.userOne.address,
        state.accounts.userTwo.address,
        state.accounts.userThree.address,
      ];
      await expect(
        state.blocklist.connect(state.accounts.owner).batchRemoveFromBlocklist(addressList)
      ).to.be.rejectedWith("Not in blocklist");

      // check if the users are still in the blocklist
      let isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne.address);
      expect(isBlocklisted).to.be.true;

      isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userTwo.address);
      expect(isBlocklisted).to.be.true;

      // check if user 3 is still not blocklisted
      isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userThree.address);
      expect(isBlocklisted).to.be.false;
    });

    it("should NOT allow a user to batch remove addresses from the blocklist", async function () {
      // add three users to the blocklist
      const addressList = [
        state.accounts.userOne.address,
        state.accounts.userTwo.address,
        state.accounts.userThree.address,
      ];
      await batchAddAddresses(addressList, state);

      // Try to remove users from the blocklist
      await expect(
        state.blocklist.connect(state.accounts.userOne).batchRemoveFromBlocklist(addressList)
      ).to.be.rejectedWith("Ownable: caller is not the owner");

      // check if the users are still in the blocklist
      let isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne.address);
      expect(isBlocklisted).to.be.true;

      isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userTwo.address);
      expect(isBlocklisted).to.be.true;

      isBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userThree.address);
      expect(isBlocklisted).to.be.true;
    });
  });

  describe("Gas costs", function () {
    it("should allow us to measure the cost of adding an address to the blocklist", async function () {
      // add userOne to the blocklist
      const tx = await state.blocklist
        .connect(state.accounts.owner)
        .addToBlocklist(state.accounts.userOne.address);

      const receipt = await tx.wait();
      const gas = Number(receipt.gasUsed);

      printDiff("Add a user to the blocklist", gasProfiling.addFirstUser, gas);
    });

    it("should allow us to measure the cost of adding another address to the blocklist", async function () {
      // add userOne to the blocklist
      await addSingleAddress(state.accounts.userOne.address, state);

      // add userTwo to the blocklist
      const tx = await state.blocklist.addToBlocklist(state.accounts.userTwo.address);

      const receipt = await tx.wait();
      const gas = Number(receipt.gasUsed);
      gasProfiling.current.addAnotherUser = gas;

      printDiff("Add another user to the blocklist", gasProfiling.addAnotherUser, gas);
    });

    it("adding a third user to the blocklist should cost the same as adding the second user", async function () {
      // add userOne to the blocklist
      await addSingleAddress(state.accounts.userOne.address, state);

      // add userTwo to the blocklist
      await addSingleAddress(state.accounts.userTwo.address, state);

      // add userThree to the blocklist
      const tx = await state.blocklist
        .connect(state.accounts.owner)
        .addToBlocklist(state.accounts.userThree.address);

      const receipt = await tx.wait();
      const gas = Number(receipt.gasUsed);

      printDiff("Add a third user to the blocklist", gasProfiling.current.addAnotherUser, gas);
    });

    it("should allow us to measure the cost of removing an address from the blocklist", async function () {
      // add userOne to the blocklist
      await addSingleAddress(state.accounts.userOne.address, state);

      // add userTwo to the blocklist
      await addSingleAddress(state.accounts.userTwo.address, state);

      // remove userOne from the blocklist
      const tx = await state.blocklist
        .connect(state.accounts.owner)
        .removeFromBlocklist(state.accounts.userOne.address);

      const receipt = await tx.wait();
      const gas = Number(receipt.gasUsed);

      printDiff("Remove a user from the blocklist", gasProfiling.removeOneUser, gas);
    });

    it("should allow us to measure the cost of removing the last address from the blocklist", async function () {
      // add userOne to the blocklist
      await addSingleAddress(state.accounts.userOne.address, state);

      // remove userOne from the blocklist
      const tx = await state.blocklist
        .connect(state.accounts.owner)
        .removeFromBlocklist(state.accounts.userOne.address);

      const receipt = await tx.wait();
      const gas = Number(receipt.gasUsed);

      printDiff("Remove the last user from the blocklist", gasProfiling.removeLastUser, gas);
    });

    it("should allow us to compare the cost of adding three addresses to the blocklist VS batch-adding them", async function () {
      let isolatedAddCost = 0;
      let isolatedRemoveCost = 0;

      // add userOne to the blocklist
      let tx = await state.blocklist
        .connect(state.accounts.owner)
        .addToBlocklist(state.accounts.userOne.address);

      let receipt = await tx.wait();
      isolatedAddCost += Number(receipt.gasUsed);

      // add userTwo to the blocklist
      tx = await state.blocklist
        .connect(state.accounts.owner)
        .addToBlocklist(state.accounts.userTwo.address);

      receipt = await tx.wait();
      isolatedAddCost += Number(receipt.gasUsed);

      // add userThree to the blocklist
      tx = await state.blocklist
        .connect(state.accounts.owner)
        .addToBlocklist(state.accounts.userThree.address);

      receipt = await tx.wait();
      isolatedAddCost += Number(receipt.gasUsed);

      // remove userOne from the blocklist
      tx = await state.blocklist
        .connect(state.accounts.owner)
        .removeFromBlocklist(state.accounts.userOne.address);

      receipt = await tx.wait();
      isolatedRemoveCost += Number(receipt.gasUsed);

      // remove userTwo from the blocklist
      tx = await state.blocklist
        .connect(state.accounts.owner)
        .removeFromBlocklist(state.accounts.userTwo.address);

      receipt = await tx.wait();
      isolatedRemoveCost += Number(receipt.gasUsed);

      // remove userThree from the blocklist
      tx = await state.blocklist
        .connect(state.accounts.owner)
        .removeFromBlocklist(state.accounts.userThree.address);

      receipt = await tx.wait();
      isolatedRemoveCost += Number(receipt.gasUsed);

      // Add three users in a batch:
      const addressList = [
        state.accounts.userOne.address,
        state.accounts.userTwo.address,
        state.accounts.userThree.address,
      ];
      tx = await state.blocklist.connect(state.accounts.owner).batchAddToBlocklist(addressList);

      receipt = await tx.wait();
      const batchAddCost = Number(receipt.gasUsed);

      // Remove three users in a batch:
      tx = await state.blocklist
        .connect(state.accounts.owner)
        .batchRemoveFromBlocklist(addressList);

      receipt = await tx.wait();
      const batchRemoveCost = Number(receipt.gasUsed);

      printDiff("Add 3 users separately VS batch add them", isolatedAddCost, batchAddCost);
      printDiff(
        "Remove 3 users separately VS batch remove them",
        isolatedRemoveCost,
        batchRemoveCost
      );
    });
  });

  describe("Convoluted flows", function () {
    it("should allow us to add and remove users to and from the blocklist sequentially and consistently", async function () {
      // Batch add 3 users to the blocklist
      const addressList = [
        state.accounts.userOne.address,
        state.accounts.userTwo.address,
        state.accounts.userThree.address,
      ];
      await batchAddAddresses(addressList, state);

      // Remove the second user from the blocklist
      await removeSingleAddress(state.accounts.userTwo.address, state);

      // Check if data is consistent
      let isUserOneBlocklisted = await state.blocklist.isBlocklisted(
        state.accounts.userOne.address
      );
      let isUserTwoBlocklisted = await state.blocklist.isBlocklisted(
        state.accounts.userTwo.address
      );
      let isUserThreeBlocklisted = await state.blocklist.isBlocklisted(
        state.accounts.userThree.address
      );
      expect(isUserOneBlocklisted).to.be.true;
      expect(isUserTwoBlocklisted).to.be.false;
      expect(isUserThreeBlocklisted).to.be.true;
      let fullList = await state.blocklist.getFullList();
      expect(fullList.length).to.be.equal(2);
      expect(fullList[0]).to.be.equal(state.accounts.userOne.address);
      expect(fullList[1]).to.be.equal(state.accounts.userThree.address);

      // Remove the first user from the blocklist
      await removeSingleAddress(state.accounts.userOne.address, state);

      // Check if data is consistent
      isUserOneBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne.address);
      isUserTwoBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userTwo.address);
      isUserThreeBlocklisted = await state.blocklist.isBlocklisted(
        state.accounts.userThree.address
      );
      expect(isUserOneBlocklisted).to.be.false;
      expect(isUserTwoBlocklisted).to.be.false;
      expect(isUserThreeBlocklisted).to.be.true;
      fullList = await state.blocklist.getFullList();
      expect(fullList.length).to.be.equal(1);
      expect(fullList[0]).to.be.equal(state.accounts.userThree.address);

      // Add the second user to the blocklist
      await addSingleAddress(state.accounts.userTwo.address, state);

      // Check if data is consistent
      isUserOneBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne.address);
      isUserTwoBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userTwo.address);
      isUserThreeBlocklisted = await state.blocklist.isBlocklisted(
        state.accounts.userThree.address
      );
      expect(isUserOneBlocklisted).to.be.false;
      expect(isUserTwoBlocklisted).to.be.true;
      expect(isUserThreeBlocklisted).to.be.true;
      fullList = await state.blocklist.getFullList();
      expect(fullList.length).to.be.equal(2);
      expect(fullList[0]).to.be.equal(state.accounts.userThree.address);
      expect(fullList[1]).to.be.equal(state.accounts.userTwo.address);

      // Add the first user to the blocklist
      await addSingleAddress(state.accounts.userOne.address, state);

      // Check if data is consistent
      isUserOneBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne.address);
      isUserTwoBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userTwo.address);
      isUserThreeBlocklisted = await state.blocklist.isBlocklisted(
        state.accounts.userThree.address
      );
      expect(isUserOneBlocklisted).to.be.true;
      expect(isUserTwoBlocklisted).to.be.true;
      expect(isUserThreeBlocklisted).to.be.true;
      fullList = await state.blocklist.getFullList();
      expect(fullList.length).to.be.equal(3);
      expect(fullList[0]).to.be.equal(state.accounts.userThree.address);
      expect(fullList[1]).to.be.equal(state.accounts.userTwo.address);
      expect(fullList[2]).to.be.equal(state.accounts.userOne.address);

      // Batch remove all users from the blocklist
      await batchRemoveAddresses(addressList, state);

      // Check if data is consistent
      isUserOneBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne.address);
      isUserTwoBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userTwo.address);
      isUserThreeBlocklisted = await state.blocklist.isBlocklisted(
        state.accounts.userThree.address
      );
      expect(isUserOneBlocklisted).to.be.false;
      expect(isUserTwoBlocklisted).to.be.false;
      expect(isUserThreeBlocklisted).to.be.false;
      fullList = await state.blocklist.getFullList();
      expect(fullList.length).to.be.equal(0);

      // Batch add 3 users to the blocklist again
      await batchAddAddresses(addressList, state);

      // Check if data is consistent
      isUserOneBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne.address);
      isUserTwoBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userTwo.address);
      isUserThreeBlocklisted = await state.blocklist.isBlocklisted(
        state.accounts.userThree.address
      );
      expect(isUserOneBlocklisted).to.be.true;
      expect(isUserTwoBlocklisted).to.be.true;
      expect(isUserThreeBlocklisted).to.be.true;
      fullList = await state.blocklist.getFullList();
      expect(fullList.length).to.be.equal(3);
      expect(fullList[0]).to.be.equal(state.accounts.userOne.address);
      expect(fullList[1]).to.be.equal(state.accounts.userTwo.address);
      expect(fullList[2]).to.be.equal(state.accounts.userThree.address);

      // Batch remove users 1 and 3 from the blocklist
      await batchRemoveAddresses(
        [state.accounts.userOne.address, state.accounts.userThree.address],
        state
      );

      // Check if data is consistent
      isUserOneBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userOne.address);
      isUserTwoBlocklisted = await state.blocklist.isBlocklisted(state.accounts.userTwo.address);
      isUserThreeBlocklisted = await state.blocklist.isBlocklisted(
        state.accounts.userThree.address
      );
      expect(isUserOneBlocklisted).to.be.false;
      expect(isUserTwoBlocklisted).to.be.true;
      expect(isUserThreeBlocklisted).to.be.false;
      fullList = await state.blocklist.getFullList();
      expect(fullList.length).to.be.equal(1);
      expect(fullList[0]).to.be.equal(state.accounts.userTwo.address);
    });
  });
});

async function addSingleAddress(address, state) {
  await expect(state.blocklist.connect(state.accounts.owner).addToBlocklist(address))
    .to.emit(state.blocklist, "addedToBlocklist")
    .withArgs(address, state.accounts.owner.address);
}

async function removeSingleAddress(address, state) {
  await expect(state.blocklist.connect(state.accounts.owner).removeFromBlocklist(address))
    .to.emit(state.blocklist, "removedFromBlocklist")
    .withArgs(address, state.accounts.owner.address);
}

async function batchAddAddresses(addressList, state) {
  const tx = await state.blocklist.connect(state.accounts.owner).batchAddToBlocklist(addressList);
  const receipt = await tx.wait();

  // check if events have been emitted correctly
  for (let i = 0; i < addressList.length; i++) {
    const event = receipt.events[i];
    const address = addressList[i];
    expect(event.event).to.be.equal("addedToBlocklist");
    expect(event.args[0]).to.be.equal(address);
  }
}

async function batchRemoveAddresses(addressList, state) {
  const tx = await state.blocklist
    .connect(state.accounts.owner)
    .batchRemoveFromBlocklist(addressList);
  const receipt = await tx.wait();

  // check if events have been emitted correctly
  for (let i = 0; i < addressList.length; i++) {
    const event = receipt.events[i];
    const address = addressList[i];
    expect(event.event).to.be.equal("removedFromBlocklist");
    expect(event.args[0]).to.be.equal(address);
  }
}

async function estimateGas(func) {
  const tx = await func();
  const receipt = await tx.wait();
  return Number(receipt.gasUsed);
}

function printDiff(title, original, current) {
  if (!gasProfiling.use) return;

  const diff = current - original;
  const pct = Math.abs((1 - current / original) * 100).toFixed(2);

  console.log(`______________________________`);
  console.log(`${title}:`);
  console.log(`-> From ${original} to ${current} GAS | Diff: ${diff} (${pct}%)`);
}
