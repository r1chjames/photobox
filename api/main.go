package main

import (
	"fmt"
	"gitlab.com/r1chjames/photobox/api/internal/apiServer"
	"gitlab.com/r1chjames/photobox/api/internal/database"
	"gitlab.com/r1chjames/photobox/api/internal/types"
	"gitlab.com/r1chjames/photobox/api/internal/utils"
	"strconv"
)

func main() {
	appConfig := parseAppVariables()
	database.InitDbConnection(appConfig)
	database.PerformDbSetup(appConfig)
	stopAllRunningJobs()
	startAPIServer(appConfig)
}

func parseAppVariables() types.AppConfig {
	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbUser := utils.GetEnv("DB_USER", "photobox")
	dbPassword := utils.GetEnv("DB_PASSWORD", "photobox")
	dbName := utils.GetEnv("DB_NAME", "photobox")
	dbURL := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPassword, dbHost, dbName)
	resetSettings, _ := strconv.ParseBool(utils.GetEnv("RESET_SETTINGS", "false"))
	debugMode, _ := strconv.ParseBool(utils.GetEnv("DEBUG_MODE", "false"))

	return types.AppConfig{
		PhotoDir:      utils.GetEnv("PHOTO_DIR", "/photos"),
		ApiBasePath:   utils.GetEnv("API_BASE_PATH", "/api"),
		DbUrl:         dbURL,
		ResetSettings: resetSettings,
		DebugMode:     debugMode,
	}
}

func stopAllRunningJobs() {
	database.StopAllRunningJobs("NOT_RUNNING")
}

func startAPIServer(appConfig types.AppConfig) {
	apiServer.Start(appConfig)
}
