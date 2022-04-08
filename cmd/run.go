package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"time"
)

type Output struct {
	AddressSender   string `json:"address_sender"`
	AddressReceiver string `json:"address_receiver"`
	Amount          string `json:"amount"`
}

type TempOutput struct {
	Out []Output `json:"Output"`
}

func ParseOutput(fp string) ([]Output, error) {
	var outputs TempOutput
	file, err := filepath.Abs(fp)
	if err != nil {
		return nil, err
	}
	o, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(o, &outputs)
	if err != nil {
		return nil, err
	}
	return outputs.Out, nil
}

func main() {
	outputs, err := ParseOutput("cmd/address-list.json")
	if err != nil {
		panic(err)
	}
	fmt.Println(outputs)
	for _, output := range outputs {
		cmd := exec.Command(
			"blogd",
			"tx",
			"bank",
			"send",
			output.AddressSender,
			output.AddressReceiver,
			output.Amount,
			"--keyring-backend=test",
			"--chain-id=localnet",
			"--fees=100000000000000000000000MYTOKEN",
			"--yes")
		var out bytes.Buffer
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			panic(err)
		}
		fmt.Println(out.String())
		time.Sleep(time.Second * 6)
	}
}
