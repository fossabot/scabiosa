package Tools

import (
	// skipcq: GSC-G501
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"os"
)

type HashType int64

const (
	SHA256 HashType = iota + 1
	MD5
)

func openFileAndCalcHash(fileData *[]byte, filePath string) error {

	var err error
	*fileData, err = os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return nil
}

func CalculateHashValue(filePath string, hashType HashType) (string, error) {
	var data []byte

	err := openFileAndCalcHash(&data, filePath)
	if err != nil {
		return "", fmt.Errorf("hash calulation failed. Error: %s", err)
	}

	switch hashType {
	case SHA256:
		return fmt.Sprintf("%x", sha256.Sum256(data)), nil
	case MD5:
		// skipcq:GSC-G401, GO-S1023
		return fmt.Sprintf("%x", md5.Sum(data)), nil
	default:
		return "", fmt.Errorf("you shouldn't be here. [Code 328]")
	}
}

func GetHashTypeFromString(name string) HashType {
	switch name {
	case HashType(SHA256).String():
		return SHA256
	case HashType(MD5).String():
		return MD5
	default:
		return -1
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
