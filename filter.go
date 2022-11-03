package hubspot

type FilterQuery struct {
	FilterGroups []FilterGroup `json:"filterGroups,omitempty"`
	Sorts        []FilterSort  `json:"sorts,omitempty"`
	Query        string        `json:"query,omitempty"`
	Properties   []string      `json:"properties,omitempty"`
	Limit        int           `json:"limit,omitempty"`
	After        int           `json:"after,omitempty"`
}

type FilterSort struct {
	PropertyName string `json:"propertyName"`
	Direction    string `json:"direction"`
}

type FilterGroup struct {
	Filters []Filter `json:"filters,omitempty"`
}

type Filter struct {
	Value        string   `json:"value,omitempty"`
	Values       []string `json:"values,omitempty"`
	PropertyName string   `json:"propertyName,omitempty"`
	Operator     string   `json:"operator,omitempty"`
}
