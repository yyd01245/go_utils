package files

import(
	"fmt"
	"os"
	"errors"
	"io/ioutil"
	"strings"
	"strconv"
	log "github.com/Sirupsen/logrus"
	"github.com/yyd01245/go_utils/process"
)

// CheckFileIsExist 检测文件是否存在
func CheckFileIsExist(filename string) bool{
	var exist =true
	if _,err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	log.Debug("file:",filename," exist=",exist)

	return exist
}

//WriteListLineToFile 写入字符串数组，每行一个
func WriteListLineToFile(filename string, paramList[]string) error{

	if len(paramList) == 0 {
		log.Warnf("write firle paramList len =0")
		return errors.New("write firle paramList len =0")
	}
	file, err := os.OpenFile(filename,os.O_CREATE|os.O_TRUNC|os.O_WRONLY,0644)
	if err != nil {
		log.Warnf("open file failed !")
		return err
	}
	defer file.Close()

	for k,v := range paramList {
		log.Debug("key=",k," value=",v)
		text := fmt.Sprintf("%s\n",v)
		log.Debug("ready to write file=",filename," text=",text)
		_,err := file.WriteString(text)
		if err != nil {
			return err
		}
	}
	return nil
}
// WriteStringToFile 写入字符串到文件
func WriteStringToFile(filename string, text string) error{

	file, err := os.OpenFile(filename,os.O_CREATE|os.O_TRUNC|os.O_WRONLY,0644)
	if err != nil {
		log.Warnf("open file failed !")
		return err
	}
	defer file.Close()
	log.Debugf("ready to write file=",filename," text=",text)
	_,err = file.WriteString(text)
	if err != nil {
		return err
	}
	return nil
}
// AppendListLineToFile 最近字符串数组到文件，一行一个
func AppendListLineToFile(filename string, paramList[]string) error{

	if len(paramList) == 0 {
		log.Warnf("write firle paramList len =0")
		return errors.New("write firle paramList len =0")
	}
	file, err := os.OpenFile(filename,os.O_CREATE|os.O_RDWR|os.O_APPEND,0644)
	if err != nil {
		log.Warnf("open file failed !")
		return err
	}
	defer file.Close()

	for k,v := range paramList {
		log.Debug("key=",k," value=",v)
		text := fmt.Sprintf("%s\n",v)
		log.Debug("ready to write file=",filename," text=",text)
		_,err := file.WriteString(text)
		if err != nil {
			log.Errorf("write string:%v Error:%v",text,err)
			return err
		}
	}
	return nil
}
//ReadFileAll 读取文件，以字符串返回
func ReadFileAll(filename string) (string,error) {

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "",err
	}
	return string(b),nil
}
//CreateDir 创建目录   
func CreateDir(path string) error {

	return os.MkdirAll(path,0744)

}
//WritePidToFile 写入当前 pid 到文件中 
func WritePidToFile(filename string) error {
	file, err := os.OpenFile(filename,os.O_CREATE|os.O_TRUNC|os.O_WRONLY,0644)
	if err != nil {
		log.Warnf("open file failed !")
		return err
	}
	defer file.Close()
	txt := fmt.Sprintf("pid=%d",os.Getpid())
	log.Debugf(txt)
	file.WriteString(txt)
	return nil
}

// CheckPidFromFile return nil success, error no pid
func CheckPidFromFile(filename string) error {
	file, err := os.OpenFile(filename,os.O_RDONLY,0644)
	if err != nil {
		log.Warnf("open file failed !")
		return err
	}
	defer file.Close()
	body, err := ioutil.ReadAll(file)
	if err != nil {
		log.Warnf("ReadAll", err)
			return err
	}
	data := strings.Split(string(body),"=")
	log.Debugf("get pid file: %v",data)
	pid,err := strconv.Atoi(data[1])
	log.Debugf("get pid = %d",pid)

	return process.FindProcess(pid)
}

