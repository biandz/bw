package mysql

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitMysql(dsn string) *gorm.DB {
	//dsn := fmt.Sprintf(
	//	"%s:%s@tcp(%s:%s)/%s?timeout=90s&parseTime=true&loc=Local&collation"+
	//		"=utf8mb4_general_ci", c.DbUserName, c.DbPassword, c.DbHost, c.DbPort, c.DbDatabase,
	//)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			Writer{},
			logger.Config{
				SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
				LogLevel:                  logger.Info,            // Log level
				IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,                   // Disable color
			},
		),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		PrepareStmt: true,
	})

	if err != nil {
		log.Fatal("数据库初始化失败:", err.Error())
	}

	dbSql, err := db.DB()
	if err != nil {
		log.Fatal("数据库连接池初始化失败:", err.Error())
	}
	dbSql.SetMaxIdleConns(10)
	dbSql.SetMaxOpenConns(100)
	dbSql.SetConnMaxLifetime(time.Hour)
	log.Println("数据库初始化成功")

	return db
}

type Writer struct {
}

func (w Writer) Printf(format string, args ...interface{}) {
	log.Println(args[len(args)-1])
	now := time.Now()
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs/orm"
	}
	if err := os.MkdirAll(logFilePath, 0777); err != nil {
		log.Println(err.Error())
	}
	logFileName := "mysql_" + now.Format("2006-01-02") + ".log"
	//日志文件
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			log.Println(err.Error())
		}
	}
	//写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Println("err", err)
	}
	// 实例化
	logger := logrus.New()
	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)
	//设置输出
	logger.Out = src

	// 日志格式
	logger.WithFields(logrus.Fields{
		"sql": args[len(args)-1],
	}).Info()
}
