package uciUtil

import (
	// "fmt"
	"strings"

	"github.com/yyd01245/go_utils/execUtil"
	// log "github.com/Sirupsen/logrus"
)

func GetUciValue(binPath string,args []string) (string, error) {

	outString,err := execUtil.ExecScripts(binPath,args)
	if err != nil {
		// log.Warnf("get uci %v ,error",args)
		return "",err
	}
	outString = strings.Replace(outString, "\n", "", -1) 
	// log.Debugf("uci show data: %v",outString)
	// outString = strings.Replace(outString, " ", "", -1) 
	// log.Debugf("uci show data: %v",outString)
	index := strings.Index(outString,"=")
	// output := strings.Split(outString[index+1:],"'")
	output := strings.Replace(outString[index+1:], "'", "", -1) 

	// log.Debugf("get data %v,len=%v",output,len(output))
	return output,nil

	// log.Debugf("output[0]=%v",output[0]);
	// return output[0],nil

}
func GetUciValueList(binPath string,args []string) ([]string, error) {

	outString,err := execUtil.ExecScripts(binPath,args)
	if err != nil {
		// log.Warnf("get uci %v ,error",args)
		return []string{},err
	}
	outString = strings.Replace(outString, "\n", "", -1) 
	index := strings.Index(outString,"=")
	outputstring := strings.Replace(outString[index+1:], "'", "", -1) 
	output := strings.Split(outputstring," ")
	// log.Debugf("get data %v,len=%v",output,len(output))
	// if len(output) < 3 {
	// 	txt := fmt.Sprintf("get uci %v, result:%s, parse error",args,outString)
	// 	return "",errors.New(txt)
	// }
	return output,nil

}

func SetUciValue(value string) error{
	// txt := fmt.Sprintf("%s=%s",key,value)
	args := []string{"set",value}
	_,err := execUtil.ExecScripts("uci",args)
	if err != nil {
		// log.Errorf("set uci %v ,error",args)
		return err
	}
	return nil
}

func CommitUciValue(value string) error{
	// txt := fmt.Sprintf("%s=%s",key,value)
	args := []string{"commit",value}
	_,err := execUtil.ExecScripts("uci",args)
	if err != nil {
		// log.Errorf("set uci %v ,error",args)
		return err
	}
	return nil
}
