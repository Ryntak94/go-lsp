package document_sync

import "github.com/Ryntak94/go-lsp.git/internal/lsp"

// type (
// 	DocumentUri string
// )

type TextDocumentItem struct {
	Uri        string `json:"uri"`
	LanguageId string `json:"languageId"`
	Version    int32  `json:"version"`
	Text       string `json:"text"`
}

type DidOpenTextDocumentParams struct {
	TextDocumentItem TextDocumentItem `json:"textDocument"`
}

type DidOpenTextMsgJSON struct {
	DidOpenTextDocumentParams DidOpenTextDocumentParams `json:"params"`
	Notification              lsp.Notification
}
