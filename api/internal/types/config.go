package types

import "time"

type AppConfig struct {
	PhotoDir      string
	ApiBasePath   string
	DbUrl         string
	ResetSettings bool
	DebugMode     bool
	Timezone      *time.Location
}
