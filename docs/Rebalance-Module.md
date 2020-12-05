# **Sifchain Minting Module**

## Changelog
-First Draft: Austin Haines

SET PARAMS IN REBALANCER

## Context
The module outlined below is necessary for Sifchain's Dynamic Reward Rebalancing Policy.  



## Types

**rebalancing.go** 

```golang
package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Minter represents the minting state.
type Rebalancing struct {
	gamma_l     sdk.Dec `json:"gamma_l" yaml:"gamma_l"`     // Target in liquidity subsystem, updated in governance module
    gamma_v     sdk.Dec `json:"gamma_l" yaml:"gamma_l"`     // Target in validator subsystem, updated in governance module
    rho_l       sdk.Dec `json:"rho_l" yaml:"rho_l"`         // Ratio of tokens in liquidity subsystem, computed within module
    rho_v       sdk.Dec `json:"rho_v" yaml:"rho_v"`         // Ratio of tokens in validator subsystem, computed within module
    error_l     sdk.Dec `json:"error_l" yaml:"error_l"`     // Error in liquidity subsystem, computed in module
    error_v     sdk.Dec `json:"error_v" yaml:"error_v"`     // Error in validator subsystem, computed in module
    S           uint64  `json:"S" yaml:"S"`                 // Total Rowan in existence
    S_c         uint64  `json:"S_c" yaml:"S_c"`             // Circulating supply, difference from all other token states
    S_v         uint64  `json:"S_v" yaml:"S_v"`             // Tokens locked in validator stake
    S_l         uint64  `json:"S_l" yaml:"S_l"`             // Rowan tokens in liquidity pools
    lambda_l    sdk.Dec `json:"lambda_l" yaml:"lambda_l"`   // 	Rebalancing coefficent for liquidity subsystem, computed in module
    lambda_v    sdk.Dec `json:"lambda_v" yaml:"lambda_v"`   // 	Rebalancing coefficent for validator subsystem, computed in module

}
```

**governance.go** 

```golang
package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Minter represents the minting state.
type Governance struct {
	Governance      address `json:"Governance" yaml:"Governance"`       // Module address
    PubKey          pubkey  `json:"PubKey" yaml:"PubKey"`               // Module Public Key
    gamma_l         sdk.Dec `json:"gamma_l" yaml:"gamma_l"`             // Target in liquidity subsystem, updated in governance module
    gamma_v         sdk.Dec `json:"gamma_v" yaml:"gamma_v"`             // Target in validator subsystem, updated in governance module
}
```

**validator.go** 

```golang
package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Minter represents the minting state.
type Validator struct {
	Validator       address `json:"Validator" yaml:"Validator"`     // Module address
    PubKey          pubkey  `json:"PubKey" yaml:"PubKey"`             // Module Public Key
    S_v             sdk.Dec `json:"S_v" yaml:"S_v"`                   // Tokens locked in validator stake
    S               sdk.Dec `json:"S" yaml:"S"`                       // Total Rowan in existence Updated from validator due to minting and burning functions
}
```

**liquidity.go** 

```golang
package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Minter represents the minting state.
type Liquidity struct {
	Liquidity       address `json:"Liquidity" yaml:"Liquidity"`       // Module address
    PubKey          pubkey  `json:"PubKey" yaml:"PubKey"`             // Module Public Key
    S_l             sdk.Dec `json:"S_l" yaml:"S_l"`                   // Rowan tokens in liquidity pools
}
```
