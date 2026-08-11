package lsp

type TypeDefinitionRequest struct {
	Request
	Params TypeDefinitionParams `json:"params"`
}

type TypeDefinitionParams struct {
	TextDocumentPositionParams
}
