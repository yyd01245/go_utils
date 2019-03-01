package execUtil

import (
	// "fmt"
	"fmt"
	// "errors"
	"strings"
	"os"
	"os/exec"
	"bytes"
	// log "github.com/Sirupsen/logrus"
)


func ExecScripts(binPath string,args []string) (string, error) {
	stderr := bytes.NewBuffer(nil)
	stdout := bytes.NewBuffer(nil)

	cmd := exec.Command(binPath, args...)
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	err := cmd.Run()
	if err != nil {
		// log.Debugf("exec command: %v  args: %v err=%v",binPath,
				// args,err)
		return "",err
	}
	// log.Infof("exec command: ",binPath," args: %v",
	// 	args," success")	
	outputErr := string(stderr.Bytes())
	if len(outputErr) > 0 {
		// log.Debugf("exec command: stderr: %v",outputErr)
	}
	output := string(stdout.Bytes())
	if len(output) > 0 {
		// log.Debugf("exec command: stdout: ",output)
		return "",err
	}	
	return output,err
}

func ExecScriptsPipe(binPath string,args []string) (string, error) {
	stderr := bytes.NewBuffer(nil)
	stdout := bytes.NewBuffer(nil)
	strArgs := ""
	for _,value := range args{
		strArgs += (" " + value)
	}
	input := fmt.Sprintf(`%s%s`,binPath,strArgs)
	cmd := exec.Command("/bin/sh", "-c", input)
	// log.Infof("pipe exec input:%s",input)
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	err := cmd.Run()
	if err != nil {
		// log.Warnf("exec command: /bin/sh -c %v  args: %v err=%v",binPath,
		// 	input,err)
		return "",err
	}
	// log.Infof("exec command: /bin/sh -c  ",binPath," args: %v",
	// 	input," success")	
	outputErr := string(stderr.Bytes())
	if len(outputErr) > 0 {
		// log.Debugf("exec command: stderr: %v",outputErr)
	}
	output := string(stdout.Bytes())
	if len(output) > 0 {
		// log.Debugf("exec command: stdout: ",output)

	}	
	return output,err
}

func ExecScriptsByEnv(binPath string,args []string, envsKey string, envsValue string) (string, error) {
	stderr := bytes.NewBuffer(nil)
	stdout := bytes.NewBuffer(nil)

	cmd := exec.Command(binPath, args...)
	cmd.Stderr = stderr
	cmd.Stdout = stdout

	BaseEnv := os.Environ()
	cmdEnv := []string{}
	for _,env := range BaseEnv {
		index := strings.Index(env,"=")
		if index > 0 && env[:index] == envsKey {
			env = env + ":" + envsValue
		}
		cmdEnv = append(cmdEnv, env)
	}
	cmd.Env = cmdEnv

	err := cmd.Run()
	if err != nil {
		// log.Debugf("exec command: %v  args: %v err=%v",binPath,
		// 		args,err)
		return "",err
	}
	// log.Infof("exec command: ",binPath," args: %v",
	// 	args," success")	
	outputErr := string(stderr.Bytes())
	if len(outputErr) > 0 {
		// log.Debugf("exec command: stderr: %v",outputErr)
	}
	output := string(stdout.Bytes())
	if len(output) > 0 {
		// log.Debugf("exec command: stdout: ",output)

	}	
	return output,err
}
