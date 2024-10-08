package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/diegomagalhaes-dev/go-service/app/tooling/docs/output/html"
	"github.com/diegomagalhaes-dev/go-service/app/tooling/docs/output/json"
	"github.com/diegomagalhaes-dev/go-service/app/tooling/docs/output/text"
	"github.com/diegomagalhaes-dev/go-service/app/tooling/docs/webapi"
)

var output = flag.String("out", "html", "json, text, html")
var browser = flag.Bool("browser", false, "start the browser automagically")

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	routes, err := webapi.Routes("v1")
	if err != nil {
		return fmt.Errorf("webapi.Routes, %w", err)
	}

	records, err := webapi.Records(routes)
	if err != nil {
		return fmt.Errorf("webapi.Records, %w", err)
	}

	switch *output {
	case "text":
		err = text.Transform(records)

	case "json":
		err = json.Transform(records)

	case "html":
		err = html.Transform(records, *browser)
	}

	if err != nil {
		return fmt.Errorf("transform, %w", err)
	}

	return nil
}
