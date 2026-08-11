package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/dhananjaylatkar/cscope_lsp/lsp"
	"github.com/dhananjaylatkar/cscope_lsp/rpc"
	"github.com/dhananjaylatkar/cscope_lsp/state"
)

const (
	initBufSize = 64 * 1024
	maxBufSize  = 10 * 1024 * 1024
)

func main() {
	logger := getLogger("/tmp/cscope_lsp.log")
	logger.Println("Started!")

	scanner := bufio.NewScanner(os.Stdin)
	buf := make([]byte, 0, initBufSize)
	scanner.Buffer(buf, maxBufSize)
	scanner.Split(rpc.Split)

	writer := os.Stdout
	state := new(state.New())

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Got an error: %s", err)
			continue
		}

		handleMessage(logger, writer, state, method, contents)
	}
}

func handleMessage(logger *log.Logger, writer io.Writer, state *state.State, method string, contents []byte) {
	logger.Printf("Received msg with method: %s", method)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Hey, we couldn't parse this: %s", err)
		}

		logger.Printf("Connected to: %s %s",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version)

		if request.Params.Options != nil {
			state.ReplaceMethods = request.Params.Options.ReplaceMethods
			logger.Printf("Client Options: %v, replaceMethods: %v", request.Params.Options, state.ReplaceMethods)
		}

		// sent initialize response
		msg := lsp.NewInitializeResponse(request.ID)
		writeResponse(writer, msg)

		logger.Print("Sent the reply")

	case "shutdown":
		var request lsp.ShutdownRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("shutdown: %s", err)
		}

		// sent shutdown response
		msg := lsp.NewShutdownResponse(request.ID)
		writeResponse(writer, msg)

		logger.Print("Shutdown")

	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/didOpen: %s", err)
			return
		}

		logger.Printf("Opened: %s", request.Params.TextDocument.URI)
		state.Update(request.Params.TextDocument.URI, request.Params.TextDocument.Text)

	case "textDocument/didChange":
		var request lsp.TextDocumentDidChangeNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/didChange: %s", err)
			return
		}

		logger.Printf("Changed: %s", request.Params.TextDocument.URI)
		for _, change := range request.Params.ContentChanges {
			state.Update(request.Params.TextDocument.URI, change.Text)
		}

	case "textDocument/definition":
		var request lsp.DefinitionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/definition: %s", err)
			return
		}

		// Create a response
		if state.ReplaceMethods {
			response := state.EmptyLocation(request.ID)
			writeResponse(writer, response)
		} else {
			response := state.Definition(request.ID, request.Params.TextDocument.URI, logger, request.Params.Position)
			writeResponse(writer, response)
		}

	case "textDocument/typeDefinition":
		var request lsp.TypeDefinitionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/typeDefinition: %s", err)
			return
		}

		// Create a response
		if state.ReplaceMethods {
			response := state.TypeDefinition(request.ID, request.Params.TextDocument.URI, logger, request.Params.Position)
			writeResponse(writer, response)
		} else {
			response := state.EmptyLocation(request.ID)
			writeResponse(writer, response)
		}

	case "textDocument/references":
		var request lsp.ReferencesRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/references: %s", err)
			return
		}

		// Create a response
		if state.ReplaceMethods {
			response := state.EmptyLocation(request.ID)
			writeResponse(writer, response)
		} else {
			response := state.References(request.ID, request.Params.TextDocument.URI, logger, request.Params.Position)
			writeResponse(writer, response)
		}

	case "textDocument/implementation":
		var request lsp.ImplementationRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/implementation: %s", err)
			return
		}

		// Create a response
		if state.ReplaceMethods {
			response := state.Implementation(request.ID, request.Params.TextDocument.URI, logger, request.Params.Position)
			writeResponse(writer, response)
		} else {
			response := state.EmptyLocation(request.ID)
			writeResponse(writer, response)
		}
	}
}

func writeResponse(writer io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)
	writer.Write([]byte(reply))
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("hey, you didnt give me a good file")
	}

	return log.New(logfile, "[cscope_lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
