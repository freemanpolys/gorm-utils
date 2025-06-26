# gorm-utils

A set of utilities for GORM to simplify model definition, pagination, and filtering in Go web applications.

## Installation

```bash
go get github.com/freemanpolys/gorm-utils
```

---

## Features

- **Enhanced GORM Model**: Use UUIDs for IDs and automatic update timestamps.
- **Pagination**: Easily paginate your GORM queries.
- **Filtering**: Build dynamic filters from HTTP requests or query strings.
- **Safe Querying**: Built-in sanitization for fields and operators to prevent SQL injection.

---

## Model Example

```go
package main

import (
	utils "github.com/freemanpolys/gorm-utils"
)

type Person struct {
	utils.Model
	Name string `json:"name"`
	Age  int    `json:"age"`
}
```

---

## Pagination and Filtering with GORM

> **Pagination approach inspired by [this article by Rafael G. Firmino](https://dev.to/rafaelgfirmino/pagination-using-gorm-scopes-3k5f).**

```go
package main

import (
	"fmt"
	utils "github.com/freemanpolys/gorm-utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	db.AutoMigrate(&Person{})

	// Create sample data
	db.Create(&Person{Name: "John", Age: 25})
	db.Create(&Person{Name: "Jane", Age: 30})
	db.Create(&Person{Name: "Doe", Age: 25})

	var people []Person

	// Pagination
	pagination := utils.Pagination{Limit: 2, Page: 1, Sort: "age desc"}
	db.Scopes(utils.Paginate(&pagination)).Find(&people)
	fmt.Println("Paginated:", people)

	// Filtering
	filters := []utils.Filter{
		{Field: "age", Operator: "=", Value: 25},
	}
	db.Scopes(utils.FilterScope(filters)).Find(&people)
	fmt.Println("Filtered:", people)

	// Combine Pagination and Filtering
	db.Scopes(utils.Paginate(&pagination), utils.FilterScope(filters)).Find(&people)
	fmt.Println("Paginated & Filtered:", people)

	// String to Filters
	filtersStr := "age=25&name=John"
	filters, err := utils.StringToFilters(filtersStr)
	if err != nil {
		fmt.Println("Error converting string to filters:", err)
		return
	}
	db.Scopes(utils.FilterScope(filters)).Find(&people)
	fmt.Println("StringToFilters:", people)
}
```

---

## HTTP Query Example

You can extract pagination and filters directly from an HTTP request using `FromRequest`:

### Example HTTP Request

```
GET /products?page=1&limit=10&sort=qty asc&filter=name::like::foo~~qty::gte::10~~price::between::100||200
```

### Go Usage

```go
import (
	"net/http"
	utils "github.com/freemanpolys/gorm-utils"
)

func handler(w http.ResponseWriter, r *http.Request) {
	pagination, filters := utils.FromRequest(r)
	// Use with GORM
	var products []Product
	db.Scopes(utils.Paginate(pagination), utils.FilterScope(filters)).Find(&products)
	// ...return products as JSON, etc.
}
```

### Special Cases: IN and BETWEEN Operators

To filter with the SQL `IN` operator, use a double-pipe `||` separated list in the value:

```
GET /products?filter=code::in::D42||L12
```

This will be parsed as:

```go
[]utils.Filter{
	{Field: "code", Operator: "in", Value: []string{"D42", "L12"}},
}
```

To filter with the SQL `BETWEEN` operator, use two values separated by `||`:

```
GET /products?filter=price::between::100||200
```

This will be parsed as:

```go
[]utils.Filter{
	{Field: "price", Operator: "between", Value: []string{"100", "200"}},
}
```

---

## Security

- All field names and operators are sanitized.
- Only a safe subset of SQL operators is allowed (`=`, `<>`, `>`, `>=`, `<`, `<=`, `LIKE`, `IN`).

---

## License

© J. K. Gaglo, 2021~time.Now

Released under the [Apache License Version 2.0](https://www.apache.org/licenses/LICENSE-2.0.txt)

