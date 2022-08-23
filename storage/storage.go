package fileStorage

import (
	"fmt"
	"os"
)

func MoveFileToFolder(fileName string) {
	os.Rename(fileName, fmt.Sprintf("dowload/%s", fileName))
}
