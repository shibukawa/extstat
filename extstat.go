package extstat

import (
	"os"
	"time"
)

type ExtraStat struct {
	AccessTime  time.Time
	ModTime     time.Time
	CreatedTime time.Time
}

func NewFromFileName(name string) (*ExtraStat, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	return New(fi), nil
}
