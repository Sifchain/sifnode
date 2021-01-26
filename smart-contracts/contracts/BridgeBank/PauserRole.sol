pragma solidity 0.5.16;

contract PauserRole {

    mapping (address => bool) public pausers;

    modifier onlyPauser() {
        require(pausers[msg.sender], "PauserRole: caller does not have the Pauser role");
        _;
    }

    function addPauser(address account) public onlyPauser {
        _addPauser(account);
    }

    function renouncePauser() public {
        _removePauser(msg.sender);
    }

    function _addPauser(address account) internal {
        pausers[account] = true;
    }

    function _removePauser(address account) internal {
        pausers[account] = false;
    }
}
