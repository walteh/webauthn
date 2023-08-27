package lambda

// BatchItemFailure is the individual record which failed processing.
type BatchItemFailure struct {
	ItemIdentifier string `json:"itemIdentifier"`
}
