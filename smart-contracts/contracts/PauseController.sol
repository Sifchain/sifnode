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
    bytes32 public constant CANCELER = keccak256("CANCELER"); // User who can cancel unpause requests
    bytes32 public constant UNPAUSER = keccak256("UNPAUSER"); // User who can unpause the bridge
    // AccessControlEnumerable has a built in DEFAULT_ADMIN_ROLE which are users that can add/remove user roles
    
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

    /**
     * @dev Event Emitted when a pause transaction is successfully submitted
     */
     event Pause(
        address indexed _pauser
     );

    /**
     * @dev Event Emitted when a unpause request is successfully submitted
     */
     event UnpauseRequest(
        address indexed _requester,
        uint256 _UnpauseRequestBlockHeight
     );

     /**
      * @dev Event Emitted when a Cancel Unpause transaction is succcessfully submitted
      */
      event CancelUnpause(
        address indexed _canceler
     );

    /**
     * @dev Event Emitted when a Unpause transaction is successfully submitted
     */
    event Unpause(
        address indexed _unpauser
    );
    
    /**
     * @dev On contract construction the bridgebank address and timelock delay must be set. These values 
     *      are immutable, you must create a new pausecontroller contract if you want to change those values.
     *      Initial roles can be set on creation as well however so long as a single admin role is created they
     *      can be changed at later times.
     * @param _bridgeBank The address of the bridgebank contract that this pause controller will pause/unpause
     * @param _timelockDelay How many blocks to wait after an unpause command before the unpause call can be made
     * @param _admins An array of addresses to give the admin role which can add/remove any users from roles (Including Admin)
     *                this is a super user, use with care.
     * @param _pausers An array of addresses which can pause the bridgebank contracts
     * @param _unpausers An array of addresses which can schedule an unpause and then execute an unpause call.
     * @param _cancelers An array of addresses which can cancel a scheduled unpause if a compromise is suspected
     */
    constructor(
        address _bridgeBank, 
        uint256 _timelockDelay, 
        address [] memory _admins, 
        address [] memory _pausers, 
        address [] memory _unpausers,
        address [] memory _cancelers
    ) {
        // Populate each role, These will be modifiable later by the admin role
        uint256 length = _admins.length;
        for (uint256 i; i<length;) {
            /**
             * @Note _setupRole has been deprecated in favor of _grantRole however 
             *       I can not use _grantRole until we upgrade versions of OpenZepplin
             *       which should not happen until after the peggy2/master merge. 
             * 
             *       AUDITORS: Please evaluate this code both as its written with _setupRole
             *                 and evaluate this code if we change it to _grantRole.
             */
            _setupRole(DEFAULT_ADMIN_ROLE, _admins[i]);
            unchecked { ++i; }
        }

        length = _pausers.length;
        for (uint256 i; i<length;) {
            // See note in earlier loop
            _setupRole(PAUSER, _pausers[i]);
            unchecked { ++i; }
        }

        length = _unpausers.length;
        for (uint256 i; i<length;) {
            // See note in earlier loop
            _setupRole(UNPAUSER, _unpausers[i]);
            unchecked { ++i; }
        }

        length = _cancelers.length;
        for (uint256 i; i<length;) {
            // See note in earlier loop
            _setupRole(CANCELER, _cancelers[i]);
            unchecked { ++i; }
        }

        // Log the bridgebank address, This will be immutable
        require(_bridgeBank != address(0), "BridgeBank address must be set");
        BridgeBank = Pausable(_bridgeBank);

        // Set the TimeLockDelay, This will be immutable
        TimeLockDelay = _timelockDelay;

        // Set the UnPauseRequestBlockHeight to a default value of 1 for gas savings
        UnpauseRequestBlockHeight = NOREQUEST;
    }

    /**
     * @dev Anyone with the pauser role will be able to pause the bridgebank
     *      assuming the bridge is not already paused. This operation should be
     *      considered a lower privlage process that trusted devs will have access
     *      to quickly stop the bridge if a problem is detected. 
     * 
     *      Only require the use of hardware wallets.
     */
    function pause() public {
        address pauser = msg.sender;
        require(hasRole(PAUSER, pauser), "User is not pauser");
        bool paused = BridgeBank.paused();
        require(paused == false, "BridgeBank already paused");
        BridgeBank.pause();
        emit Pause(pauser);
    }

    /**
     * @dev Anyone with the unpauser role will be able to request the bridgebank be resumed.
     *      This is a timelocked function based upon the block delay set at construction. 
     *      This should be considered a highly privilaged function, require the use of 
     *      multisig contracts as well as hardware wallets for the signers.
     */
    function requestUnpause() public {
        address requester = msg.sender;
        require(hasRole(UNPAUSER, requester), "User is not unpauser");
        require(UnpauseRequestBlockHeight == NOREQUEST, "Unpause request already pending");
        bool paused = BridgeBank.paused();
        require(paused == true, "BridgeBank not paused");
        uint256 RequestBlockHeight = block.number + TimeLockDelay;
        UnpauseRequestBlockHeight = RequestBlockHeight;
        emit UnpauseRequest(requester, RequestBlockHeight);
    }

    /**
     * @dev If a problem is discovered after an unpause is scheduled a user with the canceler role
     *      will be able to cancel a scheduled pause. This may be relevent if a problem still exists 
     *      after attempting to resume OR if a unpauser account is suspected of being compromised. 
     *      In the event of an unpauser compromise an admin should revoke the unpauser account controls
     *      and cancelers should cancel any unpause attempts until revokation is complete.
     * 
     *      Only require the use of hardware wallets. 
     */
    function cancelUnpause() public {
        address requester = msg.sender;
        require(hasRole(CANCELER, requester), "User is not canceler");
        UnpauseRequestBlockHeight = NOREQUEST;
        emit CancelUnpause(requester);
    }

    /**
     * @dev If a unpause request has been submited by a person with the unpauser role, and the required
     *      block delay has been waited a user with the unpause role will be able to complete the unpause
     *      request resuming the bridgebank contract. Require use of multisig contracts and hardware wallets
     *      for the signers.
     */
    function unpause() public {
        address unpauser = msg.sender;
        require(hasRole(UNPAUSER, unpauser), "User is not unpauser");
        uint256 requestBlockHeight = UnpauseRequestBlockHeight;
        require(requestBlockHeight != NOREQUEST, "No Active Unpause Request");
        require(requestBlockHeight < block.number, "TimeLock still in effect");
        bool paused = BridgeBank.paused();
        if (paused == true) {
            BridgeBank.unpause();
        }
        UnpauseRequestBlockHeight = NOREQUEST;
        emit Unpause(unpauser);
    }
}