package types

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	// ProposalTypeChange defines the type for a ParameterChangeProposal
	ProposalTypeAcceptDistribute = "AcceptDistributionProposal"
)

// Assert ParameterChangeProposal implements govtypes.Content at compile-time
var _ govtypes.Content = AcceptDistributionProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeAcceptDistribute)
	govtypes.RegisterProposalTypeCodec(AcceptDistributionProposal{}, "cosmos-sdk/ParameterChangeProposal")
}

// ParameterChangeProposal defines a proposal which contains multiple parameter
// changes.
type AcceptDistributionProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Changes     string `json:"changes" yaml:"changes"`
}

func (adp AcceptDistributionProposal) String() string {
	return adp.Title
}

func NewParameterChangeProposal(title, description string, changes string) AcceptDistributionProposal {
	return AcceptDistributionProposal{title, description, changes}
}

// GetTitle returns the title of a parameter change proposal.
func (adp AcceptDistributionProposal) GetTitle() string { return adp.Title }

// GetDescription returns the description of a parameter change proposal.
func (adp AcceptDistributionProposal) GetDescription() string { return adp.Description }

// ProposalRoute returns the routing key of a parameter change proposal.
func (adp AcceptDistributionProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a parameter change proposal.
func (adp AcceptDistributionProposal) ProposalType() string { return ProposalTypeAcceptDistribute }

// ValidateBasic validates the parameter change proposal
func (adp AcceptDistributionProposal) ValidateBasic() error {
	return nil
}
