package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Pototoooo/meterforge/tools/migrate/viewgen"
)

func main() {
	var (
		schemaPath = flag.String("schema", "./meterforge/ent/schema", "path to the ent schema package")
		outPath    = flag.String("out", viewgen.DefaultOutputPath, "output SQL file path")
	)
	flag.Parse()

	if err := viewgen.GenerateFile(*schemaPath, *outPath); err != nil {
		exitf("%v", err)
	}
}

func exitf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
