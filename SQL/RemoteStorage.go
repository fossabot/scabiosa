package SQL

type RemoteStorageType int64

const(
	REMOTE_AZURE_FILE = 1
	REMOTE_AZURE_BLOB = 2
	REMOTE_NONE = 3
)