package main

import (
	"log"

	"github.com/aretaja/godevmanapi/app"
)

var version string = "v0.0.1-devel.1"

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	a := new(app.App)
	a.Version = version
	a.Initialize()

	a.Run()
}
