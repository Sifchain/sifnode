package types

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// ClaimType is an enum used to represent the type of claim
// type ClaimType int

var ClaimTypeToString = map[ClaimType]string{
	ClaimType_CLAIM_TYPE_UNSPECIFIED: "unspecified",
	ClaimType_CLAIM_TYPE_LOCK: "lock",
	ClaimType_CLAIM_TYPE_BURN: "burn",
}

func StringToClaimType(text string) (ClaimType, error) {
	switch text {
	case "lock":
		return ClaimType_CLAIM_TYPE_LOCK, nil
	case "burn":
		return ClaimType_CLAIM_TYPE_BURN, nil
	default:
		return 0, ErrInvalidClaimType
	}
}

func SerializeClaimType(claimType ClaimType) string {
	return ClaimTypeToString[claimType]
}

func (ct ClaimType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", SerializeClaimType(ct))), nil
}

func (text *ClaimType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	stringKey, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}

	value, err := StringToClaimType(stringKey)
	if err != nil {
		return err
	}
	*text = value
	return nil
}
