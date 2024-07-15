package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/Ryntak94/go-lsp.git/internal/lsp"
	"github.com/Ryntak94/go-lsp.git/internal/rpc"
)

type ResponseMessage struct {
	Method string
}

func main() {
	logger := getLogger("../../logs/log.txt")
	logger.Println("Started...")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

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
		initialize(logger, method, contents)
	}
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)

	if err != nil {
		panic("hey, you didn't give me a good file")
	}

	return log.New(logfile, "[go-lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}

func initialize(logger *log.Logger, method string, contents []byte) {
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
