package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/etsxxx/desqueeze/internal/jpeg"
	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:        "out",
		Aliases:     []string{"o"},
		Usage:       "specify the output file path",
		Value:       "",
		DefaultText: "input file path with suffix '-desqueeze' (ex: input.jpg -> input-desqueeze.jpg)",
	},
	&cli.Float64Flag{
		Name:    "multiply",
		Aliases: []string{"m"},
		Usage:   "multiplies the image width by the specified number",
		Value:   1.33,
	},
	&cli.IntFlag{
		Name:    "quality",
		Aliases: []string{"q"},
		Usage:   "specify the output file quality",
		Value:   100,
	},
	&cli.BoolFlag{
		Name:    "overwrite",
		Aliases: []string{"O"},
		Usage:   "overwrite input file",
		Value:   false,
	},
}

func Action(c *cli.Context) error {
	switch {
	case c.Args().Len() <= 0:
		return fmt.Errorf("input file path is required")
	case c.Args().Len() > 1:
		return fmt.Errorf("too many arguments")
	}

	inputPath := c.Args().Get(0)
	var outputPath string
	switch {
	case c.String("out") != "":
		outputPath = c.String("out")
	case c.Bool("overwrite") && c.String("out") == "":
		outputPath = inputPath
	case !c.Bool("overwrite") && c.String("out") == "":
		ext := filepath.Ext(inputPath)
		outputPath = fmt.Sprintf("%s-desqueeze%s", inputPath[0:len(inputPath)-len(ext)], ext)
	default:
		return fmt.Errorf("set output path with '-o' option or use '--overwrite' option")
	}
	quality := c.Int("quality")
	multiply := c.Float64("multiply")

	switch {
	case strings.HasSuffix(strings.ToLower(inputPath), ".jpeg"):
		fallthrough
	case strings.HasSuffix(strings.ToLower(inputPath), ".jpg"):
		fmt.Printf("Process '%s' -> '%s' with desqueeze x%.2f, quality %d\n",
			inputPath, outputPath, multiply, quality)
		return jpeg.Desqueeze(inputPath, outputPath, multiply, quality)
	default:
		return fmt.Errorf("'%s' is not supported type", inputPath)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "desqueeze"
	app.Usage = "convert image width"
	app.Version = fmt.Sprintf("%s (rev:%s)", version, gitcommit)
	app.Flags = Flags
	app.Action = Action

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}
}
