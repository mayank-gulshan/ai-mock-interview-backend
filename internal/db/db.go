// Package db wires up the GORM connection.
//
// NOTE: spring.jpa.hibernate.ddl-auto=update in application.properties makes
// Hibernate create/alter tables automatically on boot to match the @Entity
// classes. GORM's AutoMigrate() is the closest idiomatic equivalent and is
// called once at startup (see cmd/server/main.go), covering User, Session,
// and Answer.
//
// spring.jpa.show-sql=true / format_sql=true map to gorm's logger.Info mode,
// which prints formatted SQL statements — enabled below to match.
package db

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func Connect(dsn string) (*gorm.DB, error) {
	newLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  gormlogger.Info, // NOTE: mirrors spring.jpa.show-sql=true
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}
	return gdb, nil
}
