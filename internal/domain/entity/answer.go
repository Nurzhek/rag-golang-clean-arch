package entity

// Answer is the result of a retrieval-augmented generation query: the generated
// text together with the source chunks that grounded it.
type Answer struct {
	// Text is the natural-language answer produced by the LLM.
	Text string
	// Sources are the retrieved chunks used as context, in relevance order.
	Sources []ScoredDocument
}
