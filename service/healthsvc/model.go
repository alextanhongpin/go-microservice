package healthsvc

import "time"

// Health model provides useful information on the app runtime.
type Health struct {
	BuildDate time.Time `json:"build_date,omitempty"`
	GitTag    string    `json:"git_tag,omitempty"`
	Uptime    string    `json:"uptime"`
}
