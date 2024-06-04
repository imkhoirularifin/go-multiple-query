package domain

type FilterCriteria int

const (
	NoMatch FilterCriteria = iota // noMatch is the default value if no filter is applied
	Equal
	NotEqual
	GreaterThanOrEqual
	LessThanOrEqual
	GreaterThan
	LessThan
	In
)

var FilterCriteriaString = map[FilterCriteria]string{
	NoMatch:            "noMatch",
	Equal:              "equal",
	NotEqual:           "notEqual",
	GreaterThanOrEqual: "greaterThanOrEqual",
	LessThanOrEqual:    "lessThanOrEqual",
	GreaterThan:        "greaterThan",
	LessThan:           "lessThan",
	In:                 "in",
}

type QueryRequest struct {
	Page      int
	Limit     int
	OrderBy   string
	SortOrder string
	Filters   []QueryFilter
}

type QueryFilter struct {
	Key    string
	Filter FilterCriteria
	Value  []string
}
