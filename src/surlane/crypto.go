package surlane

import (
	"crypto/cipher"
	"crypto/aes"
	"crypto/md5"
	"crypto/rand"
)

const (
	CES_128_CFB = iota
	CES_256_CFB
)

var (
	cipherMap = map[int]StreamBuilder{
		CES_128_CFB: &CES128StreamBuilder{},
		CES_256_CFB: nil,
	}
)

type StreamBuilder interface {
	createCipher(string, []byte)(*SurCipher, error)
	genIV()[]byte
	getIvLen()int
}

type CES128StreamBuilder struct {}

func (*CES128StreamBuilder) createCipher(password string, iv []byte) (surCipher *SurCipher, err error) {
	md5Sum := md5.Sum([]byte(password))
	block, err := aes.NewCipher(md5Sum[:])
	if err != nil {
		return
	}
	surCipher = &SurCipher{
		iv,
		password,
		CES_128_CFB,
		cipher.NewCFBEncrypter(block, iv),
	 	cipher.NewCFBDecrypter(block, iv),
	}
	return
}

func (*CES128StreamBuilder) genIV() (iv []byte) {
	iv = make([]byte, aes.BlockSize)
 	rand.Read(iv)
	return
}

func (*CES128StreamBuilder) getIvLen() int {
	return aes.BlockSize
}

type SurCipher struct {
	IV     []byte
	key    string
	method int16
	enc    cipher.Stream
	dec    cipher.Stream
}

func GetIV(config ClientConfig) []byte {
	streamBuilder := cipherMap[config.Method]
	return streamBuilder.genIV()
}

func GetIvLen(config ServerConfig) int {
	streamBuilder := cipherMap[config.Method]
	return streamBuilder.getIvLen()
}

func NewSurCipher4Client(config ClientConfig, iv []byte) (*SurCipher, error) {
	streamBuilder := cipherMap[config.Method]
	return streamBuilder.createCipher(config.Password, iv)
}

func NewSurCipher4Server(config ServerConfig, iv []byte) (*SurCipher, error) {
	return cipherMap[config.Method].createCipher(config.Password, iv)
}

func (surCipher *SurCipher) encrypt(data []byte) (cipherText []byte) {
	cipherText = make([]byte, len(data))
	surCipher.enc.XORKeyStream(cipherText, data)
	return
}

func (surCipher *SurCipher) decrypt(data []byte) (cipherText []byte) {
	cipherText = make([]byte, len(data))
	surCipher.dec.XORKeyStream(cipherText, data)
	return
}
