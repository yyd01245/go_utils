package httpUtil

import (
	"fmt"
	"io"
	URL "net/url"
	"strings"
	"sync"
	// "fmt"
	"errors"
	// "encoding/json"
	"io/ioutil"
	"net"
	// "crypto/tls"
	// "crypto/x509"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string, downtime int) error {

	if downtime == 0 {
		downtime = 10
	}
	timeout := time.Duration(downtime) * time.Second
	client := http.Client{
		Timeout: timeout,
	}
	// Get the data
	// resp, err := http.Get(url)
	resp, err := client.Get(url)
	if err != nil {
		log.Warnf("http get url:%v failed: %v ", url, err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		// Create the file
		out, err := os.Create(filepath)
		if err != nil {
			log.Warnf("create file:%v failed ", filepath)
			return err
		}
		defer out.Close()
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Warnf("copy file:%v failed ", filepath)
			return err
		}
	} else {
		txt := fmt.Sprintf("get http response code:%v  ", resp.StatusCode)
		log.Warnf(txt)
		return errors.New(txt)
	}

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 		// handle error
	// 		log.Warnf("get body error %v",body)
	// }
	// log.Infof("body: %v",body)
	// Write the body to file

	return nil
}

type HttpHandle struct {
	mux        sync.Mutex
	httpClient *http.Client
	httpTR     *http.Transport
}

type ResponData struct {
	BodyData   string
	StatusCode int
}

const TIMEOUTE = 15

func NewhttpHandle() *HttpHandle {
	var instance = new(HttpHandle)
	instance.httpTR = nil

	instance.httpTR = &http.Transport{

		// TLSClientConfig: &tls.Config{InsecureSkipVerify : true},
		DisableKeepAlives: true,
		Dial: (&net.Dialer{
			Timeout: TIMEOUTE * time.Second,
			// Deadline:  time.Now().Add(TIMEOUTE  * time.Second),
			// KeepAlive:  * time.Second,
		}).Dial,
		TLSHandshakeTimeout: TIMEOUTE * time.Second,
		DisableCompression:  true,
	}
	// timeout := time.Duration(TIMEOUTE * time.Second)

	// instance.httpClient = &http.Client{
	// 	Timeout: timeout,
	// 	Transport: instance.httpTR,
	// }
	return instance
}

func (this *HttpHandle) ResetHttpClient() {

	this.httpTR = &http.Transport{
		// TLSClientConfig: &tls.Config{InsecureSkipVerify : true},

		DisableKeepAlives: true,
		Dial: (&net.Dialer{
			Timeout: TIMEOUTE * time.Second,
			// Deadline:  time.Now().Add(TIMEOUTE  * time.Second),
			// KeepAlive:  * time.Second,
		}).Dial,
		TLSHandshakeTimeout: TIMEOUTE * time.Second,
		DisableCompression:  true,
	}
	log.Warnf("ResetHttpClient http")
	// timeout := time.Duration(TIMEOUTE * time.Second)

	// this.httpClient = &http.Client{
	// 	Timeout: timeout,
	// 	Transport: this.httpTR,
	// }
	// log.Infof("ResetHttpClient")
}

