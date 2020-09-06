package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	MysqlConnectionDSNFormat = "%s:%s@%s(%s)/%s?charset=utf8&parseTime=True"
)

type Db struct {
	instance     *gorm.DB
}

var DbWrapper *Db

// NewDb instantiates Db and connects to database.
func InitializeDB(config DBConfig) (*Db, error) {
	db := &Db{}
	if err := db.connect(config); err != nil {
		fmt.Println("Error while connecting to DataBase")
		return nil, err
	}
	return db, nil
}

// Instance returns underlying gorm db instance.
func (db *Db) Instance() *gorm.DB {
	return db.instance
}

func (db *Db) connect(config DBConfig) error {
	var err error
	if db.instance, err = gorm.Open(config.Dialect, getConnectionPath(config)); err != nil {
		fmt.Println("Error while connecting to DB :: ", err)
		return err
	}
	return nil
}

func getConnectionPath(config DBConfig) string {
	return fmt.Sprintf(MysqlConnectionDSNFormat, config.Username, config.Password, config.Protocol, config.URL, config.Database)
}

func DatabaseHealthCheck() (bool, error) {
	sampleQuery := "SELECT COUNT(*) FROM download_request"
	dbErr := RepoClient.DBFromContext(nil).Exec(sampleQuery).Error
	return true, dbErr
}