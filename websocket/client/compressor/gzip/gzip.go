package gzip

import (
	"bytes"
	"compress/gzip"
	"io"

	"github.com/CrazyThursdayV50/pkgo/websocket/client/compressor"
)

type gzipCompressor struct {
	compressLevel int
}

// level between -2 ~ 9
func (c *gzipCompressor) Compress(in []byte) ([]byte, error) {
	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, c.compressLevel)
	if err != nil {
		return nil, err
	}
	defer w.Close()

	_, err = w.Write(in)
	if err != nil {
		return nil, err
	}
	w.Flush()

	return b.Bytes(), nil
}

func (c *gzipCompressor) Uncompress(in []byte) ([]byte, error) {
	reader := bytes.NewReader(in)

	r, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	_, _ = io.Copy(&b, r)
	r.Close()

	return b.Bytes(), nil
}

func NewGzipCompressor(compressLevel int) compressor.Compressor {
	return &gzipCompressor{compressLevel: compressLevel}
}
