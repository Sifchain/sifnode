
package types

import (
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	// ProposalTypeUpdatePmtpParams defines the type for a UpdatePmtpParamsProposal
	ProposalTypeUpdatePmtpParams = "UpdatePmtpParams"
)

// Assert UpdatePmtpParamsProposal implements govtypes.Content at compile time
var _ govtypes.Content = &UpdatePmtpParamsProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeUpdatePmtpParams)
	govtypes.RegisterProposalTypeCodec(&UpdatePmtpParamsProposal{}, "sifnode/UpdatePmtpParamsProposal")
}

// UpdatePmtpParamsProposal defines a proposal to update pmtp parameters.
type UpdatePmtpParamsProposal struct {
	Title       string     `json:"title" yaml:"title"`
	Description string     `json:"description" yaml:"description"`
	RunAs       string     `json:"run_as" yaml:"run_as"`
	PmtpParams  *PmtpParams `json:"pmtp_params" yaml:"pmtp_params"`
}

// NewUpdatePmtpParamsProposal creates a new UpdatePmtpParamsProposal.
func NewUpdatePmtpParamsProposal(title, description, runAs string, pmtpParams *PmtpParams) *UpdatePmtpParamsProposal {
	return &UpdatePmtpParamsProposal{
		Title:       title,
		Description: description,
		RunAs:       runAs,
		PmtpParams:  pmtpParams,
	}
}

// GetTitle returns the title of a community pool spend proposal.
func (p *UpdatePmtpParamsProposal) GetTitle() string { return p.Title }

// GetDescription returns the description of a community pool spend proposal.
func (p *UpdatePmtpParamsProposal) GetDescription() string { return p.Description }

// ProposalRoute returns the routing key of a community pool spend proposal.
func (p *UpdatePmtpParamsProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (p *UpdatePmtpParamsProposal) ProposalType() string { return ProposalTypeUpdatePmtpParams }

// ValidateBasic runs basic governance proposal checks and checks that the proposal details are valid.
func (p *UpdatePmtpParamsProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}
	if p.PmtpParams == nil {
		return ErrInvalidPmtpParams
	}
	return p.PmtpParams.ValidateBasic()
}

// String implements the Stringer interface.
func (p *UpdatePmtpParamsProposal) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
