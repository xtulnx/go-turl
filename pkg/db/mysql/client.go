// Package mysql provides MySQL connections
package mysql

import (
	"gorm.io/driver/sqlite"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/beihai0xff/turl/configs"
)

// New create a new gorm db
func New(c *configs.MySQLConfig) (*gorm.DB, error) {
	// l := logger.New(
	// 	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	// 	logger.Config{
	// 		SlowThreshold:             time.Second, // Slow SQL threshold
	// 		LogLevel:                  logger.Warn, // Log level
	// 		IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
	// 		ParameterizedQueries:      false,       // Log the parameter values
	// 		Colorful:                  true,        // Disable color
	// 	},
	// )

	dialector := c.Dialector
	if dialector == "" {
		dialector = "mysql"
	}
	var db *gorm.DB
	var err error
	switch strings.ToLower(dialector) {
	case "mysql":
		db, err = gorm.Open(mysql.Open(c.DSN), &gorm.Config{
			Logger:                 logger.Default,
			SkipDefaultTransaction: true,
			TranslateError:         true,
		})
	case "sqlite", "sqlite3":
		db, err = gorm.Open(sqlite.Open(c.DSN), &gorm.Config{
			Logger:                 logger.Default,
			SkipDefaultTransaction: true,
			TranslateError:         true,
		})
	default:
		return nil, gorm.ErrInvalidDB
	}

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(c.MaxConn)
	sqlDB.SetMaxOpenConns(c.MaxConn)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(time.Minute)

	return db, nil
}
