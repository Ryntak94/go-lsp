package document_sync

type (
	DocumentUri string
)

type TextDocumentItem struct {
	uri        DocumentUri
	languageId string
	version    int32
	text       string
}

type DidOpenTextDocumentParams struct {
	textDocument TextDocumentItem
}
