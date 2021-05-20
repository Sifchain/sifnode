## Overview
- The module allows a user to create a Claim .
- A claim is record , stating that the user want to collect their rewards at the end of the week.
- A claim is deleted once the user is paid out.


## General use case

Any day of the week
1) Create claims through this api (on chain)
     Users send a create claim request from their address and provide the claim-type (LM/VS) as a parameter. The claims are valid until they are paid out.


On Friday ( No limitation on the day being friday )


2) Run a query to get all claims till that time , and create a list 
     - Step(1) also emits a claim_created event , which can read to get the claims and create the list
    
3) This  list is an input for a function (This function is off-chain and not part of this module)  which iterates over the list and calculates the rewards earned by that address. This data is then used to create a distribution list . 
     
4) We run a distribution using the list from step (3). 
5) The transfers happen over the next few blocks (10 per block) . 
6) An event is emitted when a transfer to a user happens.The event is an indicator that the recipient has received the funds .This will be used by an external function to reset multipliers for the recievers. 

### Out of Scope
- No mechanism present to delete a user-claim once created .The claims are automatically deleted by the module once processed.
- If a user who does not have any pending rewards creates a claim , their claims would not be processed by the module. The module would stop them from submitting requests in the future.These users would stay blocked until they do earn rewards at some point in time ,and then the rewards would be distributed to them , and their addresses unblocked.
The above is true for claim_type , if a user gets blocked due to creating a false claim for Validator Subsidy ,they would still be able to create new claims for Liquidity Mining.
- There will be a short time gap between the time the distribution list is created, and the user receives the rewards (between steps 3 and 5). The rewards are calculated at step 3, and the multipliers are reset at step 6 . If the user accumulates any reward in this duration , there is a chance the user might lose out these.
Do give an idea of the time frame , we do 10 distributions per block , and have a block time of 5s . ( The number 10 is configurable and can be increased )