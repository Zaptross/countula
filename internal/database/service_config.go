package database

import "gorm.io/gorm"

type ServiceConfig struct {
	gorm.Model
	MaintenanceMode bool
}

func SetMaintenanceMode(db *gorm.DB, enable bool) {
	var serviceCfg ServiceConfig
	db.First(&serviceCfg)

	serviceCfg.MaintenanceMode = enable
	db.Save(&serviceCfg)
}

func GetServiceConfig(db *gorm.DB) ServiceConfig {
	var serviceCfg ServiceConfig
	db.FirstOrCreate(&serviceCfg)
	return serviceCfg
}
