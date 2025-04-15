// Package main controls the user interaction logic for the application
package main

import (
	"flag"
	"ptt/pkg/models"
	"ptt/pkg/utils"
)

var version = "1.1.0"
var transformationModeArray models.MultiString

func main() {
	flag.Var(&transformationModeArray, "transformation", "Transformation mode to be used. Can be specified multiple times.")
	flag.Parse()

	utils.ReadReturnStandardInput(transformationModeArray)
}
