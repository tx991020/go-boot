package {{package}}

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/tx991020/utils/logger"
)

var db *gorm.DB


func DatabaseInit() {

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=%s dbname=%s",
		viper.GetString("db.host"),
		viper.GetInt64("db.port"),
		viper.GetString("db.user"),
		viper.GetString("db.password"),
		viper.GetString("db.sslmode"),
		viper.GetString("db.dbname"),
	)
	logger.Infoa(fmt.Sprintf("mysql connect: %s max_open=%d max_idle=%d", connStr, viper.GetInt("db.max_open"),viper.GetInt("db.max_idle")))
	var err error
	db, err = gorm.Open("postgres", connStr)
	if err != nil {
		log.Errorf("connect to postgres fails: %s", err.Error())
		panic("database error")
	}

	if viper.GetBool("db.debug") || config.GoEnv() == "debug" {
		db.LogMode(true)
	}

	db.SingularTable(true)
	db.DB().SetMaxOpenConns(viper.GetInt("db.max_open"))
	db.DB().SetMaxIdleConns(viper.GetInt("db.max_idle"))
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
}

func DB(t string) *gorm.DB {
	return db
}

func TransactionCommitErr(tx *gorm.DB) {
	if err := tx.Commit().Error; err != nil {
		log.Errorf("Transaction Err: %s", err.Error())
	}
}

func TransactionRollbackErr(tx *gorm.DB) {
	if err := tx.Rollback().Error; err != nil {
		log.Errorf("Transaction Err: %s", err.Error())
	}
}

func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); ok {
		now := time.Now()
		// FIXME: 统一处理更新时间
		scope.SetColumn("modify_time", &now)
		scope.SetColumn("utime", &now)
	}
}
