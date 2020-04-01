package main

import (
	"fmt"
	"github.com/ranabd36/project-qa/config"
)

func main() {
	fmt.Println(config.Get().Database.DatabasePort)
}
