package initialize

import (
	"Project/config"
	"Project/global"
	"Project/initialize/internal"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//func GormMysql() {
//	dsn := "root:redmi@qwe123@tcp(100.79.174.8:30123)/students?charset=utf8mb4&parseTime=True&loc=Local"
//	var err error
//	if global.BG_DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
//		fmt.Println("数据库连接失败：", err)
//	} else {
//		fmt.Println("数据库连接成功")
//	}
//	//var tables []string
//	//if tables, err = global.GVA_DB.Migrator().GetTables(); err != nil {
//	//	fmt.Println("查询所有数据库表集合失败：", err)
//	//} else {
//	//	fmt.Println("查询所有数据库表集合成功", tables)
//	//}
//	// 自动迁移模型
//	global.BG_DB.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})
//}

// GormMysql 初始化Mysql数据库
func GormMysql() *gorm.DB {
	m := global.BG_CONFIG.Mysql
	if m.Dbname == "" {
		return nil
	}
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), internal.Gorm.Config(m.Prefix, m.Singular)); err != nil {
		return nil
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE="+m.Engine)
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}

// GormMysqlByConfig 初始化Mysql数据库用过传入配置
func GormMysqlByConfig(m config.Mysql) *gorm.DB {
	if m.Dbname == "" {
		return nil
	}
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), internal.Gorm.Config(m.Prefix, m.Singular)); err != nil {
		panic(err)
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE=InnoDB")
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}
