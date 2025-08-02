package main

import (
	"github.com/geowa4/rara-membership/cmd"
	"github.com/geowa4/rara-membership/database"
)

func main() {
	database.Init()
	defer database.Close()

	cmd.Execute()
}
