package ipcUnix

import (
	"net"
	"io"
	"time"
	log "github.com/Sirupsen/logrus"

)

type IPCClient struct {
	sockFile string
	conn 	net.Conn
}


func reader(r io.Reader) string{
	buf := make([]byte, 1024)
	data := ""
	// for {
		n, err := r.Read(buf[:])
		if err != nil {
			return data
		}
		log.Debugf("Client got:%v", string(buf[0:n]))
		data = data + string(buf[0:n])
	// }
	return data
}

func writer(c io.Writer, msg string) error {
	_,err := c.Write([]byte(msg))
	if err != nil {
		log.Errorf("write error:%v",err)
		return err
	}

	return nil
}

func NewClient(sock string) *IPCClient {
	instance := new(IPCClient)
	instance.sockFile = sock

	return instance
}

func (this *IPCClient) Write(msg string) error{
	return writer(this.conn,msg)
}

func (this *IPCClient) Read() string{
	return reader(this.conn)
}

func (this *IPCClient) Close() {

	this.conn.Close()
}


func (this *IPCClient)CreateClient() error {
	log.Debugf("create client begin!")
	// ln,err := net.Dial("unix",this.sockFile)
	timeout := time.Duration(5 * time.Second)
	ln,err := net.DialTimeout("unix",this.sockFile,timeout)
	if err != nil {
		log.Fatalf("Listen error:%v",err)
		return err
	}
	this.conn = ln

	return nil
}
