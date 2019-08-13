package main

import (
	"github.com/seramirezdev/truora-test/api/application"
)

func main() {
	app := application.App{}
	app.Inicialize("sergio", "truoratest")
	app.Run("3000")
}
