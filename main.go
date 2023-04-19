package main

import (
	"log"

	"github.com/aretaja/godevmanapi/app"
)

const version string = "v0.0.1-devel.1"

// @title goDevmanAPI
// @version v0.0.1-devel.1
// @description goDevmans API

// @contact.name API Support
// @contact.email marko@aretaja.org

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	a := new(app.App)
	a.Version = version
	a.Initialize()

	a.Run()
}
