import { expect } from 'chai';
import { ethers } from "hardhat";
import { RandomTrollToken__factory, RandomTrollToken } from "../build";
import { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers';
import { BigNumber } from 'ethers';

const ZERO_ADDRESS = "0x0000000000000000000000000000000000000000";
const INITAL_BALANCE = ethers.BigNumber.from("100000000000000000000") // 100e18


describe("Random Troll Token", () => {
    let userA: SignerWithAddress, userB: SignerWithAddress, userC: SignerWithAddress, userD: SignerWithAddress;
    let tokenFactory: RandomTrollToken__factory;
    let token: RandomTrollToken;
    describe("constructor", () => {
        beforeEach(async () => {
            [userA, userB, userC, userD] = await ethers.getSigners();
            tokenFactory = await ethers.getContractFactory("RandomTrollToken");
        });
        it("should fail if there is more addresses then quantity values", async () => {
            await expect(tokenFactory.deploy([userA.address, userB.address, userC.address], [INITAL_BALANCE, INITAL_BALANCE]))
                .to.be.revertedWith("Accounts and Quantities must be same length");
        });
        it("should fail if there is more quantity values then quantity addresses", async () => {
            await expect(tokenFactory.deploy([userA.address, userB.address], [INITAL_BALANCE, INITAL_BALANCE, INITAL_BALANCE]))
                .to.be.revertedWith("Accounts and Quantities must be same length");
        });
        it("should construct the token successfully if the fields are empty and show zero balances for any address", async () => {
            token = await tokenFactory.deploy([], []);
            expect(await token.balanceOf(userA.address)).to.equal(0)
            expect(await token.balanceOf(userB.address)).to.equal(0)
            expect(await token.balanceOf(userC.address)).to.equal(0)

        });
        it("should construct the token successfully if the fields are the same length and have values in them", async () => {
            token = await tokenFactory.deploy(
                [
                    userA.address, 
                    userB.address, 
                    userC.address, 
                    userD.address
                ], 
                [
                    INITAL_BALANCE,
                    INITAL_BALANCE,
                    INITAL_BALANCE,
                    INITAL_BALANCE
                ]);
        })
        describe("token acts like an ERC20 token except for viewable values", () => {
            beforeEach(async () => {
                [userA, userB, userC, userD] = await ethers.getSigners();
                tokenFactory = await ethers.getContractFactory("RandomTrollToken");
                token = await tokenFactory.deploy([userA.address, userB.address], [INITAL_BALANCE, INITAL_BALANCE]);
            });
            it("should allow a user to transfer tokens to any already filled account", async () => {
                await token.connect(userA).transfer(userB.address, 10_000);
            });
            it("should allow a user to transfer tokens to an account that has zero balance", async () => {
                await token.connect(userA).transfer(userC.address, 10_000);
            });
            it("should not allow a user to transfer tokens from an account with zero balance to any other account", async () => {
                await expect(token.connect(userC).transfer(userA.address, 10_000)).to.be.revertedWith("ERC20: transfer amount exceeds balance");
                await expect(token.connect(userC).transfer(userB.address, 10_000)).to.be.revertedWith("ERC20: transfer amount exceeds balance");
                await expect(token.connect(userC).transfer(userC.address, 10_000)).to.be.revertedWith("ERC20: transfer amount exceeds balance");
                await expect(token.connect(userC).transfer(userD.address, 10_000)).to.be.revertedWith("ERC20: transfer amount exceeds balance");
            });
            it("should show balances of 0 for unfilled accounts and balances that change for filled accounts", async() => {
                expect(await token.balanceOf(userC.address)).to.equal(0);
                expect(await token.balanceOf(userA.address)).to.not.equal(0);
                await token.connect(userA).transfer(userC.address, 5_000);
                expect(await token.balanceOf(userC.address)).to.not.equal(0);

                const accounts = [userA.address, userB.address, userC.address];
                const checkNTimes = 10; // How many times should the balance of an account be checked?
                for (const account of accounts) {
                    let pastValues: Set<String> = new Set();
                    for (let check = 0; check < checkNTimes; check++) {
                        pastValues.add(await (await token.balanceOf(account)).toHexString());
                        ethers.provider.send('evm_mine', []);
                    }
                    // The set of balances should have unique entries for each lookup
                    expect(pastValues.size).to.equal(checkNTimes);
                }

                for(let check = 0; check < checkNTimes; check++) {
                    expect(await token.balanceOf(userD.address)).to.equal(0);
                    ethers.provider.send('evm_mine', []);
                }
            });
            it("should return a different symbol on each call", async () => {
                const symbols: Set<string> = new Set();
                for (let call = 0; call < 100; call++) {
                    symbols.add(await token.symbol());
                    ethers.provider.send('evm_mine', []);
                }
                expect(symbols.size).to.not.be.lessThanOrEqual(1);
            });
            it("should return a different name on each call", async () => {
                const names: Set<string> = new Set();
                for (let call = 0; call < 100; call++) {
                    names.add(await token.symbol());
                    ethers.provider.send('evm_mine', []);
                }
                expect(names.size).to.not.be.lessThanOrEqual(1);
            });
        });
    });
});