package main

import (
	"fmt"
	"gitlab.com/r1chjames/photobox/api/internal/apiServer"
	"gitlab.com/r1chjames/photobox/api/internal/components"
	"gitlab.com/r1chjames/photobox/api/internal/database"
	"gitlab.com/r1chjames/photobox/api/internal/types"
	"gitlab.com/r1chjames/photobox/api/internal/utils"
	"strconv"
	"time"
)

func main() {
	appConfig := parseAppVariables()

	dbEnv := database.InitDbConnection(appConfig)
	dbEnv.PerformDbSetup(appConfig)

	components.InitScheduler(appConfig)
	components.StopAllRunningJobs(dbEnv)
	components.AddScheduledJobs(appConfig, dbEnv)

	apiServer.NewServer(appConfig, dbEnv)
}

func parseAppVariables() types.AppConfig {
	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbPort := utils.GetEnv("DB_PORT", "5432")
	dbUser := utils.GetEnv("DB_USER", "photobox")
	dbPassword := utils.GetEnv("DB_PASSWORD", "photobox")
	dbName := utils.GetEnv("DB_NAME", "photobox")

	dbURL := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", dbHost, dbUser, dbPassword, dbName, dbPort)
	resetSettings, _ := strconv.ParseBool(utils.GetEnv("RESET_SETTINGS", "false"))
	debugMode, _ := strconv.ParseBool(utils.GetEnv("DEBUG_MODE", "false"))
	timezone, _ := time.LoadLocation(utils.GetEnv("TIMEZONE", "Europe/London"))

	return types.AppConfig{
		PhotoDir:      utils.GetEnv("PHOTO_DIR", "/photos"),
		ApiBasePath:   utils.GetEnv("API_BASE_PATH", "/api"),
		DbUrl:         dbURL,
		ResetSettings: resetSettings,
		DebugMode:     debugMode,
		Timezone:      timezone,
	}
}
