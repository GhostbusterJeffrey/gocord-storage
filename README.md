# gocord-storage
## A Go implementation of [jscord-storage](https://github.com/animemoeus/jscord-storage)

### Installation
`go get github.com/ghostbusterjeffrey/gocord-storage`

### Usage Example
```go
package main

import (
  "fmt"
  "github.com/ghostbusterjeffrey/gocord-storage"
)

func main() {
  res, _ := gocordstorage.Upload("uploaded_file.txt","./file.txt")
  fmt.Printf("%s\n", res)
}
```
