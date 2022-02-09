package reflect

type Swap struct {
	Amount uint32 `json:"amount,omitempty"`
}

type SifchainMsg struct {
	Swap *Swap `json:"swap,omitempty"`
}
