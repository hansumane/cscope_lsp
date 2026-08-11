package lsp

type InitializeRequest struct {
	Request
	Params InitializeRequestParams `json:"params"`
}

type InitializeRequestParams struct {
	ClientInfo *ClientInfo        `json:"clientInfo"`
	Options    *InitializeOptions `json:"initializationOptions"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeOptions struct {
	ReplaceMethods bool `json:"replaceMethods"`
}

type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

type ServerCapabilities struct {
	TextDocumentSync       int  `json:"textDocumentSync"`
	DefinitionProvider     bool `json:"definitionProvider"`
	TypeDefinitionProvider bool `json:"typeDefinitionProvider"`
	ReferencesProvider     bool `json:"referencesProvider"`
	ImplementationProvider bool `json:"implementationProvider"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewInitializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync:       1,
				DefinitionProvider:     true,
				TypeDefinitionProvider: true,
				ReferencesProvider:     true,
				ImplementationProvider: true,
			},
			ServerInfo: ServerInfo{
				Name:    "cscope_lsp",
				Version: "0.1",
			},
		},
	}
}
