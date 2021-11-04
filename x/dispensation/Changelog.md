# Change Log for Dispensation module

---
### 10/19/2021
- Removed validation check which limits dispensation to occur only in rowan .
- Ex to distribute 100 Rowan and 100 Ceth to address `sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd` ,the output.json should look like 
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
```json
{
  "Output": [
    {
      "address": "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
      "coins": [
        {
          "denom": "rowan",
          "amount": "100"
        },
        {
          "denom": "ceth",
          "amount": "100"
        }
      ]
    }
  ]
}
```


```
- The Distribution creator should have enough balances in all tokens 
- Internally only one record is created as the coins are aggregated per address , the records created wouldbe the same in either scenario.
- NOTE : Number of records created will be equivalent to the number of unique address in the output.json .
----