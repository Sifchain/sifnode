package dispensation

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	// ProposalTypeAcceptDistribute defines the type for a AcceptDistributionProposal
	ProposalTypeAcceptDistribute = "AcceptDistributionProposal"
)

var _ govtypes.Content = AcceptDistributionProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeAcceptDistribute)
	govtypes.RegisterProposalTypeCodec(AcceptDistributionProposal{}, "cosmos-sdk/AcceptDistributionProposal")
}

// AcceptDistributionProposal defines a proposal ,the resultof which would be distribution of mining rewards

type AcceptDistributionProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Changes     string `json:"changes" yaml:"changes"`
}

func (adp AcceptDistributionProposal) String() string {
	return adp.Title
}

func NewAcceptDistributionProposal(title, description string, changes string) AcceptDistributionProposal {
	return AcceptDistributionProposal{title, description, changes}
}

// GetTitle returns the title of aAcceptDistributionProposal
func (adp AcceptDistributionProposal) GetTitle() string { return adp.Title }

// GetDescription returns the description of a AcceptDistributionProposal
func (adp AcceptDistributionProposal) GetDescription() string { return adp.Description }

// ProposalRoute returns the routing key of a AcceptDistributionProposal
func (adp AcceptDistributionProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a AcceptDistributionProposal
func (adp AcceptDistributionProposal) ProposalType() string { return ProposalTypeAcceptDistribute }

// ValidateBasic validates the AcceptDistributionProposal
func (adp AcceptDistributionProposal) ValidateBasic() error {
	return nil
}
