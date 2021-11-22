package SQL

type SQLStage int64

const(
	SQLStage_Compress = 1
	SQLStage_Upload = 2
	SQLStage_DeleteTmp = 3
)
