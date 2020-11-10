package release

import "time"

type Notes struct {
	ReleaseDate  time.Time
	ChangedFiles []string
}

func New() *Notes {
	return &Notes{
		ReleaseDate: time.Now(),
	}
}

func (*Notes) Format(releaseDate time.Time) string {
	return releaseDate.Format("January 2, 2006 15:04:05 MST")
}
