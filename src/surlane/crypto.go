package surlane

import (
	"crypto/cipher"
	"crypto/aes"
	"crypto/md5"
	"crypto/rand"
	"net"
	"github.com/pkg/errors"
)

const (
	Ces128Cfb    = iota
	Ces256Cfb
)

var (
	cipherMap = map[int]StreamBuilder{
		Ces128Cfb: &CES128StreamBuilder{},
		Ces256Cfb: nil,
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
		Ces128Cfb,
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

func GetIV(config Config) []byte {
	streamBuilder := cipherMap[config.Method]
	return streamBuilder.genIV()
}

func GetIvLen(config Config) int {
	streamBuilder := cipherMap[config.Method]
	return streamBuilder.getIvLen()
}

func NewSurCipher4Client(config Config, iv []byte) (*SurCipher, error) {
	streamBuilder := cipherMap[config.Method]
	return streamBuilder.createCipher(config.Password, iv)
}

func NewSurCipher4Server(config Config, iv []byte) (*SurCipher, error) {
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

type SecureConn struct {
	*SurCipher
	net.Conn
}

func (secureConn *SecureConn) Read(buffer []byte) (n int, err error) {
	n, err = secureConn.Conn.Read(buffer)
	if n > 0 {
		copy(buffer[:n], secureConn.SurCipher.decrypt(buffer[:n]))
	}
	return
}

func (secureConn *SecureConn) Write(buffer []byte) (n int, err error) {
	return secureConn.Conn.Write(secureConn.SurCipher.encrypt(buffer))
}

func NewClientSecureConn(conn net.Conn, config Config, iv []byte) (*SecureConn, error) {
	surCipher, err := NewSurCipher4Client(config, iv)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &SecureConn{
		surCipher,
		conn,
	}, nil
}

func NewServerSecureConn(conn net.Conn, config Config, iv []byte) (*SecureConn, error) {
	surCipher, err := NewSurCipher4Server(config, iv)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &SecureConn{
		surCipher,
		conn,
	}, nil
}
