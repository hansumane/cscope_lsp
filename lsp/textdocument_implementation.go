package lsp

type ImplementationRequest struct {
	Request
	Params ImplementationParams `json:"params"`
}

type ImplementationParams struct {
	TextDocumentPositionParams
}
