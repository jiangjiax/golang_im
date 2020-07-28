## 配置数据库和日志

#### 配置数据库

在 pkg/db 目录下新建 mysql.go 文件，用于连接 mysql 数据库：（sql文件已放在sql目录下了）

``` Go
package db

import (
	"database/sql"
	"fmt"
	"golang_im/config"

	_ "github.com/go-sql-driver/mysql"
)

var DBCli *sql.DB

func Mysql_init() {
	var err error
	fmt.Println("config.LogicConf.MySQL", config.DBConf.MySQL)
	DBCli, err = sql.Open("mysql", config.DBConf.MySQL)
	if err != nil {
		panic(err)
	}
}
```

在同一个目录下新建 redis.go 文件，用于连接 redis 数据库：

``` Go
package db

import (
	"golang_im/config"

	"github.com/go-redis/redis"
)

var RedisCli *redis.Client

func Redis_init() {
	addr := config.DBConf.RedisIP
	password := config.DBConf.RedisPwd
	RedisCli = redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       0,
		Password: password,
	})

	_, err := RedisCli.Ping().Result()
	if err != nil {
		panic(err)
	}
}
```

#### 配置日志zap

zap 是 uber 开源的高性能日志库。我们将在 pkg/log 下新建 log.go 用于控制日志：

``` Go
package log

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// errorLogger
var errorLogger *zap.SugaredLogger

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func getLoggerLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.InfoLevel
}

func init() {
	level := getLoggerLevel("debug")
	now := time.Now()
	syncWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fmt.Sprintf("../../../logs/%04d年%02d月%02d日%02d时.log", now.Year(), now.Month(), now.Day(), now.Hour()),
		MaxSize:    100,
		MaxAge:     10,
		MaxBackups: 30,
		Compress: true,
	})
	runMode := gin.Mode()
	var encoder zapcore.EncoderConfig
	if runMode == "debug" {
		encoder = zap.NewDevelopmentEncoderConfig()
	} else {
		encoder = zap.NewProductionEncoderConfig()
		encoder.EncodeTime = zapcore.EpochTimeEncoder
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoder), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), syncWriter), zap.NewAtomicLevelAt(level))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	errorLogger = logger.Sugar()
}

// Debug Debug
func Debug(args ...interface{}) {
	errorLogger.Debug(args...)
}

// Debugf Debugf
func Debugf(template string, args ...interface{}) {
	errorLogger.Debugf(template, args...)
}

// Info Info
func Info(args ...interface{}) {
	errorLogger.Info(args...)
}

// Infof Infof
func Infof(template string, args ...interface{}) {
	errorLogger.Infof(template, args...)
}

// Warn Warn
func Warn(args ...interface{}) {
	errorLogger.Warn(args...)
}

// Warnf Warnf
func Warnf(template string, args ...interface{}) {
	errorLogger.Warnf(template, args...)
}

// Error Error
func Error(args ...interface{}) {
	errorLogger.Error(args...)
}

// Errorf Errorf
func Errorf(template string, args ...interface{}) {
	errorLogger.Errorf(template, args...)
}

// DPanic DPanic
func DPanic(args ...interface{}) {
	errorLogger.DPanic(args...)
}

// DPanicf DPanicf
func DPanicf(template string, args ...interface{}) {
	errorLogger.DPanicf(template, args...)
}

// Panic Panic
func Panic(args ...interface{}) {
	errorLogger.Panic(args...)
}

// Panicf Panicf
func Panicf(template string, args ...interface{}) {
	errorLogger.Panicf(template, args...)
}

// Fatal Fatal
func Fatal(args ...interface{}) {
	errorLogger.Fatal(args...)
}

// Fatalf Fatalf
func Fatalf(template string, args ...interface{}) {
	errorLogger.Fatalf(template, args...)
}
```

---

## 错误回调

在 pkg/errs 目录下新建 gerrors.go 文件，编写处理错误的工具函数：

``` Go
package gerrors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUnauthorized = newError(1, "未登录")
	ErrUnDeviceid   = newError(2, "无设备")
)

func newError(code int, message string) error {
	return status.New(codes.Code(code), message).Err()
}
```

---

## rsa权限验证

#### 获取rsa密钥

rsa 是目前使用最广泛的公钥密码体制之一，我们使用 rsa 对用户权限做验证。首先需要生成 rsa 公钥和私钥文件，可以使用以下 Go 代码生成：

``` Go
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func main() {
	//rsa 密钥文件产生
	GenRsaKey(1024)
}

//RSA公钥私钥产生
func GenRsaKey(bits int) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	file, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	file, err = os.Create("public.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}
```

执行 go run main.go 后，得到 private.pem 和 public.pem 文件，文件内分别是私钥和公钥密码。

#### 编写rsa加解密工具函数

在 pkg/util 目录下新建 rsa.go 文件：

``` Go
package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"golang_im/pkg/log"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

type TokenInfo struct {
	AppId    int64 `json:"app_id"`    // appId
	UserId   int64 `json:"user_id"`   // 用户id
	DeviceId int64 `json:"device_id"` // 设备id
	Expire   int64 `json:"expire"`    // 过期时间
}

// GetToken 获取token
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

// DecryptToken 对加密的token进行解码
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

// PrivateKey 公钥
var PrivateKey = ``

// 公钥: 根据私钥生成
var PublicKey = ``
```