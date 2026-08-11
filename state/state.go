package state

import (
	"log"
	"strings"

	"github.com/dhananjaylatkar/cscope_lsp/cscope_if"
	"github.com/dhananjaylatkar/cscope_lsp/lsp"
)

type State struct {
	// Map of file names to contents
	Documents map[string]string

	// Whether to replace:
	// - textDocument/definition with textDocument/typeDefinition
	// - textDocument/references with textDocument/implementation
	ReplaceMethods bool
}

func New() State {
	return State{Documents: map[string]string{}}
}

func (s *State) Update(uri, text string) {
	s.Documents[uri] = text
}

func isWordChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}

func extractWord(line string, pos int) string {
	start := pos
	end := pos

	for start > 0 && isWordChar(line[start]) {
		start--
	}

	if !isWordChar(line[start]) {
		start++
	}

	for end < len(line) && isWordChar(line[end]) {
		end++
	}

	if start >= end {
		return ""
	}

	return line[start:end]
}

func (s *State) EmptyLocation(id int) lsp.LocationResponse {
	return lsp.LocationResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: []lsp.Location{},
	}
}

func (s *State) Definition(id int, uri string, logger *log.Logger, position lsp.Position) lsp.LocationResponse {
	logger.Printf("uri: %s", uri)
	logger.Printf("position.Line: %d", position.Line)
	logger.Printf("position.Char: %d", position.Character)
	line := strings.Split(s.Documents[uri], "\n")[position.Line]
	word := extractWord(line, position.Character)

	logger.Printf("line: %s", line)
	logger.Printf("word: %s", word)

	defs := cscope_if.GetDefinition(logger, uri, word)

	return lsp.LocationResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: defs,
	}
}

func (s *State) TypeDefinition(id int, uri string, logger *log.Logger, position lsp.Position) lsp.LocationResponse {
	return s.Definition(id, uri, logger, position)
}

func (s *State) References(id int, uri string, logger *log.Logger, position lsp.Position) lsp.LocationResponse {
	logger.Printf("uri: %s", uri)
	logger.Printf("position.Line: %d", position.Line)
	logger.Printf("position.Char: %d", position.Character)

	line := strings.Split(s.Documents[uri], "\n")[position.Line]
	word := extractWord(line, position.Character)
	logger.Printf("line: %s", line)
	logger.Printf("word: %s", word)

	defs := cscope_if.GetReferences(logger, uri, word)

	return lsp.LocationResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: defs,
	}
}

func (s *State) Implementation(id int, uri string, logger *log.Logger, position lsp.Position) lsp.LocationResponse {
	return s.References(id, uri, logger, position)
}
