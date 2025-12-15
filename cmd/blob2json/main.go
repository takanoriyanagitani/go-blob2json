package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/takanoriyanagitani/go-blob2json"
)

type metadataFlag map[string]string

func (m *metadataFlag) String() string {
	var pairs []string
	for k, v := range *m {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(pairs, ",")
}

func (m *metadataFlag) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid metadata format; expected key=value")
	}
	(*m)[parts[0]] = parts[1]
	return nil
}

func run() error {
	var name, contentType, contentEncoding, lastModifiedStr string
	var maxBytes int64
	metadata := make(metadataFlag)

	flag.StringVar(&name, "name", "", "The name of the blob (e.g., 'blob.dat')")
	flag.StringVar(&contentType, "content-type", "application/octet-stream", "The content type of the blob")
	flag.StringVar(&contentEncoding, "content-encoding", "", "The content encoding of the blob")
	flag.StringVar(&lastModifiedStr, "last-modified", "", "The last modified time in RFC3339 format")
	flag.Int64Var(&maxBytes, "max-bytes", 1048576, "The maximum number of bytes to read from stdin")
	flag.Var(&metadata, "metadata", "Metadata as key=value pairs (can be specified multiple times)")

	flag.Parse()

	if name == "" {
		return fmt.Errorf("--name is required")
	}

	var lastModified *time.Time
	if lastModifiedStr != "" {
		t, err := time.Parse(time.RFC3339, lastModifiedStr)
		if err != nil {
			return fmt.Errorf("unable to parse last-modified time: %w", err)
		}
		lastModified = &t
	}

	builder := blob2json.BlobBuilder{
		ContentType:     contentType,
		ContentEncoding: contentEncoding,
		MaxBytes:        maxBytes,
		Metadata:        metadata,
		LastModified:    lastModified,
	}

	blob, err := builder.NewBlobFromReader(os.Stdin, name)
	if err != nil {
		return fmt.Errorf("unable to create blob: %w", err)
	}

	encoder := json.NewEncoder(os.Stdout)
	return encoder.Encode(blob)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
