package cli

import (
	"flag"
	"fmt"
	"os"
)

// globalFlags holds flags common to all modes.
type globalFlags struct {
	logFile   string
	format    string
	overwrite bool
	flatten   bool
	recursive bool
}

var gf globalFlags

// Execute is the main entry point.
func Execute(version string) {
	flag.StringVar(&gf.logFile, "log", "", "write logs to `file`")
	flag.StringVar(&gf.format, "format", "markdown", "output format: markdown or gfm")
	flag.BoolVar(&gf.overwrite, "overwrite", false, "overwrite existing output files")
	flag.BoolVar(&gf.flatten, "flatten", false, "ignore directory structure in output")
	flag.BoolVar(&gf.recursive, "recursive", true, "recursively process subdirectories")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "convert %s - Convert .docx and .pdf files to Markdown\n\n", version)
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, "  convert [flags] <input> <output.md> <images_dir>")
		fmt.Fprintln(os.Stderr, "  convert [flags] dir <input_dir> <output_dir> <images_dir>")
		fmt.Fprintln(os.Stderr, "  convert [flags] dir <input_dir> <output.zip> <images_dir> --zip")
		fmt.Fprintln(os.Stderr, "\nFlags:")
		flag.PrintDefaults()
	}

	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	var err error
	switch args[0] {
	case "version":
		fmt.Printf("convert %s\n", version)
	case "dir":
		err = runDirCmd(args[1:])
	default:
		err = runSingleCmd(args)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
