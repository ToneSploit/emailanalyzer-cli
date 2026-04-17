package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tonesploit/emailanalyzer"
)

func main() {
	outDir := flag.String("out", "Analyzed", "output directory")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: analyze-email [flags] <file.eml> ...\n\nFlags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "creating output dir: %v\n", err)
		os.Exit(1)
	}

	for _, path := range flag.Args() {
		if err := process(path, *outDir); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", path, err)
		}
	}
}

func process(path, outDir string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	result, err := emailanalyzer.Analyze(f, filepath.Base(path))
	if err != nil {
		return err
	}

	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	jsonPath := filepath.Join(outDir, base+".json")
	jf, err := os.Create(jsonPath)
	if err != nil {
		return fmt.Errorf("creating json: %w", err)
	}
	defer jf.Close()
	enc := json.NewEncoder(jf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		return fmt.Errorf("writing json: %w", err)
	}

	mdPath := filepath.Join(outDir, base+".md")
	md := emailanalyzer.ToMarkdown(result)
	if err := os.WriteFile(mdPath, []byte(md), 0o644); err != nil {
		return fmt.Errorf("writing md: %w", err)
	}

	fmt.Printf("analyzed: %s → %s, %s\n", path, jsonPath, mdPath)
	return nil
}
