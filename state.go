package main

import (
	"github.com/morauszkia/gator/internal/config"
	"github.com/morauszkia/gator/internal/database"
)

type state struct {
	db 		*database.Queries
	config 	*config.Config;
}