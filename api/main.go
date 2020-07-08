package main

import (
	"fmt"
	"gitlab.com/r1chjames/photobox/api/internal/apiServer"
	"gitlab.com/r1chjames/photobox/api/internal/database"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"gitlab.com/r1chjames/photobox/api/internal/utils"
)

func main() {
	appConfig := parseAppVariables()
	stopAllRunningJobs(appConfig)
	startApiServer(appConfig)
}

func parseAppVariables() AppConfig {
	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbUser := utils.GetEnv("DB_USER", "photobox")
	dbPassword := utils.GetEnv("DB_PASSWORD", "photobox")
	dbName := utils.GetEnv("DB_NAME", "photobox")
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPassword, dbHost, dbName)

	return AppConfig{
		PhotoDir: utils.GetEnv("PHOTO_DIR", "/photos"),
		ApiBasePath: utils.GetEnv("API_BASE_PATH", "/api"),
		DbUrl: dbUrl,
	}
}

func stopAllRunningJobs(appConfig AppConfig) {
	database.UpdateAllJobStatus(appConfig, false)
}

func startApiServer(appConfig AppConfig) {
	apiServer.Start(appConfig)
}