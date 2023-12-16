package main

import (
	"time"

	"github.com/rs/zerolog"
)

func init() {

	zerolog.TimeFieldFormat = time.RFC3339Nano

}

func main() {

}
