package spotify

import (
	"fmt"
	"strings"
	"time"
)

type PartialDate struct {
	time.Time
}

func (p *PartialDate) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)

	layouts := []string{"2006-01-02", "2006-01"}
	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.Parse(layout, str)
		if err == nil {
			p.Time = t
			return nil
		}
	}
	return fmt.Errorf("cannot parse date: %s", str)
}
