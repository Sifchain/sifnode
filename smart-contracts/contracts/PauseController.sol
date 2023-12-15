pragma solidity 0.8.17;

import "@openzeppelin/contracts/access/AccessControlEnumerable.sol";
import "./interfaces/IPausable.sol";

/**
 * @title Pause Controller
 * @dev The PauseController is a contract which should be given the pauser role of BridgeBank.
 *      This contract will delegate roles for pausing, unpausing, and adding new users; it applies
 *      a programmable timelock for unpausing. All admin and unpauser roles should be multisig
 *      contracts with hardware wallets for the signers. All Pauser roles should be EOA addresses
 *      with hardware wallets.
 */
contract PauseController is AccessControlEnumerable {
    /**
     * @dev Constant Roles for Access Control
     */
    bytes32 public constant PAUSER = keccak256("PAUSER"); // User who can pause the bridge
    bytes32 public constant PAUSER_ADDER = keccak256("PAUSER_ADDER"); // User who can add pauser accounts
    bytes32 public constant CANCELER = keccak256("CANCELER"); // User who can cancel unpause requests
    bytes32 public constant CANCELER_ADDER = keccak256("CANCELER_ADDER"); // User who can add canceler accounts
    bytes32 public constant UNPAUSER = keccak256("UNPAUSER"); // User who can unpause the bridge
    bytes32 public constant UNPAUSER_ADMIN = keccak256("UNPAUSER_ADMIN"); // User who can add and remove unpauser accounts

    // Public Variables

    /**
     * @dev Constant value for No Unpause Request
     */
    uint256 public constant NOREQUEST = 1;

    /**
     * @dev How long you must wait to resume the bridge
     */
    uint256 public immutable TimeLockDelay;

    /**
     * @dev Address of the BridgeBank
     */
    Pausable public immutable BridgeBank;

    /**
     * @dev Pending UnPause Request Block Height
     */
    uint256 public UnpauseRequestBlockHeight;

    // Events 

    /**
     * @dev Event Emitted when a pause transaction is successfully submitted
     */
    event Pause(
        address indexed _pauser, // Account which paused the bridge
        bool _messageUpdate, // Is this just a message update (true) or an actual pause event (false)
        bytes message // Message sent for why the pause occurred
    );

    /**
     * @dev Event Emitted when a unpause request is successfully submitted
     */
    event UnpauseRequest(
        address indexed _requester,
        uint256 indexed _UnpauseRequestBlockHeight
    );

    /**
     * @dev Event Emitted when a Cancel Unpause transaction is successfully submitted
     */
    event CancelUnpause(address indexed _canceler);

    /**
     * @dev Event Emitted when a Unpause transaction is successfully submitted
     */
    event Unpause(address indexed _unpauser);

    // Possible Errors

    /**
     * @dev Raised when the constructor submits the Zero (NULL) address as the BridgeBank contract address
     */
    error BridgeBankAddressIsNull();

    /**
     * @dev Raised when the BridgeBank is not currently paused but an unpause request was received.
     */
    error BridgeBankNotPaused();

    /**
     * @dev Raised when unpause request is called with an already pending unpause request.
     */
    error UnpauseRequestAlreadyPending();

    /**
     * @dev Raised when calling unpause without an unpause request having been submitted first.
     * You must first submit an unpauseRequest and wait the required timelock period before
     * submitting the unpause call.
     */
    error NoActiveUnpauseRequest();

    /**
     * @dev Raised when unpause is called before the time lock has passed. User must wait the remaining
     * time before calling the function.
     * @param remainingBlocks The number of blocks to wait before the function can be called
     */
    error TimeLock(uint256 remainingBlocks);

    /**
     * @dev User does not have the Pauser role which is required for this function call.
     */
    error UserIsNotPauser();

    /**
     * @dev User does not have the PauserAdder role which is required for this function call.
     */
    error UserIsNotPauserAdder();

    /**
     * @dev User does not have the Canceler role which is required for this function call.
     */
    error UserIsNotCanceler();

    /**
     * @dev User does not have the CancelerAdder role which is required for this function call.
     */
    error UserIsNotCancelerAdder();

    /**
     * @dev User does not have the Unpauser role which is required for this function call.
     */
    error UserIsNotUnpauser();

    // Functions

    /**
     * @dev On contract construction the bridgebank address and timelock delay must be set. These values
     *      are immutable, you must create a new pausecontroller contract if you want to change those values.
     *      Initial roles can be set on creation as well however so long as a single admin role is created they
     *      can be changed at later times.
     * @param _bridgeBank The address of the bridgebank contract that this pause controller will pause/unpause
     * @param _timelockDelay How many blocks to wait after an unpause command before the unpause call can be made
     * @param _admins An array of addresses to give the admin role which can add/remove any users from any roles (Including Admin)
     *                this is a super user, use with care.
     * @param _pausers An array of addresses which can pause the bridgebank contracts
     * @param _pauser_adder An array of addresses which can give the pauser role
     * @param _unpausers An array of addresses which can schedule an unpause and then execute an unpause call.
     * @param _unpauser_admin An array of addresses which can give and revoke the unpauser role
     * @param _cancelers An array of addresses which can cancel a scheduled unpause if a compromise is suspected
     * @param _canceler_adder An array of addresses which can give the canceler role
     */
    constructor(
        address _bridgeBank,
        uint256 _timelockDelay,
        address[] memory _admins,
        address[] memory _pausers,
        address[] memory _pauser_adder,
        address[] memory _unpausers,
        address[] memory _unpauser_admin,
        address[] memory _cancelers,
        address[] memory _canceler_adder
    ) {
        // Set Unpauser_admin to the admin role of unpausers
        _setRoleAdmin(UNPAUSER, UNPAUSER_ADMIN);
        // Populate each role, These will be modifiable later by the admin role
        uint256 length = _admins.length;
        for (uint256 i; i < length; ) {
            _grantRole(DEFAULT_ADMIN_ROLE, _admins[i]);
            _grantRole(PAUSER_ADDER, _admins[i]);
            _grantRole(UNPAUSER_ADMIN, _admins[i]);
            _grantRole(CANCELER_ADDER, _admins[i]);
            unchecked {
                ++i;
            }
        }

        length = _pausers.length;
        for (uint256 i; i < length; ) {
            _grantRole(PAUSER, _pausers[i]);
            unchecked {
                ++i;
            }
        }

        length = _pauser_adder.length;
        for (uint256 i; i < length; ) {
            _grantRole(PAUSER_ADDER, _pauser_adder[i]);
            unchecked {
                ++i;
            }
        }

        length = _unpausers.length;
        for (uint256 i; i < length; ) {
            _grantRole(UNPAUSER, _unpausers[i]);
            unchecked {
                ++i;
            }
        }

        length = _unpauser_admin.length;
        for (uint256 i; i < length; ) {
            _grantRole(UNPAUSER_ADMIN, _unpauser_admin[i]);
            unchecked {
                ++i;
            }
        }

        length = _cancelers.length;
        for (uint256 i; i < length; ) {
            _grantRole(CANCELER, _cancelers[i]);
            unchecked {
                ++i;
            }
        }

        length = _canceler_adder.length;
        for (uint256 i; i < length; ) {
            _grantRole(CANCELER_ADDER, _canceler_adder[i]);
            unchecked {
                ++i;
            }
        }

        // Log the bridgebank address, This will be immutable
        if (_bridgeBank == address(0)) {
            revert BridgeBankAddressIsNull();
        }

        BridgeBank = Pausable(_bridgeBank);

        // Set the TimeLockDelay, This will be immutable
        TimeLockDelay = _timelockDelay;

        // Set the UnPauseRequestBlockHeight to a default value of 1 for gas savings
        UnpauseRequestBlockHeight = NOREQUEST;
    }

    /**
     * @dev Anyone with the pauser role will be able to pause the bridgebank
     *      assuming the bridge is not already paused. This operation should be
     *      considered a lower privilege process that trusted devs will have access
     *      to quickly stop the bridge if a problem is detected.
     *
     *      Only require the use of hardware wallets.
     * @param message A bytes message that can be submitted with the pause request such that 
     * UI's can display the reason why the bridge was paused in a banner. This message is 
     * emitted with the pause log event.
     */
    function pause(bytes calldata message) external {
        address pauser = msg.sender;
        if (!hasRole(PAUSER, pauser)) {
            revert UserIsNotPauser();
        }
        bool paused = BridgeBank.paused();
        if (!paused) {
            BridgeBank.pause();
        }
        emit Pause(pauser, paused, message);
    }

    /**
     * @dev Anyone with the unpauser role will be able to request the bridgebank be resumed. 
     *      This is a timelocked function based upon the block delay set at construction.
     *      This should be considered a highly privileged function, require the use of
     *      multisig contracts as well as hardware wallets for the signers.
     */
    function requestUnpause() external {
        address requester = msg.sender;
        if (!hasRole(UNPAUSER, requester)) {
            revert UserIsNotUnpauser();
        }
        if (UnpauseRequestBlockHeight != NOREQUEST) {
            revert UnpauseRequestAlreadyPending();
        }
        bool paused = BridgeBank.paused();
        if (!paused) {
            revert BridgeBankNotPaused();
        }
        uint256 RequestBlockHeight;
        unchecked {
            RequestBlockHeight = block.number + TimeLockDelay;
        }
        UnpauseRequestBlockHeight = RequestBlockHeight;
        emit UnpauseRequest(requester, RequestBlockHeight);
    }

    /**
     * @dev If a problem is discovered after an unpause is scheduled a user with the canceler role
     *      will be able to cancel a scheduled pause. This may be relevant if a problem still exists
     *      after attempting to resume OR if a unpauser account is suspected of being compromised.
     *      In the event of an unpauser compromise an admin should revoke the unpauser account controls
     *      and cancelers should cancel any unpause attempts until revocation is complete.
     *
     *      Only require the use of hardware wallets.
     */
    function cancelUnpause() external {
        address requester = msg.sender;
        if (!hasRole(CANCELER, requester)) {
            revert UserIsNotCanceler();
        }
        UnpauseRequestBlockHeight = NOREQUEST;
        emit CancelUnpause(requester);
    }

    /**
     * @dev External viewable function that reports how many blocks are remaining on the time locked period.
     * @return Returns a uint256 of the remaining blocks before the timelock has passed. If 0 there is no remaining blocks
     * or there is no active time locked request. 
     */
    function TimeLockRemaining() external view returns (uint256) {
        unchecked{
            uint256 timeLock = UnpauseRequestBlockHeight;
            uint256 currentBlock = block.number;
            if (currentBlock > timeLock) {
                return 0;
            }
            return timeLock - currentBlock;
        }
    }

    /**
     * @dev If a unpause request has been submitted by a person with the unpauser role, and the required
     *      block delay has been waited anyone will be able to complete the unpause request resuming the
     *      bridgebank contract.
     *
     *      Unprivileged function, since it has to be scheduled by an unpauser anyone can complete this execution
     */
    function unpause() external {
        uint256 requestBlockHeight = UnpauseRequestBlockHeight;
        if (requestBlockHeight == NOREQUEST) {
            revert NoActiveUnpauseRequest();
        }
        if (requestBlockHeight > block.number) {
            unchecked {
                revert TimeLock(requestBlockHeight - block.number);
            }
        }
        bool paused = BridgeBank.paused();
        if (paused) {
            BridgeBank.unpause();
        }
        UnpauseRequestBlockHeight = NOREQUEST;
        emit Unpause(msg.sender);
    }

    /**
     * @dev Function to add new pausers to the system. Must be called by an account
     *      with the pauser_adder role. This function is used over the standard
     *      grantRole function for admin roles because we want these users only to
     *      be able to grant the roles and not be able to remove pausers.
     * @param account The account address to grant the pauser role to
     */
    function addPauser(address account) external {
        if (!hasRole(PAUSER_ADDER, msg.sender)) {
            revert UserIsNotPauserAdder();
        }
        _grantRole(PAUSER, account);
    }

    /**
     * @dev Function to add new canceler to the system. Must be called by an account
     *      with the canceler_adder role. This function is used over the standard
     *      grantRole function for admin roles because we want these users only to
     *      be able to grant the roles and not be able to remove pausers.
     * @param account The account address to grant the canceler role to
     */
    function addCanceler(address account) external {
        if (!hasRole(CANCELER_ADDER, msg.sender)) {
            revert UserIsNotCancelerAdder();
        }
        _grantRole(CANCELER, account);
    }
}
