package di

import (
	"database/sql"

	"r2manager/handler"
	"r2manager/repository"
	service "r2manager/service/model"
)

func CreateSettingsHandler(db *sql.DB) *handler.SettingsHandler {
	repo := repository.NewSettingsRepository(db)
	svc := service.NewSettingsService(repo)
	return handler.NewSettingsHandler(svc)
}
