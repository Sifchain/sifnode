package reflect

type Swap struct {
	Amount uint32 `json:"amount,omitempty"`
}

type ReflectCustomMsg struct {
	Swap *Swap `json:"swap,omitempty"`
}
