package model

import "time"

type Health struct {
	BuildDate time.Time `json:"build_date,omitempty"`
	GitTag    string    `json:"git_tag,omitempty"`
	Uptime    string    `json:"uptime"`
}
