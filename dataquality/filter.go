package dataquality

import (
	"strconv"

	"github.com/oudsteyn/go-predix-timeseries/filter"
)

// Creates a filter to limit returned values by data quality
func Qualities(values ...Quality) filter.Filter {
	ret := make(map[string]interface{})
	vals := make([]string, len(values))
	for i, v := range values {
		vals[i] = strconv.Itoa(int(v))
	}
	ret["values"] = vals
	return filter.New("qualities", ret)
}
