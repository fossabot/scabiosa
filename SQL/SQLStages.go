package SQL

import "fmt"

type SQLStage int64

const (
	SqlStageCompress   = 1
	SqlStageUpload     = 2
	SqlStageFinialzing = 3
)

func (e SQLStage) String() string {
	switch e {
	case SqlStageCompress:
		return "COMPRESS"
	case SqlStageUpload:
		return "UPLOAD"
	case SqlStageFinialzing:
		return "FINALIZING"
	default:
		return fmt.Sprintf("%d", e)
	}
}
