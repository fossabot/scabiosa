package SQL

import "fmt"

type SQLStage int64

const (
	SQLStage_Compress  = 1
	SQLStage_Upload    = 2
	SQLStage_DeleteTmp = 3
)

func (e SQLStage) String() string {
	switch e {
	case SQLStage_Compress:
		return "COMPRESS"
	case SQLStage_Upload:
		return "UPLOAD"
	case SQLStage_DeleteTmp:
		return "DELETE TMP"
	default:
		return fmt.Sprintf("%d", e)
	}
}
