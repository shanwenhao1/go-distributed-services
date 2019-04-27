package db

import (
	"encoding/xml"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"go-distributed-services/domain/model"
	"go-distributed-services/infra/log"
	"io/ioutil"
)

var (
	ds  gorm.DB
	tds gorm.DB
)

type xmlStruct struct {
	DbName     string `xml:"dbname"`
	DbUser     string `xml:"dbuser"`
	DbUPwd     string `xml:"dbupwd"`
	DbUrl      string `xml:"dburl"`
	DbMaxConn  int    `xml:"dbmaxconn"`
	DbMaxIdle  int    `xml:"dbmaxidle"`
	DbLogModel bool   `xml:"dblogmodel"`
}

// 生产数据库
func GetDS() gorm.DB {
	return ds
}

// 测试数据库
func GetTDS() gorm.DB {
	return tds
}

//数据库连接初始化
func InitializedDataSource(path string) {
	// TODO 考虑使用consul获取配置, 可保留config文件方式
	fmt.Println("---------------init db connect")
	dataS, rErr := ioutil.ReadFile(path)
	if rErr != nil {
		log.LogWithTag(log.ERROR, log.InitSer, "读取数据库配置文件异常:[%v]", rErr)
		panic(rErr.Error())
	}
	configData := xmlStruct{}
	xErr := xml.Unmarshal(dataS, &configData)
	if xErr != nil {
		log.LogWithTag(log.ERROR, log.InitSer, "解析数据库配置文件异常:[%v]", xErr)
		panic(xErr.Error())
	}
	db, err := gorm.Open("mysql", fmt.Sprintf("%v:%v@%v/%v?charset=utf8&parseTime=True",
		configData.DbUser, configData.DbUPwd, configData.DbUrl, configData.DbName))
	if err != nil {
		log.LogWithTag(log.ERROR, log.InitSer, "初始化数据源异常:%v", err)
		panic(err.Error())
	} else {
		db.LogMode(configData.DbLogModel)
		db.SingularTable(true)
		db.DB().SetMaxOpenConns(configData.DbMaxConn)
		db.DB().SetMaxIdleConns(configData.DbMaxIdle)
		// 自动迁移只会创建表、缺少列和索引等, 并不会执行任何删除及修改操作以保护数据
		db.AutoMigrate(&model.User{})
		ds = *db
		log.LogWithTag(log.INFO, log.InitSer, "数据源已初始化完成[最大打开连接数:%v,最大空闲连接数:%v]",
			configData.DbMaxConn, configData.DbMaxIdle)
	}
}
