import { expect } from 'chai';
import { ethers } from "hardhat";
import { CommissionToken__factory, CommissionToken } from "../build";
import { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers';
import { exec } from 'child_process';
import { executeLock } from './devenv/evm_lock_burn';

const ZERO_ADDRESS = "0x0000000000000000000000000000000000000000";


describe("Commission Token", () => {
    let userA: SignerWithAddress, userB: SignerWithAddress, userC: SignerWithAddress, userD: SignerWithAddress;
    let devAccount: SignerWithAddress;
    let tokenFactory: CommissionToken__factory;
    let token: CommissionToken;
    describe("constructor", () => {
        beforeEach(async () => {
            [userA, userB, userC, userD, devAccount] = await ethers.getSigners();
            tokenFactory = await ethers.getContractFactory("CommissionToken");
        });
        it("should revert if the dev address is the zero address", async () => {
            const tokenDeploy = tokenFactory.deploy(ZERO_ADDRESS, 500, userA.address, 100_000);
            await expect(tokenDeploy).to.be.revertedWith("Dev account must not be null address");
        });
        it("should revert if the dev fee is over 100% (10,000)", async () => {
            const tokenDeploy = tokenFactory.deploy(devAccount.address, 10_001, userA.address, 100_000);
            await expect(tokenDeploy).to.be.revertedWith("Dev Fee cannot exceed 100%");
        });
        it("should revert if the dev fee is equal to 100% (10,000)", async () => {
            const tokenDeploy = tokenFactory.deploy(devAccount.address, 10_000, userA.address, 100_000);
            await expect(tokenDeploy).to.be.revertedWith("Dev Fee cannot exceed 100%");
        });
        it("should revert if the dev fee is equal to 0% (0)", async () => {
            const tokenDeploy = tokenFactory.deploy(devAccount.address, 0, userA.address, 100_000);
            await expect(tokenDeploy).to.be.revertedWith("Dev Fee cannot be 0%");
        });
        it("should revert if the user address is the zero address", async() => {
            const tokenDeploy = tokenFactory.deploy(devAccount.address, 500, ZERO_ADDRESS, 100_000);
            await expect(tokenDeploy).to.be.revertedWith("Initial minting address must not be null address");
        });
        it("should correctly set the devFee", async () => {
            token = await tokenFactory.deploy(devAccount.address, 500, userA.address, 100_000);
            expect(await token.transferFee()).to.equal(500);
        });
        it("should correctly set the balance of the user", async () => {
            token = await tokenFactory.deploy(devAccount.address, 500, userA.address, 100_000);
            expect(await token.balanceOf(userA.address)).to.equal(100_000);
        });
    });
    describe("contract devFee of 500 (5%)", () => {
        beforeEach(async () => {
            [userA, userB, userC, userD, devAccount] = await ethers.getSigners();
            tokenFactory = await ethers.getContractFactory("CommissionToken");
            token = await tokenFactory.deploy(devAccount.address, 500, userA.address, 100_000);
        });
        it("should charge a 5% commission when transferring between accounts", async () => {
            await token.connect(userA).transfer(userB.address, 10_000);
            expect(await token.balanceOf(userA.address)).to.equal(90_000);
            expect(await token.balanceOf(userB.address)).to.equal(9_500);
            expect(await token.balanceOf(devAccount.address)).to.equal(500);
        });
        it("should charge a 5% commission when doing a transferFrom between accounts", async () => {
            await token.connect(userA).approve(userB.address, 10_000);
            await token.connect(userB).transferFrom(userA.address, userB.address, 10_000);
            expect(await token.balanceOf(userA.address)).to.equal(90_000);
            expect(await token.balanceOf(userB.address)).to.equal(9_500);
            expect(await token.balanceOf(devAccount.address)).to.equal(500);
        });
    });
});