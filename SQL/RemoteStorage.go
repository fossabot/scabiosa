package SQL

import "fmt"

type RemoteStorageType int64

const (
	RemoteAzureFile = 1
	RemoteNone      = 2
)

func (e RemoteStorageType) String() string {
	switch e {
	case RemoteAzureFile:
		return "AZURE-FILE"
	case RemoteNone:
		return "LOCAL"
	default:
		return fmt.Sprintf("%d", e)
	}
}
