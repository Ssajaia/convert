package types

// ConvertOptions holds all options for a conversion operation.
type ConvertOptions struct {
	InputPath      string
	OutputPath     string
	ImagesDir      string
	Recursive      bool
	Overwrite      bool
	Flatten        bool
	LogFile        string
	Format         string
	ImagesRelative bool
	ImagesAbsolute bool
	ZipOutput      bool
}

// ConvertResult represents the outcome of a single file conversion.
type ConvertResult struct {
	InputFile  string
	OutputFile string
	ImagesDir  string
	ImageCount int
	Error      error
}

// SupportedExtensions is the list of file types the tool can convert.
var SupportedExtensions = map[string]bool{
	".docx": true,
	".pdf":  true,
	".odt": true,
}
