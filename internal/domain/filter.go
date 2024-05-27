package domain

type FilterCriteria int

const (
	Equal FilterCriteria = iota
	NotEqual
	GreaterThanOrEqual
	LessThanOrEqual
	GreaterThan
	LessThan
	In
)

var FilterCriteriaString = map[FilterCriteria]string{
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
