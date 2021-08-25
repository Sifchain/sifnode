package txs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// ParseDecimalFile parse the file to a map
func ParseDecimalFile(fileName string) map[string]int {
	decimalFile, err := ioutil.ReadFile(fileName)

	if err != nil {
		fmt.Printf("error as %s\n", err.Error())
		return map[string]int{}
	}

	var data map[string]int
	err = json.Unmarshal(decimalFile, &data)
	if err != nil {
		fmt.Printf("error as %s\n", err.Error())
		return map[string]int{}
	}

	return data
}
