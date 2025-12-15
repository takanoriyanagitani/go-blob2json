package blob2json_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/takanoriyanagitani/go-blob2json"
)

func TestBlobBuilder_PlainText(t *testing.T) {
	input := "Hello, Go!"
	expectedBody := base64.StdEncoding.EncodeToString([]byte(input))

	builder := blob2json.BlobBuilder{
		ContentType: "text/plain",
		MaxBytes:    1024,
	}

	blob, err := builder.NewBlobFromReader(strings.NewReader(input), "greeting.txt")
	if err != nil {
		t.Fatalf("Failed to create blob: %v", err)
	}

	if blob.Name != "greeting.txt" {
		t.Errorf("Expected name 'greeting.txt', got '%s'", blob.Name)
	}
	if blob.ContentType != "text/plain" {
		t.Errorf("Expected content type 'text/plain', got '%s'", blob.ContentType)
	}
	if blob.Body != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, blob.Body)
	}
	if blob.ContentTransferEncoding != "base64" {
		t.Errorf("Expected content transfer encoding 'base64', got '%s'", blob.ContentTransferEncoding)
	}
	if blob.Metadata != nil {
		t.Errorf("Expected no metadata, got %s", blob.Metadata)
	}
	if *blob.ContentLength != int64(len(input)) {
		t.Errorf("Expected content length %d, got %d", len(input), *blob.ContentLength)
	}
}

func TestBlobBuilder_WithMetadata(t *testing.T) {
	input := "data"
	expectedBody := base64.StdEncoding.EncodeToString([]byte(input))
	metadata := map[string]string{
		"source": "test",
		"user":   "gemini",
	}

	builder := blob2json.BlobBuilder{
		ContentType: "application/octet-stream",
		MaxBytes:    1024,
		Metadata:    metadata,
	}

	blob, err := builder.NewBlobFromReader(strings.NewReader(input), "file.bin")
	if err != nil {
		t.Fatalf("Failed to create blob: %v", err)
	}

	if blob.Name != "file.bin" {
		t.Errorf("Expected name 'file.bin', got '%s'", blob.Name)
	}
	if blob.Body != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, blob.Body)
	}

	var parsedMeta map[string]string
	if err := json.Unmarshal(blob.Metadata, &parsedMeta); err != nil {
		t.Fatalf("Failed to unmarshal metadata: %v", err)
	}

	if val, ok := parsedMeta["source"]; !ok || val != "test" {
		t.Errorf("Expected metadata 'source=test', got 'source=%s'", val)
	}
	if val, ok := parsedMeta["user"]; !ok || val != "gemini" {
		t.Errorf("Expected metadata 'user=gemini', got 'user=%s'", val)
	}
}

func TestBlobBuilder_MaxBytes(t *testing.T) {
	input := "This is a long string that should be truncated."
	expectedRead := "This is a " // 10 bytes
	expectedBody := base64.StdEncoding.EncodeToString([]byte(expectedRead))

	builder := blob2json.BlobBuilder{
		ContentType: "text/plain",
		MaxBytes:    10,
	}

	blob, err := builder.NewBlobFromReader(strings.NewReader(input), "truncated.txt")
	if err != nil {
		t.Fatalf("Failed to create blob: %v", err)
	}

	if blob.Body != expectedBody {
		t.Errorf("Expected truncated body '%s', got '%s'", expectedBody, blob.Body)
	}
	if *blob.ContentLength != 10 {
		t.Errorf("Expected truncated content length 10, got %d", *blob.ContentLength)
	}
}

func TestBlobBuilder_BinaryData(t *testing.T) {
	binaryData := []byte{0xDE, 0xAD, 0xBE, 0xEF, 0x01, 0x02, 0x03, 0x04}
	expectedBody := base64.StdEncoding.EncodeToString(binaryData)

	builder := blob2json.BlobBuilder{
		ContentType: "application/octet-stream",
		MaxBytes:    1024,
	}

	blob, err := builder.NewBlobFromReader(bytes.NewReader(binaryData), "binary.bin")
	if err != nil {
		t.Fatalf("Failed to create blob: %v", err)
	}

	if blob.Name != "binary.bin" {
		t.Errorf("Expected name 'binary.bin', got '%s'", blob.Name)
	}
	if blob.Body != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, blob.Body)
	}
}

func TestBlobBuilder_WithTypedMetadata(t *testing.T) {
	input := "some typed data"
	now := time.Now().Truncate(time.Second)

	builder := blob2json.BlobBuilder{
		ContentType:  "text/plain",
		MaxBytes:     1024,
		LastModified: &now,
	}

	blob, err := builder.NewBlobFromReader(strings.NewReader(input), "typed.txt")
	if err != nil {
		t.Fatalf("Failed to create blob: %v", err)
	}

	if *blob.ContentLength != int64(len(input)) {
		t.Errorf("Expected content length %d, got %d", len(input), *blob.ContentLength)
	}

	if blob.LastModified == nil {
		t.Fatal("Expected last modified to be set, but it was nil")
	}

	if !blob.LastModified.Equal(now) {
		t.Errorf("Expected last modified to be %v, got %v", now, *blob.LastModified)
	}

	// Also check JSON marshaling
	marshaled, err := json.Marshal(blob)
	if err != nil {
		t.Fatalf("Failed to marshal blob: %v", err)
	}

	var unmarshaled map[string]json.RawMessage
	if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal to map: %v", err)
	}

	if _, ok := unmarshaled["last_modified"]; !ok {
		t.Error("Expected 'last_modified' field in JSON output, but it was missing")
	}

	if _, ok := unmarshaled["content_length"]; !ok {
		t.Error("Expected 'content_length' field in JSON output, but it was missing")
	}
}
