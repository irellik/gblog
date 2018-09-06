package service

import (
	"flag"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	}
}

var globalConfig = Config{}

// 载入配置文件
func LoadConfig() {

	confFile := flag.String("c", "", "配置文件")
	flag.Parse()
	yamlFile, err := ioutil.ReadFile(*confFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &globalConfig)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
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
