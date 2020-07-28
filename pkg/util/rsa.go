package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	gerrors "golang_im/pkg/errs"
	"golang_im/pkg/log"
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

type TokenInfo struct {
	AppId    int64 `json:"app_id"`    // appId
	UserId   int64 `json:"user_id"`   // 用户id
	DeviceId int64 `json:"device_id"` // 设备id
	Expire   int64 `json:"expire"`    // 过期时间
}

// 获取token
func GetToken(appId, userId, deviceId int64, expire int64, publicKey string) (string, error) {
	info := TokenInfo{
		AppId:    appId,
		UserId:   userId,
		DeviceId: deviceId,
		Expire:   expire,
	}
	bytes, err := json.Marshal(info)
	if err != nil {
		log.Error(err)
		return "", err
	}

	token, err := RsaEncrypt(bytes, []byte(publicKey))
	if err != nil {
		log.Error(err)
		return "", err
	}
	return base64.StdEncoding.EncodeToString(token), nil
}

// 对加密的token进行解码
func DecryptToken(token string, privateKey string) (*TokenInfo, error) {
	bytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	result, err := RsaDecrypt(bytes, Str2bytes(privateKey))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var info TokenInfo
	err = jsoniter.Unmarshal(result, &info)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &info, nil
}

// 加密
func RsaEncrypt(origData []byte, publicKey []byte) ([]byte, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 解密
func RsaDecrypt(ciphertext []byte, privateKey []byte) ([]byte, error) {
	//解密
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// VerifySecretKey 对用户秘钥进行校验
func VerifyToken(appId, userId, deviceId int64, token string) error {
	info, err := DecryptToken(token, PrivateKey)
	if err != nil {
		return gerrors.ErrUnauthorized
	}

	if !(info.AppId == appId && info.UserId == userId && info.DeviceId == deviceId) {
		return gerrors.ErrUnauthorized
	}

	if info.Expire < time.Now().Unix() {
		return gerrors.ErrUnauthorized
	}
	return nil
}

// PrivateKey 公钥
var PrivateKey = `
-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDBpOO3+WA/rzLEqXcLX0IBJhg0w3Pf5WhWDLAvgVKA1P7vGGXP
qmb5oRe0cbH7Q5Ef89JRS9/+v4cZEALN2oAKQOUFsRoH5yvFSHFUT6tKNG43ooRE
OnIAXZ4Tb8cNHbF5XSMNtf2+Ne0NfXOEZaqNPlvlnWSjujRUPLTCAdgP6wIDAQAB
AoGBAL62h7Psffeas/RuNplTou0AuLxWduvew2hkLK1Mv5W0sLOIItVorOxT1MXZ
aAHf5LFEcDGy+ZOqzAJJ+4kEFi7MIuWsZ9lVtuytRNIiBIWtdZzZs40kQ43UU1xy
izG7FGI35gU77L9esbGu69KdJU2vUqyZVvHHqNDQSFtC0rwpAkEA50ibRz0c7PEF
cVFdZGsXXCjWW94+PjY4F+Rh2IsE2drrk8UxlWtouXNuwq2GgICUFAaYpQaSJbwD
QLCWkurgnQJBANZWjScJtaEmqK1DasHVNf13ElO3NKgWWUUlQjYrpY94N0xHZQeW
o9xuAWQV6LpA8kP/Shc06B5dhs5D5JFruCcCQQCSz+sJaIiw+znqOazf7n7QmHeh
r0yxbvdiay2VKIH2zFmX3qff4mOCvPyFBWOItJXKtHk24BnrbBJggPfD4OadAkA8
difVJkkFD3mvfoAD85gKSudxlBGXhM5j0fHOhBts0DWRH+ag8F6C1MkxqXh/6cgt
ZDtLNpJv1mQrlT1JxEArAkEAg2ZXyT+bCFPm7KtVVrX9ilqNBUv1Q8fQLAgEvEpy
l1H8NQdcgwtmPmRLsp3rKpmPIV1VuSacdkdOvrlb7nYt+w==
-----END RSA PRIVATE KEY-----
`

// 公钥: 根据私钥生成
var PublicKey = `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDBpOO3+WA/rzLEqXcLX0IBJhg0
w3Pf5WhWDLAvgVKA1P7vGGXPqmb5oRe0cbH7Q5Ef89JRS9/+v4cZEALN2oAKQOUF
sRoH5yvFSHFUT6tKNG43ooREOnIAXZ4Tb8cNHbF5XSMNtf2+Ne0NfXOEZaqNPlvl
nWSjujRUPLTCAdgP6wIDAQAB
-----END PUBLIC KEY-----
`
