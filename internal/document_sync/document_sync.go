package document_sync

import (
	"encoding/json"
	"log"
)

func DidOpenHandler(logger *log.Logger, contents []byte, stateMap map[string]int) {
	var didOpenTextMsgJSON DidOpenTextMsgJSON
	if err := json.Unmarshal(contents, &didOpenTextMsgJSON); err != nil {
		logger.Printf("Could not parse message content: %s", err)
	}

	languageId := didOpenTextMsgJSON.DidOpenTextDocumentParams.TextDocumentItem.LanguageId
	uri := string(didOpenTextMsgJSON.DidOpenTextDocumentParams.TextDocumentItem.DocumentUri)

	// Throw Error for non go file
	if languageId != "go" {
		panic("INVALID LANGUAGE ID: " + languageId)
	} else {
		logger.Printf("go file detected")
		stateMap[uri] = 1 // declare file as open
		logger.Printf("after statemap")
		ParseText(logger, didOpenTextMsgJSON.DidOpenTextDocumentParams.TextDocumentItem.Text)
	}
}

func ParseText(logger *log.Logger, text string) {
	logger.Printf("text: %s", text)
}
