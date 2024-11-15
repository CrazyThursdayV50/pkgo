package compressor

type Compressor interface {
	Compress(in []byte) (out []byte, err error)
	Uncompress(in []byte) (out []byte, err error)
}
