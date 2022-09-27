// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.8.17;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

/**
 * @title Commission Token
 * @dev This token will charge a programmable fee on every transfer that is credited to the the devolpoers address
 *      this test token will help test various tokens that may transfer a different amount then the transfer requests
 *      to the recipient.
 **/
contract CommissionToken is ERC20 {
    address private dev;
    uint256 public transferFee;
    /**
     * @dev This token needs a dev account to be credited with the dev fee, a initial user account for the minting, and a 
     *      one time quantity to mint. 
     * @param _dev The address of the developer that gets paid the devFee
     * @param _devFee The fee as a thousandth of a percent charged per transfer. (e.g. 500 would be 5% transfer fee)
     * @param _user The address to mint the initial tokens to, this is a fixed supply token
     * @param _quantity The quantity to mint, this is a fixed supply token
     */
    constructor(address _dev, uint256 _devFee, address _user, uint256 _quantity) ERC20("Commission Token", "CMT") {
        require(_dev != address(0), "Dev account must not be null address");
        require(_devFee < 10_000, "Dev Fee cannot exceed 100%");
        require(_devFee > 0, "Dev Fee cannot be 0%");
        require(_user != address(0), "Initial minting address must not be null address");
        dev = _dev;
        transferFee = _devFee;
        _mint(_user, _quantity);
    }

    function _transfer(address sender, address recipient, uint256 amount) internal override {
        uint256 devFee = amount / 10_000;
        devFee *= transferFee;
        uint256 transferAmount = amount - devFee;

        // Send dev fee to dev address
        super._transfer(sender, dev, devFee);
        // Send remainder to intended recipient
        super._transfer(sender, recipient, transferAmount);
    }
}