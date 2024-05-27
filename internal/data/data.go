package data

import (
	"github.com/acmestack/gorm-plus/gplus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"kratos-demo/internal/conf"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// 第 4 步，用 wire 注入代码，修改 原来的 NewSet
// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewGormDB, NewstudentRepo, NewGinRepo)

// Data .
// 第 1 步引入 *gorm.DB
type Data struct {
	// TODO wrapped database client
	gormDB *gorm.DB
	log    *log.Helper
}

// 第 2 步初始化 gorm
func NewGormDB(c *conf.Data) (*gorm.DB, error) {
	dsn := c.Database.Source
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(50)
	sqlDB.SetMaxOpenConns(150)
	sqlDB.SetConnMaxLifetime(time.Second * 25)

	gplus.Init(db)
	return db, err
}

// NewData .
// 第 3 步，初始化 Data
func NewData(logger log.Logger, db *gorm.DB) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}

	return &Data{gormDB: db}, cleanup, nil
}

func Paginate(pageNo int32, pageSize int32) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageNo == 0 {
			pageNo = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (pageNo - 1) * pageSize
		return db.Offset(int(offset)).Limit(int(pageSize))
	}
}
