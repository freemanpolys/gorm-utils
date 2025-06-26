# Gorm Utils
```bash
go get github.com/andiwork/gorm-utils
```
## Overview
* model.go : update  base gorm.Model with uuid for ID and set update time to now() 

## Getting Started
### gorm-utils.Model
```
package hello

import (
	utils "github.com/andiwork/gorm-utils"
)

// User is just a sample type
type Person struct {
	utils.Model
	Name string `json:"name" description:"name of the user" default:"john"`
	Age  int    `json:"age" description:"age of the user" default:"21"`
}
```

### Pagination and Filtering

This utility also provides a simple way to paginate and filter your GORM queries.

```go
package main

import (
	"fmt"

	utils "github.com/andiwork/gorm-utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Person{})

	// Create
	db.Create(&Person{Name: "John", Age: 25})
	db.Create(&Person{Name: "Jane", Age: 30})
	db.Create(&Person{Name: "Doe", Age: 25})

	var people []Person

	// Paginate
	pagination := utils.Pagination{Limit: 2, Page: 1, Sort: "age desc"}
	db.Scopes(utils.Paginate(&pagination)).Find(&people)
	fmt.Println(people)

	// Filter
	filters := []utils.Filter{
		{Field: "age", Operator: "=", Value: 25},
	}
	db.Scopes(utils.FilterScope(filters)).Find(&people)
	fmt.Println(people)
	 // Combine Pagination and Filter
	db.Scopes(utils.Paginate(&pagination), utils.FilterScope(filters)).Find(&people)
	fmt.Println(people)

	// String to StringToFilters
	 filtersStr := "age=25&name=John"
	  filters, err := utils.StringToFilters(filtersStr)
	   if err != nil {
	fmt.Println("Error converting string to filters:", err)
	return
	}
	db.Scopes(utils.FilterScope(filters)).Find(&people)
	fmt.Println(people)
}
```

## License

© J. K. Gaglo, 2021~time.Now

Released under the [Apache License Version 2.0](https://www.apache.org/licenses/LICENSE-2.0.txt)

