package SQL

import "fmt"

type SQLStage int64

const (
	SqlstageCompress = 1
	SqlstageUpload   = 2
)

func (e SQLStage) String() string {
	switch e {
	case SqlstageCompress:
		return "COMPRESS"
	case SqlstageUpload:
		return "UPLOAD"
	default:
		return fmt.Sprintf("%d", e)
	}
}
