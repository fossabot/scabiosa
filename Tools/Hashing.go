package Compressor

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
)

type HashType int64

const (
	SHA256 HashType = iota
	MD5
)

func calculateHashValue(data []byte, hashType HashType) string {
	switch hashType {
	case SHA256:
		return fmt.Sprintf("%x", sha256.Sum256(data))
	case MD5:
		return fmt.Sprintf("%x", md5.Sum(data))
	default:
		return ""
	}
}

func (e HashType) String() string {
	switch e {
	case SHA256:
		return "SHA256"
	case MD5:
		return "MD5"
	default:
		return fmt.Sprintf("%d", e)
	}
}
