# Symmetric Adds

When adding symmetrically to a pool the fraction of total pool units owned by the Liquidity Provider (LP)
equals the amount of native token added to the pool as a fraction of total native asset token in the
pool:

```
l / (P + l) = r / (r + R)
```

where:

l - LP units
P - total pool units (before)
r - amount of native token added
R - native asset pool depth (before)

Rearranging gives:

```
(1) l = r * P / R
```

# Asymmetric adds

In the asymmetric case, by definition:

```
R/A =/= r/a
```

(this includes the case where the division is not defined i.e. when a=0 the division is not defined
in which case the add is considered asymmetric)

Where:

R - native asset pool depth (before adding liquidity)
A - external asset pool depth (before adding liquidity)
r - amount of native token added
a - amount of external token added

Currently sifnoded blocks asymmetric adds. The following procedure is proposed to enable
asymmetric adds.

## Proposed method

If the pool is not in the same ratio as the add then either:

i. Some r must be swapped for a, such that after the swap the add is symmetric
ii. Some a must be swapped for r, such that after the swap the add is symmetric

### Swap native token for external token

Swap an amount, s, of native token such that:

```
(R + s) / (A - g.s) = (R + r) / (A + a) = (r − s) / (a + g.s)
```

where g is the swap formula

```
g.x = s * Y / (x + X)
```

Solving for s (using mathematica!) gives:

```
s =  abs((sqrt(pow((-1*f*r*X*y-f*r*X*Y-f*X*y-f*X*Y+r*X*y+r*X*Y+2*x*Y+2*X*Y), 2)-4*(x+X)*(x*Y*Y-X*y*Y)) + f*r*X*y + f*r*X*Y + f*X*y + f*X*Y - r*X*y - r*X*Y - 2*x*Y - 2*X*Y) / (2 * (x + X)))
```

The number of pool units is then given by the symmetric formula (1):

```
l = (r - s) * P / (R + s)
```

### Swap external token for native token

Swap an amount, s, of native token such that:

```
(R - s) / (A + g.s) = (R + r) / (A + a) = (r + g.s) / (a - s)
```

Solving for s (using mathematica!) gives:

```
s = abs((sqrt(Y*(-1*(x+X))*(-1*f*f*x*Y-f*f*X*Y-2*f*r*x*Y+4*f*r*X*y+2*f*r*X*Y+4*f*X*y+4*f*X*Y-r*r*x*Y-r*r*X*Y-4*r*X*y-4*r*X*Y-4*X*y-4*X*Y)) + f*x*Y + f*X*Y + r*x*Y - 2*r*X*y - r*X*Y - 2*X*y - 2*X*Y) / (2 * (r + 1) * (y + Y)))
```

The number of pool units is then given by the symmetric formula (1):

```
l = (a - s) * P / (A + s)
```

## Equivalence with swapping

Any procedure which assigns LP units should guarantee that if an LP adds (x,y) then removes all their
liquidity from the pool, receiving (x',y') then it is never the case that x' > x and y' > y (all else being equal i.e.
no LPD, rewards etc.). Furthermore
assuming (without loss of generality) that x' =< x, if instead of adding to the pool then removing all liquidity
the LP had swapped (x - x'), giving them a total of y'' (i.e. y'' = y + g.(x - x')), then y'' should equal y'. (Certainly y' cannot be greater than y'' otherwise 
the LP has achieved a cheap swap.)

In the case of the proposed add liquidity procedure the amount the LP would receive by adding then removing would equal the amounts
of each token after the internal swap (at this stage the add is symmetric and with symmetric adds x' = x and y' = y), that is:

(2) x' = x - s
(3) y' = y + g.s

Plugging these into the equation for y'', y'' = y + g.(x - x')):

y'' = y + g.(x - x')
    = y + g.s by rearranging (2) and substituting 
    = y' by substituting (3)

## Liquidity Protection

Since the add liquidity process involves swapping then the Liquidity protection procedure must be applied.

