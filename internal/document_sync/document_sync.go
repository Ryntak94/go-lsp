package document_sync

import "log"

func DidOpenHandler(logger *log.Logger, contents []byte) {
	logger.Printf("Hello from didOpenHandler")
}
