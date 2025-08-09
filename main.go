package main

import (
	"github.com/geowa4/rara-membership/cmd"
	"github.com/geowa4/rara-membership/database"
	_ "github.com/geowa4/rara-membership/ent/runtime"
)

func main() {
	database.Init()
	defer database.Close()

	cmd.Execute()
}
