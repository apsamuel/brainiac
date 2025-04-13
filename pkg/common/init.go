package common

import "os"

/* initialize is executed on module loading */
func init() {
	AppRoot := getAppRoot()
	if AppRoot != "" {
		os.Setenv("MODULE_PATH", AppRoot)
	}

	Logger := GetLogger()
	Logger.Info().Msg("module path set to: " + AppRoot)
}
