package main

import "log"

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	a := App{}
	a.Initialize("postgres://godevman:godevman@localhost/godevman?sslmode=verify-full")

	a.Run(":48888")
}
