package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Exported code type numbers
var (
	ErrProphecyNotFound              = sdkerrors.Register(ModuleName, 2, "prophecy with given id not found")
	ErrMinimumConsensusNeededInvalid = sdkerrors.Register(ModuleName, 3,
		"minimum consensus proportion of validator staking power must be > 0 and <= 1")
	ErrNoClaims          = sdkerrors.Register(ModuleName, 4, "cannot create prophecy without initial claim")
	ErrInvalidIdentifier = sdkerrors.Register(ModuleName, 5,
		"invalid identifier provided, must be a nonempty string")
	ErrProphecyFinalized = sdkerrors.Register(ModuleName, 6, "prophecy already finalized")
	ErrDuplicateMessage  = sdkerrors.Register(ModuleName, 7,
		"already processed message from validator for this id")
	ErrInvalidClaim            = sdkerrors.Register(ModuleName, 8, "claim cannot be empty string")
	ErrInvalidValidator        = sdkerrors.Register(ModuleName, 9, "claim must be made by actively bonded validator")
	ErrInternalDB              = sdkerrors.Register(ModuleName, 10, " failed prophecy serialization/deserialization")
	ErrValidatorNotInWhiteList = sdkerrors.Register(ModuleName, 11, "validator must be in whitelist")
	ErrNotAdminAccount         = sdkerrors.Register(ModuleName, 12, "Not an admin account")
	ErrInvalidOperationType    = sdkerrors.Register(ModuleName, 13, "invalid operation type for validator whitelist")
	ErrInvalidProphecyStatus   = sdkerrors.Register(ModuleName, 14, "invalid prophecy status")
	ErrValidatorPowerOverflow  = sdkerrors.Register(ModuleName, 15, "validator power setting is overflow")
)
