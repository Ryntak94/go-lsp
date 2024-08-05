package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/Ryntak94/go-lsp.git/internal/keywords"
	"github.com/Ryntak94/go-lsp.git/internal/lsp"
	"github.com/Ryntak94/go-lsp.git/internal/rpc"
)

type ResponseMessage struct {
	Method string
}

func main() {
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	pathToLogger := path[:strings.Index(path, "/go-lsp")+7] + "/logs/log.txt"

	logger := getLogger(pathToLogger)
	logger.Println("Started...")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	keywordTree := keywords.GenerateKeywords(logger)
	for _, keyword := range keywordTree.FindWords() {
		_ = keyword
	}
	for scanner.Scan() {
		msg := scanner.Bytes()

		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Got an error: %s", err)
			continue
		}

		logger.Printf("Received message with method: %s", method)
		handleMessage(logger, method, contents)
	}
}

func handleMessage(logger *log.Logger, method string, contents []byte) {
	switch method {
	case "initialize":
		initialize(logger, contents)
	// adding stubbing for other lifecycle methods
	case "initialized":
		logger.Printf("initialized message received")
	case "client/registerCapability":
		logger.Printf("client/registerCapability message received")
	case "client/unregisterCapability":
		logger.Printf("client/unregisterCapability messsage received")
	case "$/setTrace":
		logger.Printf("$/setTrace message received")
	case "$/logTrace":
		logger.Printf("$/logTrace message received")
	case "shutdown":
		logger.Printf("shutdown message received")
	case "exit":
		logger.Printf("exit message received")
	default:
		logger.Printf("No handling for message with method: %s", method)
	}

}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("hey, you didn't give me a good file")
	}

	return log.New(logfile, "[go-lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}

func initialize(logger *log.Logger, contents []byte) {
	var request lsp.InitializeRequest
	if err := json.Unmarshal(contents, &request); err != nil {
		logger.Printf("Could not parse message content: %s", err)
	}
	logger.Printf("Connected to: %s %s",
		request.Params.ClientInfo.Name,
		*request.Params.ClientInfo.Version)

	msg := lsp.NewInitializeResponse(request.ID)
	response := rpc.EncodeMessage(msg)

	writer := os.Stdout
	writer.Write([]byte(response))

	logger.Printf("Sent the reply...")
	logger.Printf(response)
}
