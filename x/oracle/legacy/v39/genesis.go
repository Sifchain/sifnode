package v39

import (
	"encoding/json"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const ModuleName = "oracle"

type GenesisState struct {
	AddressWhitelist []sdk.ValAddress `json:"address_whitelist"`
	AdminAddress     sdk.AccAddress   `json:"admin_address"`
	Prophecies       []DBProphecy     `json:"prophecies"`
}

type Prophecy struct {
	ID     string `json:"id"`
	Status Status `json:"status"`

	//WARNING: Mappings are nondeterministic in Amino,
	// an so iterating over them could result in consensus failure. New code should not iterate over the below 2 mappings.

	//This is a mapping from a claim to the list of validators that made that claim.
	ClaimValidators map[string][]sdk.ValAddress `json:"claim_validators"`
	//This is a mapping from a validator bech32 address to their claim
	ValidatorClaims map[string]string `json:"validator_claims"`
}

// DBProphecy is what the prophecy becomes when being saved to the database.
//  Tendermint/Amino does not support maps so we must serialize those variables into bytes.
type DBProphecy struct {
	ID              string `json:"id"`
	Status          Status `json:"status"`
	ClaimValidators []byte `json:"claim_validators"`
	ValidatorClaims []byte `json:"validator_claims"`
}

type Status struct {
	Text       StatusText `json:"text"`
	FinalClaim string     `json:"final_claim"`
}

type StatusText int

const (
	PendingStatusText StatusText = iota
	SuccessStatusText
	FailedStatusText
)

var StatusTextToString = [...]string{"pending", "success", "failed"}
var StringToStatusText = map[string]StatusText{
	"pending": PendingStatusText,
	"success": SuccessStatusText,
	"failed":  FailedStatusText,
}

func (text StatusText) String() string {
	return StatusTextToString[text]
}

func (text StatusText) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", text.String())), nil
}

func (text *StatusText) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	stringKey, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}

	// Note that if the string cannot be found then it will be set to the zero value, 'pending' in this case.
	*text = StringToStatusText[stringKey]
	return nil
}
