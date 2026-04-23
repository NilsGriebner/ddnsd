/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"ddnsd/cmd"
	"os"

	"github.com/rs/zerolog/log"
)

func main() {

	err := cmd.NewRootCmd().Execute()
	if err != nil {
		log.Error().Err(err)
		os.Exit(1)
	}
}
