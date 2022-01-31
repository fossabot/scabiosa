package SQL

import "fmt"

type RemoteStorageType int64

const (
	REMOTE_AZURE_FILE = 1
	REMOTE_AZURE_BLOB = 2
	REMOTE_NONE       = 3
)

func (e RemoteStorageType) String() string {
	switch e {
	case REMOTE_AZURE_FILE:
		return "AZURE-FILE"
	case REMOTE_AZURE_BLOB:
		return "AZURE-BLOB"
	case REMOTE_NONE:
		return "NONE"
	default:
		return fmt.Sprintf("%d", e)
	}
}
