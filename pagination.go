package gorm

import (
	"fmt"
	"strings"

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
				// Basic protection against SQL injection
				// You might want to use a more robust solution for production
				// For example, a whitelist of allowed fields and operators
				field := Sanitize(filter.Field)
				operator := SanitizeOperator(filter.Operator)

				if operator == "IN" {
					db = db.Where(fmt.Sprintf("%s %s (?)", field, operator), filter.Value)
				} else {
					db = db.Where(fmt.Sprintf("%s %s ?", field, operator), filter.Value)
				}
			}
		}
		return db
	}
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
		"=":    true,
		"<>":   true,
		">":    true,
		">=":   true,
		"<":    true,
		"<=":   true,
		"LIKE": true,
		"IN":   true,
	}
	if allowedOperators[strings.ToUpper(s)] {
		return strings.ToUpper(s)
	}
	return "=" // default to "=" if the operator is not allowed
}

