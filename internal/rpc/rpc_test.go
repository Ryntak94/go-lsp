package rpc_test

import (
	"testing"

	"github.com/Ryntak94/go-lsp.git/internal/rpc"
)

type EncodingExample struct {
	Testing bool
}

func TestEncode(t *testing.T) {
	expected := "Content-Length: 16\r\n\r\n{\"Testing\":true}"

	actual := rpc.EncodeMessage(EncodingExample{Testing: true})

	if expected != actual {
		t.Fatalf("Expected: %s, Actual: %s", expected, actual)
	}
}

func TestDecode(t *testing.T) {
	incomingMessage := "Content-Length: 17\r\n\r\n{\"Method\":\"test\"}"

	method, content, err := rpc.DecodeMessage([]byte(incomingMessage))
	contentLength := len(content)

	if err != nil {
		t.Fatal(err)
	}

	if contentLength != 17 {
		t.Fatalf("Expected contentLength: 17, Actual: %d", contentLength)
	}

	if method != "test" {
		t.Fatalf("Expected method: test, Actual %s", method)
	}
}
