package auth

type BlobStore interface {
	GetBlob(id string) ([]byte, error)
	PutBlob(id string, val []byte) error
	DelBlob(id string) error
}
