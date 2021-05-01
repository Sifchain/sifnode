package network

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Sifchain/sifnode/tools/sifgen/common"
	"github.com/Sifchain/sifnode/tools/sifgen/genesis"
	"github.com/Sifchain/sifnode/tools/sifgen/network/types"
	"github.com/Sifchain/sifnode/tools/sifgen/utils"

	"github.com/BurntSushi/toml"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"gopkg.in/yaml.v3"
)

type Network struct {
	ChainID    string
	CLI        utils.CLI
	BondAmount string
	MintAmount string
}

func Reset(chainID, networkDir string) error {
	s := sha256.New()
	_, err := s.Write([]byte(chainID))
	if err != nil {
		log.Println("Write failed: ", err.Error())
	}
	file := fmt.Sprintf("%x", s.Sum(nil))

	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.yml", networkDir, file))
	if err != nil {
		return err
	}

	var nodes types.Nodes
	if err := yaml.Unmarshal(data, &nodes); err != nil {
		return err
	}

	for _, node := range nodes {
		nodeDir := fmt.Sprintf("%s/validators/%s/%s/.sifnoded", networkDir, chainID, node.Moniker)
		_, err = utils.NewCLI(chainID, keyring.BackendTest).ResetState(nodeDir)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewNetwork(chainID string) *Network {
	return &Network{
		ChainID: chainID,
		CLI:     utils.NewCLI(chainID, keyring.BackendTest),
	}
}

func (n *Network) Build(count int, outputDir, seedIPv4Addr string) (*string, error) {
	if err := n.CLI.Reset([]string{outputDir}); err != nil {
		return nil, err
	}

	initDirs := []string{
		fmt.Sprintf("%s/%s", outputDir, ValidatorsDir),
		fmt.Sprintf("%s/%s", outputDir, GentxsDir),
	}

	if err := n.createDirs(initDirs); err != nil {
		return nil, err
	}

	gentxDir := fmt.Sprintf("%s/%s", outputDir, GentxsDir)
	validators := n.initValidators(count, outputDir, seedIPv4Addr)

	for _, validator := range validators {
		appDirs := []string{validator.NodeHomeDir}
		if err := n.createDirs(appDirs); err != nil {
			return nil, err
		}

		if err := n.generateKey(validator); err != nil {
			return nil, err
		}

		if err := n.initChain(validator); err != nil {
			return nil, err
		}

		if err := n.setValidatorAddress(validator); err != nil {
			return nil, err
		}

		if err := n.setValidatorConsensusAddress(validator); err != nil {
			return nil, err
		}

		if err := genesis.ReplaceStakingBondDenom(validator.NodeHomeDir); err != nil {
			return nil, err
		}

		if err := n.setValidatorID(validator); err != nil {
			return nil, err
		}

		if !validator.Seed {
			seedValidator := n.getSeedValidator(validators)

			if err := n.addValidatorKeyToSeed(validator, seedValidator); err != nil {
				return nil, err
			}

			if err := n.addGenesis(validator.Address, seedValidator.NodeHomeDir); err != nil {
				return nil, err
			}

			if err := n.generateTx(validator, seedValidator.NodeHomeDir, gentxDir); err != nil {
				return nil, err
			}
		} else {
			if err := n.addGenesis(validator.Address, validator.NodeHomeDir); err != nil {
				return nil, err
			}

			if err := n.generateTx(validator, validator.NodeHomeDir, gentxDir); err != nil {
				return nil, err
			}
		}
	}

	seedValidator := n.getSeedValidator(validators)
	if err := n.collectGenTxs(gentxDir, seedValidator.NodeHomeDir); err != nil {
		return nil, err
	}

	if err := n.setPeers(validators); err != nil {
		return nil, err
	}

	if err := n.copyGenesis(validators); err != nil {
		return nil, err
	}

	summary := n.summary(validators)
	return &summary, nil
}

func (n *Network) initValidators(count int, outputDir, seedIPv4Addr string) []*Validator {
	var validators []*Validator
	var lastIPv4Addr string

	for i := 0; i < count; i++ {
		seed := false
		if i == 0 {
			seed = true
		}

		if seed {
			lastIPv4Addr = seedIPv4Addr
		}

		validator := NewValidator(outputDir, n.ChainID, seed, lastIPv4Addr)
		validators = append(validators, validator)

		lastIPv4Addr = validator.IPv4Address
	}

	return validators
}

func (n *Network) createDirs(toCreate []string) error {
	for _, dir := range toCreate {
		if err := n.CLI.CreateDir(dir); err != nil {
			return err
		}
	}

	return nil
}

func (n *Network) generateKey(validator *Validator) error {
	output, err := n.CLI.AddKey(validator.Moniker, validator.Mnemonic, validator.Password, fmt.Sprintf("%s/%s", validator.HomeDir, ".sifnoded"))
	if err != nil {
		return err
	}

	yml, err := ioutil.ReadAll(strings.NewReader(*output))
	if err != nil {
		return err
	}

	var keys common.Keys

	err = yaml.Unmarshal(yml, &keys)
	if err != nil {
		return err
	}

	validator.Address = keys[0].Address
	validator.PubKey = keys[0].PubKey

	return nil
}

func (n *Network) addValidatorKeyToSeed(validator, seedValidator *Validator) error {
	_, err := n.CLI.AddKey(validator.Moniker, validator.Mnemonic, seedValidator.Password, fmt.Sprintf("%s/%s", seedValidator.HomeDir, ".sifnoded"))
	if err != nil {
		return err
	}

	return nil
}

func (n *Network) initChain(validator *Validator) error {
	_, err := n.CLI.InitChain(validator.ChainID, validator.Moniker, validator.NodeHomeDir)
	if err != nil {
		return err
	}

	return nil
}

func (n *Network) setValidatorAddress(validator *Validator) error {
	output, err := n.CLI.ValidatorAddress(validator.NodeHomeDir)
	if err != nil {
		return err
	}

	validator.ValidatorAddress = strings.TrimSuffix(*output, "\n")

	return nil
}

func (n *Network) setValidatorConsensusAddress(validator *Validator) error {
	output, err := n.CLI.ValidatorConsensusAddress(validator.NodeHomeDir)
	if err != nil {
		return err
	}

	validator.ValidatorConsensusAddress = strings.TrimSuffix(*output, "\n")

	return nil
}

func (n *Network) setValidatorID(validator *Validator) error {
	output, err := n.CLI.NodeID(validator.NodeHomeDir)
	if err != nil {
		return err
	}

	validator.NodeID = strings.TrimSuffix(*output, "\n")

	return nil
}

func (n *Network) getSeedValidator(validators []*Validator) *Validator {
	for _, validator := range validators {
		if validator.Seed {
			return validator
		}
	}

	return &Validator{}
}

func (n *Network) addGenesis(address, validatorHomeDir string) error {
	_, err := n.CLI.AddGenesisAccount(address, validatorHomeDir, []string{n.MintAmount})
	if err != nil {
		return err
	}

	return nil
}

func (n *Network) generateTx(validator *Validator, validatorDir, outputDir string) error {
	_, err := n.CLI.GenerateGenesisTxn(
		validator.Moniker,
		validator.Password,
		n.BondAmount,
		validatorDir,
		fmt.Sprintf("%s/%s.json", outputDir, validator.Moniker),
		validator.NodeID,
		validator.ValidatorAddress,
		validator.IPv4Address,
		n.ChainID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (n *Network) collectGenTxs(gentxDir, validatorDir string) error {
	_, err := n.CLI.CollectGenesisTxns(gentxDir, validatorDir)
	if err != nil {
		return err
	}

	return nil
}

func (n *Network) generatePeerList(validators []*Validator, idx int) []string {
	var peers []string
	for i, validator := range validators {
		if i != idx {
			peers = append(peers, fmt.Sprintf("%s@%s:26656", validator.NodeID, validator.IPv4Address))
		}
	}

	return peers
}

func (n *Network) setPeers(validators []*Validator) error {
	for i, validator := range validators {
		var config common.NodeConfig

		configFile := fmt.Sprintf("%s/%s/%s", validator.NodeHomeDir, ConfigDir, utils.ConfigFile)

		content, err := ioutil.ReadFile(configFile)
		if err != nil {
			return err
		}

		if _, err := toml.Decode(string(content), &config); err != nil {
			return err
		}

		file, err := os.Create(configFile)
		if err != nil {
			return err
		}

		config.P2P.PersistentPeers = strings.Join(n.generatePeerList(validators, i)[:], ",")
		if err := toml.NewEncoder(file).Encode(config); err != nil {
			return err
		}

		if err := file.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (n *Network) copyGenesis(validators []*Validator) error {
	seedValidator := n.getSeedValidator(validators)
	srcFile := fmt.Sprintf("%s/%s/%s", seedValidator.NodeHomeDir, ConfigDir, utils.GenesisFile)

	for _, validator := range validators {
		if !validator.Seed {
			input, err := ioutil.ReadFile(srcFile)
			if err != nil {
				return err
			}

			err = ioutil.WriteFile(
				fmt.Sprintf("%s/%s/%s", validator.NodeHomeDir, ConfigDir, utils.GenesisFile),
				input,
				0600,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (n *Network) summary(validators []*Validator) string {
	yml, _ := yaml.Marshal(validators)
	return string(yml)
}
