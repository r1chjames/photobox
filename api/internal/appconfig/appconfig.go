package appconfig

import (
	"fmt"
	"gitlab.com/r1chjames/photobox/api/internal/core/utils"
	"strconv"
	"time"
)

type AppConfig struct {
	PhotoDir      string
	ApiBasePath   string
	DbUrl         string
	ResetSettings bool
	DebugMode     bool
	Timezone      *time.Location
	Token         string
	TokenDuration time.Duration
}

func New() *AppConfig {
	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbPort := utils.GetEnv("DB_PORT", "5432")
	dbUser := utils.GetEnv("DB_USER", "photobox")
	dbPassword := utils.GetEnv("DB_PASSWORD", "photobox")
	dbName := utils.GetEnv("DB_NAME", "photobox")

	dbURL := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", dbHost, dbUser, dbPassword, dbName, dbPort)
	resetSettings, _ := strconv.ParseBool(utils.GetEnv("RESET_SETTINGS", "false"))
	debugMode, _ := strconv.ParseBool(utils.GetEnv("DEBUG_MODE", "false"))
	timezone, _ := time.LoadLocation(utils.GetEnv("TIMEZONE", "Europe/London"))

	token := utils.GetEnv("TOKEN", "")
	tokenDuration, _ := time.ParseDuration(utils.GetEnv("TOKEN_DURATION", "1h"))

	return &AppConfig{
		PhotoDir:      utils.GetEnv("PHOTO_DIR", "/photos"),
		ApiBasePath:   utils.GetEnv("API_BASE_PATH", "/api"),
		DbUrl:         dbURL,
		ResetSettings: resetSettings,
		DebugMode:     debugMode,
		Timezone:      timezone,
		Token:         token,
		TokenDuration: tokenDuration,
	}
}
