package main

import (
	_ "github.com/go-sql-driver/mysql"

	"github.com/sliveryou/grom/cmd"
)

func main() {
	cmd.Execute()
}
