package setting

import "time"

var (
	// AppVersion is the version of the current build.
	AppVersion string
	// AppAuthor is the author of the project.
	AppAuthor string
	// AppStartTS is the timestamp when the application started.
	AppStartTS time.Time
	// IsProd is true if the application is running in production mode.
	IsProd bool
)

func init() {
	if AppVersion == "" {
		AppVersion = "development"
	}
}
