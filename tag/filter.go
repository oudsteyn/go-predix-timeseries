package tag

import "github.com/oudsteyn/go-predix-timeseries/filter"

type attr struct {
	name   string
	values []string
}

func Attr(name string, values ...string) attr {
	return attr{name, values}
}

// Filters query results by attributes and their values
func Attributes(attributes ...attr) filter.Filter {
	ret := make(map[string]interface{})
	for _, attr := range attributes {
		ret[attr.name] = attr.values
	}
	return filter.New("attributes", ret)
}
