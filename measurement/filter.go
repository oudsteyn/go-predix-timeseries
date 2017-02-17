package measurement

import (
	"github.com/oudsteyn/go-predix-timeseries/filter"
)

func condition(condition string, values ...Measurement) filter.Filter {
	ret := make(map[string]interface{})
	ret["condition"] = condition
	vals := make([]string, len(values))
	for i, v := range values {
		vals[i] = v.Value()
	}
	ret["values"] = vals

	return filter.New("measurements", ret)
}

// Filters data points greater than ( > ) a specific data point value
func Gt(values ...Measurement) filter.Filter { return condition("gt", values...) }

// Filters data points greater than ( > ) or equal to ( = ) a specific data point value
func Ge(values ...Measurement) filter.Filter { return condition("ge", values...) }

// Filters data points less than ( < ) a specific data point value
func Lt(values ...Measurement) filter.Filter { return condition("lt", values...) }

// Filters data points less than ( < ) or equal to ( = ) a specific data point value
func Le(values ...Measurement) filter.Filter { return condition("le", values...) }

// Filters data points equal to ( = ) a specific data point value
func Eq(values ...Measurement) filter.Filter { return condition("eq", values...) }

// Filters data points not equal to ( = ) a specific data point value
func Ne(values ...Measurement) filter.Filter { return condition("ne", values...) }
