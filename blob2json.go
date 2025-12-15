package blob2json

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"time"
)

// Blob represents the structure for the JSON output.
type Blob struct {
	Name                    string          `json:"name"`
	ContentType             string          `json:"content_type"`
	ContentEncoding         string          `json:"content_encoding"`
	ContentTransferEncoding string          `json:"content_transfer_encoding"`
	Body                    string          `json:"body"`
	Metadata                json.RawMessage `json:"metadata,omitempty"`

	// Typed metadata
	ContentLength *int64     `json:"content_length,omitempty"`
	LastModified  *time.Time `json:"last_modified,omitempty"`
}

// BlobBuilder helps construct a Blob.
type BlobBuilder struct {
	ContentType     string
	ContentEncoding string
	MaxBytes        int64
	Metadata        map[string]string
	LastModified    *time.Time
}

// NewBlobFromReader creates a Blob by reading from an io.Reader.
func (b *BlobBuilder) NewBlobFromReader(r io.Reader, name string) (*Blob, error) {
	reader := io.LimitReader(r, b.MaxBytes)
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	encodedBody := base64.StdEncoding.EncodeToString(data)
	contentLength := int64(len(data))

	blob := &Blob{
		Name:                    name,
		ContentType:             b.ContentType,
		ContentEncoding:         b.ContentEncoding,
		ContentTransferEncoding: "base64",
		Body:                    encodedBody,
		ContentLength:           &contentLength,
		LastModified:            b.LastModified,
	}

	if len(b.Metadata) > 0 {
		metaJSON, err := json.Marshal(b.Metadata)
		if err != nil {
			return nil, err
		}
		blob.Metadata = metaJSON
	}

	return blob, nil
}
