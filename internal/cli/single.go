package cli

import (
	"fmt"

	"github.com/convert/internal/converter"
	"github.com/convert/internal/logger"
	"github.com/convert/internal/pandoc"
)

// runSingleCmd handles: convert <input> <output.md> <images_dir>
func runSingleCmd(args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("single-file mode requires: <input> <output.md> <images_dir>\ngot: %v", args)
	}

	if err := logger.Init(gf.logFile, logger.INFO); err != nil {
		return err
	}
	if err := pandoc.CheckAvailable(); err != nil {
		return err
	}

	res, err := converter.ConvertFile(converter.Options{
		InputFile:     args[0],
		OutputFile:    args[1],
		ImagesRootDir: args[2],
		Format:        gf.format,
		Overwrite:     gf.overwrite,
	})
	if err != nil {
		return err
	}

	fmt.Printf("converted: %s → %s (%d image(s))\n", res.InputFile, res.OutputFile, res.ImageCount)
	return nil
}
