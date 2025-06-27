package gorm

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Pagination struct for pagination query
type Pagination struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Sort  string `json:"sort"`
}

// Filter struct for filtering query
type Filter struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    any    `json:"value"`
}

// Paginate returns a function that can be used as a GORM scope
func Paginate(pagination *Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pagination.Page == 0 {
			pagination.Page = 1
		}

		if pagination.Limit == 0 {
			pagination.Limit = 10
		}
		offset := (pagination.Page - 1) * pagination.Limit
		return db.Offset(offset).Limit(pagination.Limit).Order(pagination.Sort)
	}
}

// FilterScope returns a function that can be used as a GORM scope for filtering
func FilterScope(filters []Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, filter := range filters {
			if filter.Field != "" && filter.Operator != "" && filter.Value != nil {
				field := filter.Field
				operator := SanitizeOperator(filter.Operator)
				rawValue, ok := filter.Value.(string)
				if !ok {
					continue // skip non-string values
				}

				// Detect cast type from string value
				castType := detectCastType(rawValue)

				if strings.Contains(field, ".") {
					parts := strings.SplitN(field, ".", 2)
					jsonbCol := Sanitize(parts[0])
					jsonAttr := Sanitize(parts[1])
					jsonbExpr := fmt.Sprintf("%s->>'%s'", jsonbCol, jsonAttr)

					if castType != "" {
						jsonbExpr = fmt.Sprintf("CAST(%s AS %s)", jsonbExpr, castType)
					}

					if operator == "IN" {
						db = db.Where(fmt.Sprintf("%s %s (?)", jsonbExpr, operator), strings.Split(rawValue, "||"))
					} else {
						db = db.Where(fmt.Sprintf("%s %s ?", jsonbExpr, operator), rawValue)
					}
				} else {
					field = Sanitize(field)
					if operator == "IN" {
						db = db.Where(fmt.Sprintf("%s %s (?)", field, operator), strings.Split(rawValue, "||"))
					} else {
						db = db.Where(fmt.Sprintf("%s %s ?", field, operator), rawValue)
					}
				}
			}
		}
		return db
	}
}

// detectCastType attempts to guess PostgreSQL type from string
func detectCastType(value string) string {
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return "numeric"
	}
	if _, err := strconv.ParseBool(value); err == nil {
		return "boolean"
	}
	if _, err := time.Parse("2006-01-02", value); err == nil {
		return "date"
	}
	if _, err := time.Parse(time.RFC3339, value); err == nil {
		return "timestamp"
	}
	return "" // treat as plain string
}


// Sanitize removes any character that is not a letter, a number or an underscore.
func Sanitize(s string) string {
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return -1
	}, s)
}

// SanitizeOperator removes any character that is not a letter, a number or an underscore.
func SanitizeOperator(s string) string {
	// Allow only a subset of operators to prevent SQL injection.
	// You can extend this list if you need more operators.
	allowedOperators := map[string]bool{
		"=":       true,
		"<>":      true,
		">":       true,
		">=":      true,
		"<":       true,
		"<=":      true,
		"LIKE":    true,
		"IN":      true,
		"BETWEEN": true,
	}
	if allowedOperators[strings.ToUpper(s)] {
		return strings.ToUpper(s)
	}
	return "=" // default to "=" if the operator is not allowed
}
