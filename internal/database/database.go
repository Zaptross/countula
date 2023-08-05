package database

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Ssl      string
	Timezone string
}

func GetConfig() DatabaseConfig {
	var config DatabaseConfig
	envconfig.Process("database", &config)
	return config
}

func Connect() *gorm.DB {
	config := GetConfig()

	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Database,
		config.Ssl,
		config.Timezone,
	)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Turn{}, &AuditLog{}, &StatisticRow{})

	return db
}
