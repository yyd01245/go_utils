package files

import (
	"encoding/hex"
	"crypto/md5"
	// "crypto/aes"
	// "encoding/base64"
	// "crypto/cipher"
	"os"
	"io"
	"io/ioutil"
	"strings"
	"errors"
	log "github.com/Sirupsen/logrus"	
)


func GetFileMD5Value(filename string ) string {
	f, err := os.Open(filename)
	if err != nil {
			log.Warnf("Open", err)
			return ""
	}

	defer f.Close()

	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
			log.Warnf("Copy", err)
			return ""
	}
	ret := md5hash.Sum(nil)
	return hex.EncodeToString(ret)

}

func GetFileValue(filename string ) []byte {
	f, err := os.Open(filename)
	if err != nil {
			log.Warnf("Open", err)
			return nil
	}

	defer f.Close()

	body, err := ioutil.ReadAll(f)
	if err != nil {
		log.Warnf("ReadAll", err)
			return nil
	}
	return body

}

func GetByteMD5Value(data []byte) string {
	md5hash := md5.New()
	md5hash.Write(data)

	ret := md5hash.Sum(nil)
	return hex.EncodeToString(ret)

}

func GetStringMD5Value(data string ) string {

	md5hash := md5.New()
	md5hash.Write([]byte(data))

	ret := md5hash.Sum(nil)
	return hex.EncodeToString(ret)

}


func ParseFileToMap(filename string, mapData map[string]string) error{
	// key=value
	txtTmpCloudConfig := GetFileValue(filename)

	cloudConfig := string(txtTmpCloudConfig)
	log.Debugf("get file string :%v",cloudConfig)
	if cloudConfig == "" {
		log.Errorf("ParseFileToMap get file string error")
		return errors.New("ParseFileToMap get file string error")
	}
	configLine := strings.Split(cloudConfig,"\n")

	for _,value := range configLine {
		value = strings.Replace(value, " ", "", -1)
		value = strings.Replace(value, "\n", "", -1)
		if value == "" {
			continue
		}
		data := strings.Split(value,"=")
		if len(data) != 2 {
			log.Errorf("parse key value failed:%v",value)
			continue
		}
		mapData[data[0]] = data[1]
	}
	return nil
}

// func AESEncrypt(src []byte, key []byte) (encrypted []byte) {
// 	cipher, _ := aes.NewCipher(generateKey(key))
// 	length := (len(src) + aes.BlockSize) / aes.BlockSize
// 	plain := make([]byte, length*aes.BlockSize)
// 	copy(plain, src)
// 	pad := byte(len(plain) - len(src))
// 	for i := len(src); i < len(plain); i++ {
// 		plain[i] = pad
// 	}
// 	encrypted = make([]byte, len(plain))

// 	for bs, be := 0, cipher.BlockSize(); bs <= len(src); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
// 		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
// 	}

// 	return encrypted
// }

// func generateKey(key []byte) (genKey []byte) {
// 	genKey = make([]byte, 16)
// 	copy(genKey, key)
// 	for i := 16; i < len(key); {
// 		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
// 			genKey[j] ^= key[i]
// 		}
// 	}
// 	return genKey
// }

// func AESDecrypt(encrypted []byte, key []byte) (decrypted []byte) {
// 	cipher, _ := aes.NewCipher(generateKey(key))
// 	decrypted = make([]byte, len(encrypted))
// 	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
// 		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
// 	}
// 	log.Infof("---decrpted:%v, len(decrpted)=%v",decrypted,len(decrypted))
// 	trim := 0
// 	if len(decrypted) > 0 {
// 		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
// 	}
// 	log.Infof("---trim:%v",trim)
// 	return decrypted[:trim]
// }

// func Decrypt(key []byte, securemess string) (decodedmess string,err error){
// 	cipherText, err := base64.URLEncoding.DecodeString(securemess)
// 	if err != nil {
// 		log.Errorf("decode string cipherText error :%v",err)
// 		return
// 	}

// 	block, err := aes.NewCipher(generateKey(key))
// 	if err != nil {
// 		log.Errorf("NewCipher error :%v",err)
// 		return
// 	}

// 	if len(cipherText) < aes.BlockSize {
// 		err = errors.New("Ciphertext block size is too short!")
// 		log.Errorf(" error :%v",err)
// 		return
// 	}

// 	//IV needs to be unique, but doesn't have to be secure.
// 	//It's common to put it at the beginning of the ciphertext.
// 	iv := cipherText[:aes.BlockSize]
// 	cipherText = cipherText[aes.BlockSize:]
// 	log.Infof("in:%v",iv)
// 	log.Infof("text: %v",cipherText)
// 	stream := cipher.NewCFBDecrypter(block, iv)
// 	// XORKeyStream can work in-place if the two arguments are the same.
// 	stream.XORKeyStream(cipherText, cipherText)

// 	decodedmess = string(cipherText)
// 	return
// }