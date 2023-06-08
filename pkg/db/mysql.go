package db

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"social/config"
)

var _ IDataSource = new(defaultMysqlDataSource)

// IDataSource 定义数据库数据源接口，按照业务需求可以返回主库链接Master和从库链接Slave
type IDataSource interface {
	Master() *gorm.DB
	Slave() *gorm.DB
	Close()
}

// defaultMysqlDataSource 默认mysql数据源实现
type defaultMysqlDataSource struct {
	master *gorm.DB
	slave  *gorm.DB
}

func (d *defaultMysqlDataSource) Master() *gorm.DB {
	if d.master == nil {
		panic("The [master] connection is nil, Please initialize it first.")
	}
	return d.master
}

func (d *defaultMysqlDataSource) Slave() *gorm.DB {
	if d.master == nil {
		panic("The [slave] connection is nil, Please initialize it first.")
	}
	return d.slave
}

func (d *defaultMysqlDataSource) Close() {
	// 关闭主库链接
	if d.master != nil {
		m, err := d.master.DB()
		if err != nil {
			m.Close()
		}
	}
	// 关闭从库链接
	if d.slave != nil {
		s, err := d.slave.DB()
		if err != nil {
			s.Close()
		}
	}
}

func NewDefaultMysql(c config.MysqlConfig) *defaultMysqlDataSource {
	return &defaultMysqlDataSource{
		master: connect(
			c.Username,
			c.Password,
			c.Host,
			c.Dbname,
			c.Port,
			c.MaximumPoolSize,
			c.MaximumIdleSize),
	}
}

func connect(user, password, host, dbname string, port, maxPoolSize, maxIdle int) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		password,
		host,
		port,
		dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true, // 缓存每一条sql语句，提高执行速度
	})
	if err != nil {
		panic(err)
	}
	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}

	alive(sqlDb)

	sqlDb.SetConnMaxLifetime(time.Hour)
	// 设置连接池大小
	sqlDb.SetMaxOpenConns(maxPoolSize)
	sqlDb.SetMaxIdleConns(maxIdle)
	return db
}

func alive(db *sql.DB) {
	log.Println("connecting to database... ")
	for {
		_, err := db.Exec("SELECT true")
		if err == nil {
			log.Println("database connected")
			return
		}

		// exponential backoff
		base, capacity := time.Second, time.Minute
		for backoff := base; err != nil; backoff <<= 1 {
			if backoff > capacity {
				backoff = capacity
			}

			jitter := rand.Int63n(int64(backoff * 3))
			sleep := base + time.Duration(jitter)
			time.Sleep(sleep)
			_, err := db.Exec("SELECT true")
			if err == nil {
				log.Println("database connected")
				return
			}
		}
	}
}
