package lsp

type ReferencesRequest struct {
	Request
	Params ReferencesParams `json:"params"`
}

type ReferencesParams struct {
	TextDocumentPositionParams
}
