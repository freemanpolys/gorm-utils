package gorm

import (
	"net/http"
	"strconv"
	"strings"
)

// parseFilterConditions parses the filter string into conditions, handling 'in' operator values with commas.
func parseFilterConditions(filterStr string) []string {
	if filterStr == "" {
		return nil
	}
	tokens := strings.Split(filterStr, ",")
	var conditions []string
	i := 0
	for i < len(tokens) {
		// Try to find a token with at least two colons (field:operator:value)
		parts := strings.SplitN(tokens[i], ":", 3)
		if len(parts) == 3 {
			field := parts[0]
			operator := parts[1]
			value := parts[2]
			if strings.EqualFold(operator, "in") {
				// Collect all tokens until the next token with at least two colons or end
				j := i + 1
				for j < len(tokens) && !strings.Contains(tokens[j], ":") {
					value += "," + tokens[j]
					j++
				}
				conditions = append(conditions, field+":"+operator+":"+value)
				i = j
				continue
			}
			conditions = append(conditions, tokens[i])
		}
		i++
	}
	return conditions
}

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

	// Parse filters from new format: field::operator::value~~field2::operator2::value2
	var filters []Filter
	filterStr := q.Get("filter")
	if filterStr != "" {
		conditions := strings.Split(filterStr, "~~")
		for _, cond := range conditions {
			parts := strings.SplitN(cond, "::", 3)
			if len(parts) == 3 {
				field := parts[0]
				operator := parts[1]
				value := parts[2]
				var filterValue any
				if strings.EqualFold(operator, "in") || strings.EqualFold(operator, "between") {
					filterValue = strings.Split(value, "||")
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

	return pagination, filters
}
