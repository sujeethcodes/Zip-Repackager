package entity

type FileData struct {
	Name     string
	Content  []byte
	FileSize int64
	SHA256   [32]byte
}
