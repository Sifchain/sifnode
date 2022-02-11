import { expect } from 'chai';
import { ethers } from "hardhat";
import { UnicodeToken, UnicodeToken__factory } from "../build";
import { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers';

describe("Unicode Token", () => {
    let userA: SignerWithAddress, userB: SignerWithAddress, userC: SignerWithAddress, userD: SignerWithAddress;
    let tokenFactory: UnicodeToken__factory;
    let token: UnicodeToken;
    describe("constructor", () => {
        beforeEach(async () => {
            [userA, userB, userC, userD] = await ethers.getSigners();
            tokenFactory = await ethers.getContractFactory("UnicodeToken");
        });
        it("should deploy without any parameters being set", async () => {
            token = await tokenFactory.deploy();
        });
    });
    describe("contract functions and properties", () => {
        beforeEach(async () => {
            [userA, userB, userC, userD] = await ethers.getSigners();
            tokenFactory = await ethers.getContractFactory("UnicodeToken");
            token = await tokenFactory.deploy();
        });
        it("should allow tokens to be minted and transferred like normal", async () => {
            const mintAmount = 10_000;
            const transferAmount = 5_000;

            // Mint Tokens
            expect(await token.balanceOf(userA.address)).to.equal(0);
            await token.mint(userA.address, mintAmount);
            expect(await token.balanceOf(userA.address)).to.equal(mintAmount);
            expect(await token.balanceOf(userB.address)).to.equal(0);

            // Transfer Minted Tokens
            await token.connect(userA).transfer(userB.address, transferAmount);
            expect(await token.balanceOf(userA.address)).to.equal(transferAmount);
            expect(await token.balanceOf(userB.address)).to.equal(transferAmount);
            
        });
        it("should have a unicode symbol and name string", async () => {
            // Should have a nasty unicode string that causes problems on some older systems
            expect(await token.symbol()).to.equal("ܝܘܚܢܢ ܒܝܬ ܐܦܪܝܡ");
            expect(await token.name()).to.equal("لُلُصّبُلُلصّبُررً ॣ ॣh ॣ ॣ 冗");
        });
    });
});