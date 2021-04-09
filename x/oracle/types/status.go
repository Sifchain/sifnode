package types

import (
	"encoding/json"
	"fmt"
	"strconv"
)

var StatusTextToString = map[StatusText]string{
	StatusText_STATUS_TEXT_UNSPECIFIED: "unspecified",
	StatusText_STATUS_TEXT_PENDING: "pending",
	StatusText_STATUS_TEXT_SUCCESS: "success",
	StatusText_STATUS_TEXT_FAILED: "failed",
}
var StringToStatusText = map[string]StatusText{
	"pending": StatusText_STATUS_TEXT_PENDING,
	"success": StatusText_STATUS_TEXT_SUCCESS,
	"failed":  StatusText_STATUS_TEXT_FAILED,
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
