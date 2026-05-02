# convert

A cross-platform CLI tool that converts `.docx` and `.pdf` files to Markdown using [Pandoc](https://pandoc.org), with automatic image extraction and path normalisation.

## Requirements

- Go 1.22+
- [Pandoc](https://pandoc.org/installing.html) in `PATH`

## Build

```bash
make build
# or
go build -o convert ./cmd/convert
```

Cross-platform binaries:

```bash
make cross
```

## Usage

### Single file

```bash
convert input.docx output.md images/
convert report.pdf report.md images/
```

Extracts images to `images/<input_basename>/image1.ext`, `image2.ext`, … and rewrites paths in the Markdown.

### Directory mode

```bash
convert dir input_dir/ output_dir/ images/
```

Recursively converts all `.docx` and `.pdf` files. Preserves subdirectory structure in the output.

### ZIP output

```bash
convert dir input_dir/ output.zip images/ --zip
```

Same as directory mode but packages the result into a zip archive. No temp files are left behind.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-format` | `markdown` | Pandoc output format (`markdown` or `gfm`) |
| `-overwrite` | false | Overwrite existing output files |
| `-recursive` | true | Recurse into subdirectories |
| `-flatten` | false | Collapse directory structure in output |
| `-log <file>` | | Write log output to a file |

## Architecture

```
cmd/convert/         main entry point
internal/
  cli/               argument parsing, command dispatch
  converter/         single-file and directory conversion pipelines
  pandoc/            pandoc invocation via os/exec
  images/            image normalisation and Markdown path rewriting
  filesystem/        path helpers, directory walking, safe temp dirs
  zip/               zip archive creation
  logger/            structured levelled logging
pkg/types/           shared types (ConvertOptions, ConvertResult)
```

The Go application is a thin orchestrator. All document parsing is delegated to Pandoc via `os/exec`. The tool's responsibilities are:

- CLI interface and flag parsing
- Input validation and directory traversal
- Invoking Pandoc with the correct arguments
- Moving extracted images to canonical paths
- Rewriting image references in the generated Markdown
- ZIP packaging and temp directory cleanup

## Testing

```bash
make test
# or
go test ./... -race
```
