package main

import (
	"github.com/airenas/get-text/internal/pkg/extractor"
	"github.com/airenas/get-text/internal/pkg/file"
	"github.com/airenas/get-text/internal/pkg/service"
	"github.com/airenas/go-app/pkg/goapp"
	"github.com/labstack/gommon/color"
	"github.com/pkg/errors"
)

func main() {
	goapp.StartWithDefault()

	data := service.Data{}
	data.Port = goapp.Config.GetInt("port")

	var err error
	goapp.Log.Infof("Temp dir: %s", goapp.Config.GetString("tempDir"))
	data.Saver, err = file.NewSaver(goapp.Config.GetString("tempDir"))
	if err != nil {
		goapp.Log.Fatal(errors.Wrap(err, "can't init file saver"))
	}
	data.Extractor, err = extractor.NewEBookConverter()
	if err != nil {
		goapp.Log.Fatal(errors.Wrap(err, "can't init ebook extractor wrapper"))
	}

	printBanner()

	err = service.StartWebServer(&data)
	if err != nil {
		goapp.Log.Fatal(errors.Wrap(err, "can't start the service"))
	}
}

var (
	version = "DEV"
)

func printBanner() {
	banner := `
                __        __            __ 
    ____ ____  / /_      / /____  _  __/ /_
   / __ ` + "`" + `/ _ \/ __/_____/ __/ _ \| |/_/ __/
  / /_/ /  __/ /_/_____/ /_/  __/>  </ /_  
  \__, /\___/\__/      \__/\___/_/|_|\__/  v: %s    
 /____/   

%s
________________________________________________________                                                 

`
	cl := color.New()
	cl.Printf(banner, cl.Red(version), cl.Green("https://github.com/airenas/get-text"))
}
