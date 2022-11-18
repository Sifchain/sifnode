# Asymmetric Liquidity Adds

Sifnoded does not currently support asymmetric liquidity adds. This document proposes a procedure
which would allow asymmetric adds.

## Symmetric Adds

When adding symmetrically to a pool the fraction of total pool units owned by the Liquidity Provider (LP)
after the add must equal the amount of native token added to the pool as a fraction of total native asset token in the
pool (after the add):

```
l / (P + l) = r / (r + R)
```

Where:
```
l - LP units
P - total pool units (before)
r - amount of native token added
R - native asset pool depth (before)
```
Rearranging gives:

```
(1) l = r * P / R
```

## Asymmetric adds

In the asymmetric case, by definition:

```
R/A =/= r/a
```

(this includes the case where the division is not defined i.e. when a=0 the division is not defined
in which case the add is considered asymmetric)

Where:
```
R - native asset pool depth (before adding liquidity)
A - external asset pool depth (before adding liquidity)
r - amount of native token added
a - amount of external token added
```
Currently sifnoded blocks asymmetric adds. The following procedure is proposed to enable
asymmetric adds.

### Proposed method

In the following formulas:

```
p - current ratio shifting running rate
f - swap fee rate
```

If the pool is not in the same ratio as the add then either:

1. Some r must be swapped for a, such that after the swap the add is symmetric
2. Some a must be swapped for r, such that after the swap the add is symmetric

#### Swap native token for external token

Swap an amount, s, of native token such that:

```
(R + s) / (A - g.s) = (R + r) / (A + a) = (r − s) / (a + g.s)
```

where g is the swap formula:

```
g.x = (1 - f) * (1 + r) * x * Y / (x + X)
```

Solving for s (using mathematica!) gives:

```
s = abs((sqrt(pow((-1*f*p*A*r-f*p*A*R-f*A*r-f*A*R+p*A*r+p*A*R+2*a*R+2*A*R), 2)-4*(a+A)*(a*R*R-A*r*R)) + f*p*A*r + f*p*A*R + f*A*r + f*A*R - p*A*r - p*A*R - 2*a*R - 2*A*R) / (2 * (a + A))).
```

The number of pool units is then given by the symmetric formula (1):

```
l = (r - s) * P / (R + s)
```

#### Swap external token for native token

Swap an amount, s, of native token such that:

```
(R - s) / (A + g'.s) = (R + r) / (A + a) = (r + g'.s) / (a - s)
```

Where g' is the swap formula:

```
g' = (1 - f) * x * Y / ((x + X) * (1 + r))
```

Solving for s (using mathematica!) gives:

```
s = abs((sqrt(R*(-1*(a+A))*(-1*f*f*a*R-f*f*A*R-2*f*p*a*R+4*f*p*A*r+2*f*p*A*R+4*f*A*r+4*f*A*R-p*p*a*R-p*p*A*R-4*p*A*r-4*p*A*R-4*A*r-4*A*R)) + f*a*R + f*A*R + p*a*R - 2*p*A*r - p*A*R - 2*A*r - 2*A*R) / (2 * (p + 1) * (r + R)))
```

The number of pool units is then given by the symmetric formula (1):

```
l = (a - s) * P / (A + s)
```

### Equivalence with swapping

Any procedure which assigns LP units should guarantee that if an LP adds (x,y) then removes all their
liquidity from the pool, receiving (x',y') then it is never the case that x' > x and y' > y (all else being equal i.e.
no LPD, rewards etc.). Furthermore
assuming (without loss of generality) that x' =< x, if instead of adding to the pool then removing all liquidity
the LP had swapped (x - x'), giving them a total of y'' (i.e. y'' = y + g.(x - x')), then y'' should equal y'. (Certainly y' cannot be greater than y'' otherwise 
the LP has achieved a cheap swap.)

In the case of the proposed add liquidity procedure the amount the LP would receive by adding then removing would equal the amounts
of each token after the internal swap (at this stage the add is symmetric and with symmetric adds x' = x and y' = y), that is:
```
(2) x' = x - s
(3) y' = y + g.s
```
Plugging these into the equation for y'', y'' = y + g.(x - x'):
```
y'' = y + g.(x - x')
    = y + g.s by rearranging (2) and substituting 
    = y' by substituting (3)
```
### Liquidity Protection

Since the add liquidity process involves swapping then the Liquidity protection procedure must be applied.

## CLI

### Get Estimation of pool share

```bash
sifnoded query clp estimate-pool-share \
  --externalAmount=0 \
  --nativeAmount=1000000000000000000 \
  --symbol ceth \
  --output json
```

```json
{
	"percentage": "0.183227703974514619",
	"native_asset_amount": "549683111923543857",
	"external_asset_amount": "366455407949029238",
	"swap_info": {
		"status": "SELL_NATIVE",
		"fee": "1102674246586848",
		"fee_rate": "0.003000000000000000",
		"amount": "450316888076456143",
		"result": "366455407949029237"
	}
}
```

Where `swap_info` `swap_status` is one of `NO_SWAP`, `SELL_NATIVe`, `BUY_NATIVE`. 

## References

Detailed derivation of formulas https://hackmd.io/NjvaZY1qQiS17s_uEgZmTw?both
