# Gorm Utils
```bash
go get github.com/andiwork/gorm-utils
```
## Overview
* model.go : update  base gorm.Model with uuid for ID and set update time to now() and set 

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
## License

© J. K. Gaglo, 2021~time.Now

Released under the [Apache License Version 2.0](https://www.apache.org/licenses/LICENSE-2.0.txt)

