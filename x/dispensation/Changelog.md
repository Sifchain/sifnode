# Change Log for Dispensation module

---
### 10/19/2021
- Removed validation check which limits dispensation to occur only in rowan . 
- Every dispensation record still has a limit of 1 token .
- Ex to send 100 Rowan and 100 Ceth to address `sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd` ,the output.json should look like 
```json
{
 "Output": [
  {
   "address": "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
   "coins": [
    {
     "denom": "rowan",
     "amount": "100"
    }
   ]
  },
  {
   "address": "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
   "coins": [
    {
     "denom": "ceth",
     "amount": "100"
    }
   ]
  }
 ]
}
```
- And the Distribution creator should have enough balances in all tokens 
- Internally only one record is created as the coins are aggregated per address.
- NOTE : Number of records created will be equivalent to the the number of unique address in the output.json and not the number of entries .
----