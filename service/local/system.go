package local

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// 定义配置文件yaml结构
type Config struct {
	Mysql struct {
		Host     string `yaml:"host"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Port     int    `yaml:"port"`
		Database string `yaml:"database"`
	} `yaml:"mysql"`
	Site struct {
		PageSize int    `yaml:"pageSize"`
		Address  string `yaml:"address"`
		RememberDays int `yaml:"remember_days"`
		Domain string `yaml:"domain"`
	} `yaml:"site"`
	Session struct{
		MaxAge int `yaml:"maxAge"`
	}`yaml:"session"`
	AppKey string `yaml:"appKey"`
}

var globalConfig = Config{}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 载入配置文件
func LoadConfig() Config{
	currentPath := GetCurrentPath()
	confFile := flag.String("config", fmt.Sprintf("%s/config/app.yaml", currentPath), "配置文件")
	flag.Parse()
	yamlFile, err := ioutil.ReadFile(*confFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &globalConfig)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	// 检查AppKey是否合法
	if len(globalConfig.AppKey) != 32 {
		globalConfig.AppKey = RandStr(32)
		yamlText, err := yaml.Marshal(globalConfig)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
		// 重写配置文件
		ioutil.WriteFile(*confFile, yamlText, 0644)
	}
	return globalConfig
}

// 获取配置
func GetConfig() Config {
	return globalConfig
}

// 获取当前页面页码
func GetPage(c *gin.Context) int {
	// 获取页码
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}
	return page
}

func GetCurrentPath() string {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("%s", err)
	}
	return strings.Replace(path, "\\", "/", -1)
}

func RandStr(length int) string {
	letter := []byte("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	randStr := make([]byte, length)
	for i := range randStr {
		randStr[i] = letter[rand.Intn(len(letter))]
	}
	return string(randStr)
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

func Encrypt(data []byte) ([]byte, error) {
	return AesEncrypt(data, []byte(globalConfig.AppKey))
}

func Decrypt(data []byte) ([]byte, error) {
	return AesDecrypt(data, []byte(globalConfig.AppKey))
}
