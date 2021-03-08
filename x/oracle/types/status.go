package types

import (
	"encoding/json"
	"fmt"
	"strconv"
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
