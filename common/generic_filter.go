/**
 * Generic DB filter for get all records
 * Should support various operators such as eq, ne, gt, lt, gte, lte, in, nin, like, nlike, between, nbetween
 * Should support various logical operators such as and, or, not
 * Should support various data types such as string, number, boolean, date, array, object
**/
package common

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-chassis/openlog"
	"gorm.io/gorm"
)

const SPACE = " "
const QMARK = "?"

var Operators = map[string]string{
	"eq":       "=",
	"ne":       "!=",
	"gt":       ">",
	"lt":       "<",
	"gte":      ">=",
	"lte":      "<=",
	"in":       "IN",
	"nin":      "NOT IN",
	"like":     "LIKE",
	"ilike":    "ILIKE",
	"nlike":    "NOT LIKE",
	"between":  "BETWEEN",     // for range
	"nbetween": "NOT BETWEEN", // for range
}

func validateOperator(operator string) bool {
	if _, ok := Operators[operator]; ok {
		return true
	}
	return false
}

// validate values does not contain sql injection
func validateValue(value string) bool {
	return strings.ContainsAny(value, ";") || strings.ContainsAny(value, "--") || strings.ContainsAny(value, "/*") || strings.ContainsAny(value, "*/")
}

type TableFilters struct {
	Queries []ColumnQuery `json:"queries"`
	Sort    []Sort        `json:"sort"`
}

type ColumnQuery struct {
	Col       string `json:"col"`
	Condition string `json:"cond"`
	Value     string `json:"value"`
}

type Sort struct {
	Col   string `json:"col"`
	Order string `json:"order"`
}

// GetWhereScope closure for returining a scope in gorm [for Reference :- https://gorm.io/docs/advanced_query.html#Scopes]
func GetWhereScope(query ColumnQuery) func(*gorm.DB) *gorm.DB {
	// validate Operator
	cond, ok := Operators[query.Condition]
	if !ok {
		return nil
	}
	openlog.Debug("Operator is " + cond)
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(query.Col+SPACE+cond+SPACE+QMARK, query.Value) // ex: db.Where("amount > ?", 1000)
	}
}

// Function to manage all conditions
func GetAllCondition(filterQuery string) []func(*gorm.DB) *gorm.DB {
	// Parse filters
	filters := TableFilters{}
	print(filterQuery)
	err := json.Unmarshal([]byte(filterQuery), &filters)
	if err != nil {
		fmt.Println("Error in parsing filters", err.Error())
		return nil
	}

	// Scopes
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)
	// Conditions
	for _, filters := range filters.Queries {
		// validate Operator
		if !validateOperator(filters.Condition) {
			continue //Ignore incase of invalid operator
		}
		// TODO: validate value if it contains SQL injections.
		scopes = append(scopes, GetWhereScope(filters))
	}

	// Sorting
	for _, sort := range filters.Sort {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Order(sort.Col + SPACE + sort.Order) // ex: db.Order("amount desc")
		})
	}

	return scopes
}
