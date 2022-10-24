pragma solidity 0.5.16;

contract ToLower {

    function toLower(string memory str) public pure returns (string memory) {
		bytes memory bStr = bytes(str);
		bytes memory bLower = new bytes(bStr.length);
		for (uint i = 0; i < bStr.length; i++) {
			// Uppercase character...
			if ((bStr[i] >= bytes1(uint8(65))) && (bStr[i] <= bytes1(uint8(90)))) {
				// So we add 32 to make it lowercase
				bLower[i] = bytes1(uint8(bStr[i]) + 32);
			} else {
				bLower[i] = bStr[i];
			}
		}
		return string(bLower);
	}
}