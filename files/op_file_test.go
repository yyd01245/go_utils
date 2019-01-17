package files_test


import (
	"fmt"
	"testing"
	"github.com/yyd01245/go_utils/files"

)

func TestFile(t *testing.T) {
	const localFile = "./test.log"
	err := files.WriteStringToFile(localFile,"hello world!")
	if err != nil {
		fmt.Println("-- error write to file")
	}
	if files.CheckFileIsExist(localFile) {
		fmt.Println("--file exit")
	}

	output,err := files.ReadFileAll(localFile)
	fmt.Println("read: ",output)

	const pidFile = "./test.pid"
	err = files.WritePidToFile(pidFile)
	if err != nil {
		fmt.Println("-- error WritePidToFile to file")
	}
	output,err = files.ReadFileAll(pidFile)
	fmt.Println("pid file read: ",output)

	err = files.CheckPidFromFile(pidFile)
	if err != nil {
		fmt.Println("-- error CheckPidFromFile to file")
	}

}