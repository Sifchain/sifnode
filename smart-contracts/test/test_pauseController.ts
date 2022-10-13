import { ethers } from "hardhat";
import { use, expect, should } from "chai";
import { solidity } from "ethereum-waffle";
import { BridgeBank, BridgeBank__factory, PauseController__factory } from "../build";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { Signer } from "ethers";

use(solidity);

const ZERO_ADDRESS = "0x0000000000000000000000000000000000000000";
const TimeLockDelay = 250; // Just a default used in the tests

describe("Test Pause Controller", function () {
    let BridgeBank: BridgeBank;
    let BridgeBankFactory: BridgeBank__factory;
    let PauseControllerFactory: PauseController__factory;
    let Deployer: SignerWithAddress, Admin1:SignerWithAddress, Admin2:SignerWithAddress, 
        Pauser1: SignerWithAddress, Pauser2: SignerWithAddress,  Pauser3: SignerWithAddress,
        Unpauser1:SignerWithAddress, Unpauser2:SignerWithAddress, Unpauser3: SignerWithAddress,
        Canceler1:SignerWithAddress, Canceler2:SignerWithAddress, PauserAdder1: SignerWithAddress,
        PauserAdder2: SignerWithAddress, UnpauserAdmin1: SignerWithAddress, UnpauserAdmin2: SignerWithAddress,
        CancelerAdder1: SignerWithAddress, CancelerAdder2: SignerWithAddress;

    beforeEach(async () => {
        [
            Deployer,
            Admin1, 
            Admin2, 
            Pauser1, 
            Pauser2, 
            Pauser3, 
            PauserAdder1,
            PauserAdder2,
            Unpauser1, 
            Unpauser2, 
            Unpauser3, 
            UnpauserAdmin1,
            UnpauserAdmin2,
            Canceler1, 
            Canceler2,
            CancelerAdder1,
            CancelerAdder2
        ] = await ethers.getSigners()
        BridgeBankFactory = await ethers.getContractFactory("BridgeBank");
        BridgeBank = await BridgeBankFactory.connect(Deployer).deploy();
        PauseControllerFactory = await ethers.getContractFactory("PauseController");
    });
    describe("Constructor", function () {
        it("should deploy properly with all fields unique and correct", async () => {
            await expect(
                PauseControllerFactory.connect(Deployer).deploy(
                    BridgeBank.address, // _bridgeBank
                    TimeLockDelay, // _timelockDelay
                    [Admin1.address, Admin2.address], // _admins
                    [Pauser1.address, Pauser2.address, Pauser3.address], // _pausers
                    [PauserAdder1.address, PauserAdder2.address], // _pauser_adder
                    [Unpauser1.address, Unpauser2.address, Unpauser3.address], // _unpausers
                    [UnpauserAdmin1.address, UnpauserAdmin2.address], // _unpauser_admin
                    [Canceler1.address, Canceler2.address], // _cancelers
                    [CancelerAdder1.address, CancelerAdder2.address] // _canceler_adder
                )
            ).to.not.be.reverted;
        });
        it("should deploy properly with all fields empty but timelockdelay (set to zero) and bridgebank", async () => {
            await expect(
                PauseControllerFactory.connect(Deployer).deploy(
                    BridgeBank.address, // _bridgeBank
                    0, // _timelockDelay
                    [], // _admins
                    [], // _pausers
                    [], // _pauser_adder
                    [], // _unpausers
                    [], // _unpauser_admin
                    [], // _cancelers
                    [] // _canceler_adder
                )
            ).to.not.be.reverted;
        });
        it("should deploy properly with one account having all the same roles", async () => {
            await expect(
                PauseControllerFactory.connect(Deployer).deploy(
                    BridgeBank.address, // _bridgeBank
                    TimeLockDelay, // _timelockDelay
                    [Admin1.address], // _admins
                    [Admin1.address], // _pausers
                    [Admin1.address], // _pauser_adder
                    [Admin1.address], // _unpausers
                    [Admin1.address], // _unpauser_admin
                    [Admin1.address], // _cancelers
                    [Admin1.address] // _canceler_adder
                )
            ).to.not.be.reverted;
        });
        it("should REVERT with the same user twice in any field", async () => {
            await expect(
                PauseControllerFactory.connect(Deployer).deploy(
                    BridgeBank.address, // _bridgeBank
                    TimeLockDelay, // _timelockDelay
                    [Admin1.address, Admin1.address], // _admins
                    [], // _pausers
                    [], // _pauser_adder
                    [], // _unpausers
                    [], // _unpauser_admin
                    [], // _cancelers
                    [] // _canceler_adder
                )
            ).to.be.reverted;
        });
        it("should REVERT with a null bridgebank address", async () => {
            await expect(
                PauseControllerFactory.connect(Deployer).deploy(
                    ZERO_ADDRESS, // _bridgeBank
                    TimeLockDelay, // _timelockDelay
                    [Admin1.address, Admin2.address], // _admins
                    [Pauser1.address, Pauser2.address, Pauser3.address], // _pausers
                    [PauserAdder1.address, PauserAdder2.address], // _pauser_adder
                    [Unpauser1.address, Unpauser2.address, Unpauser3.address], // _unpausers
                    [UnpauserAdmin1.address, UnpauserAdmin2.address], // _unpauser_admin
                    [Canceler1.address, Canceler2.address], // _cancelers
                    [CancelerAdder1.address, CancelerAdder2.address] // _canceler_adder
                )
            ).to.be.revertedWith("BridgeBank address must be set");
        });
        it("It should meh", async () => {
            await expect(
                PauseControllerFactory.connect(Deployer).deploy(
                    BridgeBank.address, // _bridgeBank
                    TimeLockDelay, // _timelockDelay
                    [Admin1.address, Admin2.address], // _admins
                    [Pauser1.address, Pauser2.address, Pauser3.address], // _pausers
                    [PauserAdder1.address, PauserAdder2.address], // _pauser_adder
                    [Unpauser1.address, Unpauser2.address, Unpauser3.address], // _unpausers
                    [UnpauserAdmin1.address, UnpauserAdmin2.address], // _unpauser_admin
                    [Canceler1.address, Canceler2.address], // _cancelers
                    [CancelerAdder1.address, CancelerAdder2.address] // _canceler_adder

                )
            ).to.be.revertedWith("BridgeBank address must be set");
        });
    })
});