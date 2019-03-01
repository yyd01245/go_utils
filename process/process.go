package process

import(
	"syscall"
	// log "github.com/Sirupsen/logrus"

)

func FindProcess(pid int) error{
	err := syscall.Kill(pid, 0);
	if err == nil {
		// log.Debugf("find process success ")

	}else {
		// log.Infof("Failed to find process: %v\n", err)
		return err	
	}
	return nil
}

func KillProcess(pid int) error {
	err := syscall.Kill(pid, 9);
	if err == nil {
		// log.Debugf("kill process success ")
	}else {
		// log.Warnf("Failed to kill process: %v\n", err)
		return err
	}
	return nil
}