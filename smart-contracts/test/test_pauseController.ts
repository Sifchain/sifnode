import { ethers } from "hardhat";
import hre from "hardhat";
import { use, expect, should } from "chai";
import { solidity } from "ethereum-waffle";
import { BridgeBank, BridgeBank__factory, PauseController, PauseController__factory } from "../build";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { Signer } from "ethers";
import { BytesLike } from "ethers";
import { string } from "yargs";
import { BigNumber } from "ethers";

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
        it("should deploy with the same user twice in any field", async () => {
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
            ).to.not.be.reverted;
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
        it("should correctly set the timelock values and bridgebank address", async () => {
            const PauseController = await PauseControllerFactory.connect(Deployer).deploy(
                    BridgeBank.address, // _bridgeBank
                    TimeLockDelay, // _timelockDelay
                    [Admin1.address, Admin2.address], // _admins
                    [Pauser1.address, Pauser2.address, Pauser3.address], // _pausers
                    [PauserAdder1.address, PauserAdder2.address], // _pauser_adder
                    [Unpauser1.address, Unpauser2.address, Unpauser3.address], // _unpausers
                    [UnpauserAdmin1.address, UnpauserAdmin2.address], // _unpauser_admin
                    [Canceler1.address, Canceler2.address], // _cancelers
                    [CancelerAdder1.address, CancelerAdder2.address] // _canceler_adder
            );
            expect(await PauseController.BridgeBank()).to.equal(BridgeBank.address);
            expect(await PauseController.TimeLockDelay()).to.equal(TimeLockDelay);
            const NO_REQUEST = await PauseController.NOREQUEST();
            expect(NO_REQUEST).to.equal(1); // 1 is the default NO_REQUEST ENUM
            expect(await PauseController.UnpauseRequestBlockHeight()).to.equal(NO_REQUEST);
        });
    });

    let PauseController: PauseController;
    let PAUSER: string, CANCELER: string, ADMIN: string, PAUSER_ADDER: string, CANCELER_ADDER: string,
        UNPAUSER: string, UNPAUSER_ADMIN: string, NO_REQUEST: BigNumber;
    beforeEach(async() => {
        // One Pauser preset, no cancelers preset for adder testing
        PauseController = await PauseControllerFactory.connect(Deployer).deploy(
                BridgeBank.address, // _bridgeBank
                TimeLockDelay, // _timelockDelay
                [Admin1.address], // _admins
                [Pauser1.address], // _pausers
                [PauserAdder1.address], // _pauser_adder
                [Unpauser1.address], // _unpausers
                [UnpauserAdmin1.address], // _unpauser_admin
                [Canceler1.address], // _cancelers
                [CancelerAdder1.address] // _canceler_adder
        );
        await BridgeBank["initialize(address,address,address,address)"](Deployer.address, ZERO_ADDRESS, Deployer.address, Deployer.address);
        await BridgeBank.connect(Deployer).addPauser(PauseController.address);
        const pauserPromise = PauseController.PAUSER();
        const cancelerPromise = PauseController.CANCELER()
        const adminPromise = PauseController.DEFAULT_ADMIN_ROLE();
        const pauserAdderPromise = PauseController.PAUSER_ADDER();
        const cancelerAdderPromise = PauseController.CANCELER_ADDER();
        const unpauserPromise = PauseController.UNPAUSER();
        const unpauserAdminPromise = PauseController.UNPAUSER_ADMIN();
        const noRequestPromise = PauseController.NOREQUEST();

        [PAUSER, CANCELER, ADMIN, PAUSER_ADDER, CANCELER_ADDER, UNPAUSER, UNPAUSER_ADMIN, NO_REQUEST] = await Promise.all
            ([pauserPromise, cancelerPromise, adminPromise, pauserAdderPromise, cancelerAdderPromise, unpauserPromise, unpauserAdminPromise, noRequestPromise]);
    });

    describe("Adder Accounts", () => {
        it("pauser_adder should be able to add pausers by calling addPauser", async () => {
            expect(await PauseController.connect(PauserAdder1).addPauser(Pauser2.address))
              .to.emit(PauseController, "RoleGranted").withArgs(PAUSER, Pauser2.address, PauserAdder1.address);
        });
        it("canceler_adder should be able to add cancelers by calling addCanceler", async () => {
            expect(await PauseController.connect(CancelerAdder1).addCanceler(Canceler2.address))
              .to.emit(PauseController, "RoleGranted").withArgs(CANCELER, Canceler2.address, CancelerAdder1.address);
        });
        it("pauser_adder should not be able to remove pausers", async () => {
            await expect(PauseController.connect(PauserAdder1).revokeRole(PAUSER, Pauser1.address))
                .to.be.reverted;
        });
        it("canceler_adder should not be able to remove cancelers", async () => {
            await expect(PauseController.connect(CancelerAdder1).revokeRole(PAUSER, Canceler1.address))
                .to.be.reverted;
        });
        it("pauser_adder should not be able to add or remove other pauser_adders", async () => {
            await expect(PauseController.connect(PauserAdder1).revokeRole(PAUSER_ADDER, PauserAdder2.address))
                .to.be.reverted;
            await expect(PauseController.connect(PauserAdder1).grantRole(PAUSER_ADDER, Pauser1.address))
                .to.be.reverted;
        });
        it("canceler_adder should not be able to add or remove other canceler_adders", async () => {
            await expect(PauseController.connect(CancelerAdder1).revokeRole(CANCELER_ADDER, CancelerAdder2.address))
                .to.be.reverted;
            await expect(PauseController.connect(CancelerAdder1).grantRole(CANCELER_ADDER, Canceler1.address))
                .to.be.reverted;
        });
        it("should not let anyone without a pauser_adder role from calling addPauser", async () => {
            await expect(PauseController.connect(Deployer).addPauser(Pauser2.address))
              .to.be.reverted;
            await expect(PauseController.connect(Pauser1).addPauser(Pauser2.address))
              .to.be.reverted;
            await expect(PauseController.connect(Pauser2).addPauser(Pauser2.address))
              .to.be.reverted;
            await expect(PauseController.connect(CancelerAdder1).addPauser(Pauser2.address))
              .to.be.reverted;
            await expect(PauseController.connect(Unpauser1).addPauser(Pauser2.address))
              .to.be.reverted;
        });
        it("should not let anyone without a canceler_adder role from calling addCanceler", async () => {
            await expect(PauseController.connect(Deployer).addCanceler(Canceler2.address))
              .to.be.reverted;
            await expect(PauseController.connect(Canceler1).addCanceler(Canceler2.address))
              .to.be.reverted;
            await expect(PauseController.connect(Pauser1).addCanceler(Canceler2.address))
              .to.be.reverted;
            await expect(PauseController.connect(PauserAdder1).addCanceler(Canceler2.address))
              .to.be.reverted;
            await expect(PauseController.connect(Unpauser1).addCanceler(Canceler2.address))
              .to.be.reverted;
        });
    });

    describe("Admin Accounts", () => {
        it("should allow an admin account to add pausers and cancelers with the special addPauser/addCanceler functions", async () => {
            expect(await PauseController.connect(Admin1).addPauser(Pauser2.address))
              .to.emit(PauseController, "RoleGranted").withArgs(PAUSER, Pauser2.address, Admin1.address);
            expect(await PauseController.connect(Admin1).addCanceler(Canceler2.address))
              .to.emit(PauseController, "RoleGranted").withArgs(CANCELER, Canceler2.address, Admin1.address);
        });
        it("should allow an admin account to add/revoke pausers with grantRole/revokeRole", async () => {
            expect(await PauseController.connect(Admin1).grantRole(PAUSER, Pauser2.address))
              .to.emit(PauseController, "RoleGranted").withArgs(PAUSER, Pauser2.address, Admin1.address);
            expect(await PauseController.connect(Admin1).revokeRole(PAUSER, Pauser2.address))
              .to.emit(PauseController, "RoleRevoked").withArgs(PAUSER, Pauser2.address, Admin1.address);
        });
        it("should allow an admin account to add/revoke cancelers with grantRole/revokeRole", async() => {
            expect(await PauseController.connect(Admin1).grantRole(CANCELER, Canceler2.address))
              .to.emit(PauseController, "RoleGranted").withArgs(CANCELER, Canceler2.address, Admin1.address);
            expect(await PauseController.connect(Admin1).revokeRole(CANCELER, Canceler2.address))
              .to.emit(PauseController, "RoleRevoked").withArgs(CANCELER, Canceler2.address, Admin1.address);
        });
        it("should allow an admin account to add/revoke unpausers with grantRole/revokeRole", async () => {
            expect(await PauseController.connect(Admin1).grantRole(UNPAUSER, Unpauser2.address))
              .to.emit(PauseController, "RoleGranted").withArgs(UNPAUSER, Unpauser2.address, Admin1.address);
            expect(await PauseController.connect(Admin1).revokeRole(UNPAUSER, Unpauser2.address))
              .to.emit(PauseController, "RoleRevoked").withArgs(UNPAUSER, Unpauser2.address, Admin1.address);
        });
        it("should allow an admin account to add/revoke pauser adders with grantRole/revokeRole", async () => {
            expect(await PauseController.connect(Admin1).grantRole(PAUSER_ADDER, PauserAdder2.address))
              .to.emit(PauseController, "RoleGranted").withArgs(PAUSER_ADDER, PauserAdder2.address, Admin1.address);
            expect(await PauseController.connect(Admin1).revokeRole(PAUSER_ADDER, PauserAdder2.address))
              .to.emit(PauseController, "RoleRevoked").withArgs(PAUSER_ADDER, PauserAdder2.address, Admin1.address);
        });
        it("should allow an admin account to add/revoke canceler adders with grantRole/revokeRole", async () => {
            expect(await PauseController.connect(Admin1).grantRole(CANCELER_ADDER, CancelerAdder2.address))
              .to.emit(PauseController, "RoleGranted").withArgs(CANCELER_ADDER, CancelerAdder2.address, Admin1.address);
            expect(await PauseController.connect(Admin1).revokeRole(CANCELER_ADDER, CancelerAdder2.address))
              .to.emit(PauseController, "RoleRevoked").withArgs(CANCELER_ADDER, CancelerAdder2.address, Admin1.address);
        });
        it("should allow an admin account to add/remove unpauser admins with grantRole/revokeRole", async () => {
            expect(await PauseController.connect(Admin1).grantRole(UNPAUSER_ADMIN, UnpauserAdmin2.address))
              .to.emit(PauseController, "RoleGranted").withArgs(UNPAUSER_ADMIN, UnpauserAdmin2.address, Admin1.address);
            expect(await PauseController.connect(Admin1).revokeRole(UNPAUSER_ADMIN, UnpauserAdmin2.address))
              .to.emit(PauseController, "RoleRevoked").withArgs(UNPAUSER_ADMIN, UnpauserAdmin2.address, Admin1.address);
        });
        it("should allow an admin account to add/removke default admin role with grantRole/revokeRole", async () => {
            expect(await PauseController.connect(Admin1).grantRole(ADMIN, Admin2.address))
              .to.emit(PauseController, "RoleGranted").withArgs(ADMIN, Admin2.address, Admin1.address);
            expect(await PauseController.connect(Admin1).revokeRole(ADMIN, Admin2.address))
              .to.emit(PauseController, "RoleRevoked").withArgs(ADMIN, Admin2.address, Admin1.address);
        });
        it("smart contract deployer should have no special admin privileges", async () => {
            // Deployer should not be able to add any roles
            await expect(PauseController.connect(Deployer).grantRole(PAUSER, Pauser2.address)).to.be.reverted;
            await expect(PauseController.connect(Deployer).grantRole(PAUSER_ADDER, PauserAdder2.address)).to.be.reverted;
            await expect(PauseController.connect(Deployer).grantRole(CANCELER, Canceler2.address)).to.be.reverted;
            await expect(PauseController.connect(Deployer).grantRole(CANCELER_ADDER, CancelerAdder2.address)).to.be.reverted;
            await expect(PauseController.connect(Deployer).grantRole(UNPAUSER, Unpauser2.address)).to.be.reverted;
            await expect(PauseController.connect(Deployer).grantRole(UNPAUSER_ADMIN, UnpauserAdmin2.address)).to.be.reverted;
            await expect(PauseController.connect(Deployer).grantRole(ADMIN, Admin2.address)).to.be.reverted;
            // Deployer should not be able to revoke any roles
            await expect(PauseController.connect(Deployer).revokeRole(PAUSER, Pauser1.address)).to.be.reverted;
            await expect(PauseController.connect(Deployer).revokeRole(PAUSER_ADDER, PauserAdder1.address)).to.be.reverted;
            await expect(PauseController.connect(Deployer).revokeRole(CANCELER, Canceler1.address)).to.be.reverted;
            await expect(PauseController.connect(Deployer).revokeRole(CANCELER_ADDER, CancelerAdder1.address)).to.be.reverted;
            await expect(PauseController.connect(Deployer).revokeRole(UNPAUSER, Unpauser1.address)).to.be.reverted;
            await expect(PauseController.connect(Deployer).revokeRole(UNPAUSER_ADMIN, UnpauserAdmin1.address)).to.be.reverted;
            await expect(PauseController.connect(Deployer).revokeRole(ADMIN, Admin1.address)).to.be.reverted;
        });
        it("should allow an unpause admin to add or remove unpausers", async () => {
            expect(await PauseController.connect(UnpauserAdmin1).grantRole(UNPAUSER, Unpauser2.address))
              .to.emit(PauseController, "RoleGranted").withArgs(UNPAUSER, Unpauser2.address, UnpauserAdmin1.address);
            expect(await PauseController.connect(UnpauserAdmin1).revokeRole(UNPAUSER, Unpauser2.address))
              .to.emit(PauseController, "RoleRevoked").withArgs(UNPAUSER, Unpauser2.address, UnpauserAdmin1.address);
        });
        it("should not allow unpause admins to add or remove other unpause admins", async () => {
            await expect(PauseController.connect(UnpauserAdmin1).grantRole(UNPAUSER_ADMIN, UnpauserAdmin2.address)).to.be.reverted;
            await expect(PauseController.connect(UnpauserAdmin1).revokeRole(UNPAUSER_ADMIN, UnpauserAdmin1.address)).to.be.reverted;
        });
        it("should not allow an unpauser to add or remove other unpausers or unpauser admins", async () => {
            await expect(PauseController.connect(Unpauser1).grantRole(UNPAUSER, Unpauser2.address)).to.be.reverted;
            await expect(PauseController.connect(Unpauser1).revokeRole(UNPAUSER, Unpauser1.address)).to.be.reverted;
            await expect(PauseController.connect(Unpauser1).grantRole(UNPAUSER_ADMIN, UnpauserAdmin2.address)).to.be.reverted;
            await expect(PauseController.connect(Unpauser1).revokeRole(UNPAUSER_ADMIN, UnpauserAdmin1.address)).to.be.reverted;
        });
    });

    describe("Bridge Pausing", () => {
        function stringToBytes(message: string): BytesLike {
            return ethers.utils.hexlify(ethers.utils.toUtf8Bytes(message));
        }
        it("should allow a user with the pauser role to pause the bridge with an empty message", async () => {
            await expect(PauseController.connect(Pauser1).pause([]))
              .to.emit(PauseController, "Pause").withArgs(Pauser1.address, false, []);
        });
        it("should allow a user with the pauser role to pause the bridge with a message", async () => {
            const message = stringToBytes("Bridgebank under maintenance");
            await expect(PauseController.connect(Pauser1).pause(message))    
              .to.emit(PauseController, "Pause").withArgs(Pauser1.address, false, message);
        });
        it("should allow a user to update the pauser message if the first one was blank", async () => {
            await expect(PauseController.connect(Pauser1).pause([]))
              .to.emit(PauseController, "Pause").withArgs(Pauser1.address, false, []);
            const message = stringToBytes("Bridgebank under maintenance");
            await expect(PauseController.connect(Pauser1).pause(message))    
              .to.emit(PauseController, "Pause").withArgs(Pauser1.address, true, message);
        });
        it("should allow a second user to update the pauser message if the first one was blank", async () => {
            await PauseController.connect(PauserAdder1).addPauser(Pauser2.address);
            await expect(PauseController.connect(Pauser1).pause([]))
              .to.emit(PauseController, "Pause").withArgs(Pauser1.address, false, []);
            const message = stringToBytes("Bridgebank under maintenance");
            await expect(PauseController.connect(Pauser2).pause(message))    
              .to.emit(PauseController, "Pause").withArgs(Pauser2.address, true, message);
        });
        it("should allow a very large pause message", async () => {
            const loremIpsum = `
            Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. 
            Erat imperdiet sed euismod nisi. Vitae justo eget magna fermentum iaculis eu non. Suspendisse faucibus interdum posuere lorem 
            ipsum dolor sit amet consectetur. Cursus euismod quis viverra nibh cras pulvinar mattis. Dignissim convallis aenean et tortor. 
            In ante metus dictum at tempor. Eget lorem dolor sed viverra ipsum nunc aliquet. Vel facilisis volutpat est velit egestas dui id. 
            Aliquet nec ullamcorper sit amet risus nullam eget felis eget. Nulla aliquet enim tortor at auctor urna nunc.

            Sit amet consectetur adipiscing elit ut aliquam purus. Lacus viverra vitae congue eu consequat ac felis. Nullam vehicula ipsum 
            a arcu cursus. Feugiat in fermentum posuere urna nec tincidunt praesent semper feugiat. Vulputate sapien nec sagittis aliquam 
            malesuada bibendum. Convallis convallis tellus id interdum velit laoreet id donec. Lectus arcu bibendum at varius vel. Sit amet 
            porttitor eget dolor morbi non arcu risus quis. Facilisis mauris sit amet massa vitae tortor. Nec dui nunc mattis enim ut tellus. 
            Suspendisse in est ante in nibh mauris cursus. Arcu felis bibendum ut tristique et egestas quis ipsum. Cursus in hac habitasse platea. 
            Pellentesque nec nam aliquam sem et tortor consequat. Sit amet volutpat consequat mauris nunc congue nisi vitae suscipit. 
            Id eu nisl nunc mi ipsum. Nec ultrices dui sapien eget mi proin sed. Convallis aenean et tortor at risus.
            `
            const message = stringToBytes(loremIpsum);
            await expect(PauseController.connect(Pauser1).pause(message))    
              .to.emit(PauseController, "Pause").withArgs(Pauser1.address, false, message);
        });
        it("should not allow a user without the pauser role to pause the bridge", async () => {
            await expect(PauseController.connect(Deployer).pause([]))
              .to.be.revertedWith("User is not pauser");
        });
    });
    
    describe("Request Unpause", () => {
        it("should revert if it is called when the bridgebank is not already paused", async () => {
            await expect(PauseController.connect(Unpauser1).requestUnpause())
              .to.be.revertedWith("BridgeBank not paused");
        });
        it("should allow an unpauser to request an unpause if the bridgebank is paused", async () => {
            await PauseController.connect(Pauser1).pause([]);
            const tx = await PauseController.connect(Unpauser1).requestUnpause()
            expect(tx).to.emit(PauseController, "UnpauseRequest")
              .withArgs(Unpauser1.address, Number(tx.blockNumber) + TimeLockDelay);
        });
        it("should revert if the requester does not have the unpauser role", async () => {
            await PauseController.connect(Pauser1).pause([]);
            await expect(PauseController.connect(Pauser1).requestUnpause())
              .to.be.revertedWith("User is not unpauser");
        });
        it("should revert if a curent unpause request is pending", async () => {
            await PauseController.connect(Pauser1).pause([]);
            await PauseController.connect(Unpauser1).requestUnpause();
            await expect(PauseController.connect(Unpauser1).requestUnpause())
              .to.be.revertedWith("Unpause request already pending");
        });
    });
    
    describe("Cancel Unpause", () => {
        beforeEach(async () => {
            await PauseController.connect(Pauser1).pause([]);
            await PauseController.connect(Unpauser1).requestUnpause()
        });
        it("should allow a unpause request to be canceled when called by a user with the canceler role", async () => {
            expect(await PauseController.UnpauseRequestBlockHeight()).to.not.equal(NO_REQUEST)
            expect(await PauseController.connect(Canceler1).cancelUnpause())
              .to.emit(PauseController, "CancelUnpause").withArgs(Canceler1.address);
            expect(await PauseController.UnpauseRequestBlockHeight()).to.equal(NO_REQUEST);
        });
        it("should allow a unpause request to be repeatedly called by a user with the canceler role", async () => {
            expect(await PauseController.UnpauseRequestBlockHeight()).to.not.equal(NO_REQUEST)
            expect(await PauseController.connect(Canceler1).cancelUnpause())
              .to.emit(PauseController, "CancelUnpause").withArgs(Canceler1.address);
            expect(await PauseController.UnpauseRequestBlockHeight()).to.equal(NO_REQUEST);
            expect(await PauseController.connect(Canceler1).cancelUnpause())
              .to.emit(PauseController, "CancelUnpause").withArgs(Canceler1.address);
            expect(await PauseController.UnpauseRequestBlockHeight()).to.equal(NO_REQUEST);
        });
        it("should not allow a user without the canceler role to call cancelUnpause", async () => {
            await expect(PauseController.connect(Deployer).cancelUnpause())
              .to.be.revertedWith("User is not canceler");
        });
    });

    describe("Unpause", () => {
        const provider = hre.network.provider;
        describe("Before an unpause request", () => {
            it("should revert with no active request", async () => {
               await expect(PauseController.connect(Deployer).unpause())
                 .to.be.revertedWith("No Active Unpause Request");
            });
        });
        describe("After Pause Request", () => {
            beforeEach(async () => {
                await PauseController.connect(Pauser1).pause([]);
                await PauseController.connect(Unpauser1).requestUnpause()
            });
            it("should revert if timelock is still in effect", async() => {
                await expect(PauseController.connect(Deployer).unpause())
                  .to.be.revertedWith("TimeLock still in effect");
            });
            it("should revert a single block before the TimeLockDelay has passed", async () => {
                const currentBlock = await PauseController.provider.getBlockNumber()
                const requestBlock = await PauseController.UnpauseRequestBlockHeight();
                const oneShort = requestBlock.sub(currentBlock + 1)
                for (let i=0; oneShort.gt(i); i++) {
                    provider.send("evm_mine");
                }
                await expect(PauseController.connect(Deployer).unpause())
                .to.be.revertedWith("TimeLock still in effect");
            });
            it("should allow any user to call if a single block has passed the TimeLockDelay", async () => {
                const currentBlock = await PauseController.provider.getBlockNumber()
                const requestBlock = await PauseController.UnpauseRequestBlockHeight();
                const exactDelay = requestBlock.sub(currentBlock)
                for (let i=0; exactDelay.gt(i); i++) {
                    provider.send("evm_mine");
                }
                await expect(PauseController.connect(Deployer).unpause())
                .to.emit(PauseController, "Unpause").withArgs(Deployer.address);
            });
            it("should allow any user to call if a many blocks have passed the TimeLockDelay", async () => {
                const currentBlock = await PauseController.provider.getBlockNumber()
                const requestBlock = await PauseController.UnpauseRequestBlockHeight();
                const manyOver = requestBlock.add(currentBlock)
                for (let i=0; manyOver.gt(i); i++) {
                    provider.send("evm_mine");
                }
                await expect(PauseController.connect(Deployer).unpause())
                .to.emit(PauseController, "Unpause").withArgs(Deployer.address);
            });
            it("should allow any user to call even if the bridge has already been unpaused", async () => {
                const currentBlock = await PauseController.provider.getBlockNumber()
                const requestBlock = await PauseController.UnpauseRequestBlockHeight();
                const exactDelay = requestBlock.sub(currentBlock)
                for (let i=0; exactDelay.gt(i); i++) {
                    provider.send("evm_mine");
                }
                await BridgeBank.connect(Deployer).unpause();
                await expect(PauseController.connect(Deployer).unpause())
                .to.emit(PauseController, "Unpause").withArgs(Deployer.address);
            });
        });
    });
});