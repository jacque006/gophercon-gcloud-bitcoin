package messages

// Represents a transaction we will write to the graph database.
type TransactionMessage struct {
	Hash string `json:"hash"`
	Time int32  `json:"time"`
	// key: address, value: value
	In  []TransactionMessagePart `json:"in"`
	Out []TransactionMessagePart `json:"out"`
}

type TransactionMessagePart struct {
	Address string  `json:"address"`
	Amount  float64 `json:"amount"`
}