func (this *HttpHandle) PostRequest(url string, data map[string]string) (*ResponData, error) {
	//var jsonStr = []byte(json)
	form := URL.Values{}
	for key, value := range data {

		// log.Debugf("map key=%v,value=%v",key,value)
		form.Add(key, value)
	}
	// log.Debugf("post : %v",form)
	respData := &ResponData{
		BodyData:   "",
		StatusCode: 0,
	}

	// if this.httpTR == nil {
	// 	this.ResetHttpClient()
	// }
	timeout := time.Duration(TIMEOUTE * time.Second)

	// httpClient := this.httpClient
	this.mux.Lock()
	httpClient := &http.Client{
		Timeout:   timeout,
		Transport: this.httpTR,
	}
	log.Debugf("--------begin %p, lock url:%v=--------", this, url)
	// defer this.mux.Unlock()
	req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := httpClient.Do(req)
	defer func() {
		// log.Debugf("--------begin %p, unlock url:%v=--------",this,url)
		this.mux.Unlock()
		req.Close = true
	}()
	if err != nil {
		log.Errorf("PostRequest unable ro reach the serve:%v, err:%v", url, err)
		// this.ResetHttpClient()
		return respData, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	respData.BodyData = string(body)
	respData.StatusCode = resp.StatusCode
	log.Debugf("resp=%v", respData)
	return respData, err
}

func (this *HttpHandle) GetHttpDataWitchCode(url string) (*ResponData, error) {

	data := &ResponData{
		BodyData:   "",
		StatusCode: 0,
	}

	// if this.httpTR == nil {
	// 	this.ResetHttpClient()
	// }
	timeout := time.Duration(TIMEOUTE * time.Second)

	// httpClient := this.httpClient
	this.mux.Lock()
	httpClient := &http.Client{
		Timeout:   timeout,
		Transport: this.httpTR,
	}
	// log.Debugf("--------begin %p, lock url:%v=--------",this,url)
	// defer this.mux.Unlock()
	req, err := http.NewRequest("GET", url, nil)
	resp, err := httpClient.Do(req)
	defer func() {
		// log.Debugf("--------begin %p, unlock url:%v=--------",this,url)
		this.mux.Unlock()
		req.Close = true
	}()
	if err != nil {
		// log.Errorf("GetHttpDataWitchCode unable ro reach the serve:%v, err:%v",url,err)
		// this.ResetHttpClient()
		return data, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	data.BodyData = string(body)
	data.StatusCode = resp.StatusCode
	// log.Debugf("respon data=%v", data)
	// txt := fmt.Sprintf("get http response code:%v",resp.StatusCode)
	// log.Warnf(txt)
	return data, err
}

func (this *HttpHandle) GetHttpData(url string) (string, error) {

	data := ""

	// if this.httpTR == nil {
	// 	this.ResetHttpClient()
	// }
	timeout := time.Duration(TIMEOUTE * time.Second)

	// httpClient := this.httpClient
	this.mux.Lock()
	httpClient := &http.Client{
		Timeout:   timeout,
		Transport: this.httpTR,
	}
	log.Debugf("--------begin %p, lock url:%v=--------", this, url)
	// defer this.mux.Unlock()
	req, err := http.NewRequest("GET", url, nil)
	resp, err := httpClient.Do(req)
	defer func() {
		log.Debugf("--------begin %p, unlock url:%v=--------", this, url)
		this.mux.Unlock()
		req.Close = true
	}()
	if err != nil {
		log.Errorf("GetHttpData unable ro reach the serve:%v, err:%v", url, err)
		this.ResetHttpClient()
		return data, err
	}
	defer resp.Body.Close()
	// body, _ := ioutil.ReadAll(resp.Body)
	// data = string(body)
	// log.Debugf("body=%v", data)
	// // txt := fmt.Sprintf("get http response code:%v",resp.StatusCode)
	// // log.Warnf(txt)
	// return data, err
	if resp.StatusCode == 200 || resp.StatusCode == 404 {
		body, _ := ioutil.ReadAll(resp.Body)
		data = string(body)
		log.Debugf("body=%v", data)
		err = nil
	} else {
		txt := fmt.Sprintf("get http response code:%v", resp.StatusCode)
		log.Warnf(txt)
		err = errors.New(txt)
	}

	return data, err
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFileOperatorError(filepath string, url string) error {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		log.Warnf("create file:%v failed ", filepath)
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		log.Warnf("http get url:%v failed: %v ", url, err)
		resCode := resp.StatusCode
		if resCode != 404 {
			// 404 not found
			return err
		}
		log.Warnf("respone code status : %v", resCode)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		log.Warnf("get body error %v", body)
	}
	log.Infof("body: %v", body)
	// Write the body to file

	out.WriteString(string(body))
	return nil
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadOkFile(filepath string, url string, downtime int) error {

	if downtime == 0 {
		downtime = 10
	}
	timeout := time.Duration(downtime) * time.Second
	client := http.Client{
		Timeout: timeout,
	}
	// Get the data
	// resp, err := http.Get(url)
	resp, err := client.Get(url)
	if err != nil {
		log.Warnf("http get url:%v failed: %v ", url, err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		// Create the file
		out, err := os.Create(filepath)
		if err != nil {
			log.Warnf("create file:%v failed ", filepath)
			return err
		}
		defer out.Close()
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Warnf("copy file:%v failed ", filepath)
			return err
		}
	} else {
		txt := fmt.Sprintf("get http response code:%v  ", resp.StatusCode)
		log.Warnf(txt)
		return errors.New(txt)
	}

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 		// handle error
	// 		log.Warnf("get body error %v",body)
	// }
	// log.Infof("body: %v",body)
	// Write the body to file

	return nil
}
