package zlib

import (
	"bytes"
	"compress/zlib"
	"io"

	"github.com/CrazyThursdayV50/pkgo/websocket/compressor"
)

type zlibCompressor struct {
	compressLevel int
}

// level between -2 ~ 9
func (c *zlibCompressor) Compress(in []byte) ([]byte, error) {
	var b bytes.Buffer
	w, err := zlib.NewWriterLevel(&b, c.compressLevel)
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

func (c *zlibCompressor) Uncompress(in []byte) ([]byte, error) {
	reader := bytes.NewReader(in)

	r, err := zlib.NewReader(reader)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	_, _ = io.Copy(&b, r)
	r.Close()

	return b.Bytes(), nil
}

func NewZlibCompressor(compressLevel int) compressor.Compressor {
	return &zlibCompressor{compressLevel: compressLevel}
}
