package ipcUnix

import (
	"net"
	// "os"
	// "os/signal"
	// "syscall"
	log "github.com/Sirupsen/logrus"

)

type CallbackFunc func(name string) string

type IPCServer struct {
	sockFile string
	// 外部注册函数
	Call CallbackFunc
	lnSock net.Listener
}


func echoServer(this *IPCServer, c net.Conn) {
	// for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}
		data := string(buf[0:nr])
		log.Debugf("Server got:", data)

		msg := this.Call(data)
		_, err = c.Write([]byte(msg))
		if err != nil {
			log.Errorf("Writing client error: %v", err)
		}
		c.Close()
	// }
}

func NewServer(sock string) *IPCServer {
	instance := new(IPCServer)
	instance.sockFile = sock
	instance.Call = nil
	return instance
}

func (this *IPCServer)RegsiterCallback(call CallbackFunc) error {
	if this.Call != nil {
		log.Errorf("cannot regsiter again!!!!")
		return nil
	}
	this.Call = call
	return nil
}

func (this *IPCServer)Close() error {
	// pid := syscall.Getppid()
	// log.Infof("main: Killing pid: %v", pid)
	// syscall.Kill(pid, syscall.SIGTERM)
	this.lnSock.Close()
	return nil
}

func (this *IPCServer)CreateServer() error {
	log.Debugf("create begin!")
	ln,err := net.Listen("unix",this.sockFile)
	if err != nil {
		log.Fatalf("Listen error:%v",err)
		return err
	}
	this.lnSock = ln
	// sigc := make(chan os.Signal, 1)
	// signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	// go func(ln net.Listener, c chan os.Signal) {
	// 	sig := <-c
	// 	log.Infof("Caught signal %s: shutting down.", sig)
	// 	ln.Close()
	// 	os.Exit(0)
	// }(ln, sigc)

	for {
		fd, err := ln.Accept()
		if err != nil {
			log.Errorf("Accept error: ", err)
			continue
		}

		go echoServer(this,fd)
	}

	return nil
}
