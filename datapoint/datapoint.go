// Package datapoint contains data point representation.
package datapoint

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/oudsteyn/go-predix-timeseries/dataquality"
	"github.com/oudsteyn/go-predix-timeseries/measurement"
)

type Datapoint struct {
	Measure   measurement.Measurement
	Timestamp time.Time
	Quality   dataquality.Quality
}

// Custom JSON marshaling because Time Series expects a data point to be a 3-element array
func (p *Datapoint) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("[%d,%s,%d]", p.Timestamp.UnixNano()/int64(time.Millisecond), p.Measure.Value(), p.Quality)), nil
}

// Custom JSON unmarshalling because Time Series returns a datapoint as a 3-element array
func (p *Datapoint) UnmarshalJSON(bs []byte) error {
	s := string(bs)
	if !strings.HasPrefix(s, "[") && !strings.HasSuffix(s, "]") {
		return errors.New("Not a datapoint")
	}
	values := strings.Split(strings.Trim(s, "[]"), ",")
	if len(values) != 3 {
		return errors.New("Not a datapoint")
	}
	timestamp, err := strconv.ParseInt(values[0], 10, 64)
	if err != nil {
		return nil
	}
	p.Timestamp = time.Unix(0, timestamp*int64(time.Millisecond))
	p.Measure, err = measurement.FromString(values[1])
	if err != nil {
		return err
	}
	q, err := strconv.Atoi(values[2])
	if err != nil {
		return err
	}
	p.Quality = dataquality.Quality(q)
	return nil
}
