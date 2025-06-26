package gorm

import (
	"net/http"
	"strconv"
	"strings"
)

// FromRequest extracts pagination and filters from an HTTP request.
func FromRequest(r *http.Request) (*Pagination, []Filter) {
	q := r.URL.Query()
	// Parse pagination
	limit, _ := strconv.Atoi(q.Get("limit"))
	page, _ := strconv.Atoi(q.Get("page"))
	sort := q.Get("sort")

	if limit == 0 {
		limit = 10
	}
	if page == 0 {
		page = 1
	}

	pagination := &Pagination{
		Limit: limit,
		Page:  page,
		Sort:  sort,
	}

	// Parse filters
	var filters []Filter
	for key, values := range q {
		if strings.HasPrefix(key, "filter[") && strings.HasSuffix(key, "]") {
			// Extract field from "filter[field]"
			field := key[7 : len(key)-1]
			if len(values) > 0 {
				// The value is expected to be in the format "operator,value"
				parts := strings.SplitN(values[0], ",", 2)
				if len(parts) == 2 {
					operator := parts[0]
					value := parts[1]
					var filterValue any
					if strings.EqualFold(operator, "in") {
						// Split value by comma for IN operator
						filterValue = strings.Split(value, ",")
					} else {
						filterValue = value
					}
					filters = append(filters, Filter{
						Field:    field,
						Operator: operator,
						Value:    filterValue,
					})
				}
			}
		}
	}

	return pagination, filters
}
