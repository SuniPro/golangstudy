package main

import "nlabsoft__assignment/controller"

func main() {
	r := controller.SetupRouter()
	r.Run(":8080")
}
