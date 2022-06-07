package db

import (
	"fmt"
	"github.com/go-redis/redis"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"path"
	"runtime"
)

// 配置文件
type config struct {
	Mysql mysqlConfig `yaml:"mysql"`
	Redis redisConfig `yaml:"redis"`
}

// mysql 配置
type mysqlConfig struct {
	Addr   string `yaml:"addr"`
	User   string `yaml:"user"`
	Passwd string `yaml:"passwd"`
	DBName string `yaml:"db-name"`
	DSN    string
}

// redis 配置
type redisConfig struct {
	Addr   string `yaml:"addr"`
	Passwd string `yaml:"passwd"`
}

var DB *gorm.DB
var Redis *redis.Client

func Init() {
	DatabaseInit()
}

func DatabaseInit() {
	// 加载配置文件获取 DSN
	c := &config{}
	_, filename, _, _ := runtime.Caller(0)
	loadConfig(path.Dir(path.Dir(path.Dir(filename)))+"/config/config.yml", c)

	// 初始化两个数据库
	var err error
	DB, err = gorm.Open(mysql.Open(c.Mysql.DSN), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	})
	Redis = redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Passwd,
	})

	if err != nil {
		log.Fatal("failed to connect database")
	}
}

// loadConfig 读取 yaml 配置文件
func loadConfig(path string, c *config) {
	// 读取 yaml 文件
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("打开 yaml 文件失败", err)
		return
	}

	// 解析 yaml 文件得到 mysql 和 redis 配置
	err = yaml.NewDecoder(file).Decode(c)
	if err != nil {
		fmt.Println("解析 yaml 文件失败", err)
	}

	// 拼接 dsn
	c.Mysql.DSN = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Mysql.User, c.Mysql.Passwd, c.Mysql.Addr, c.Mysql.DBName)
}
