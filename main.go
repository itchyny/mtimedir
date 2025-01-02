package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const name = "mtimedir"

var verbose bool

func main() {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.BoolVar(&verbose, "verbose", false, "enable verbose output")
	if err := fs.Parse(os.Args[1:]); err != nil {
		if err == flag.ErrHelp {
			return
		}
		os.Exit(1)
	}
	args := fs.Args()
	if len(args) == 0 {
		args = []string{"."}
	}
	for _, arg := range args {
		if _, err := run(arg); err != nil {
			if verbose {
				fmt.Fprintf(os.Stderr, "%s: %s\n", name, err)
			}
			os.Exit(1)
		}
	}
}

func run(dir string) (*time.Time, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	fis, err := f.Readdir(10000)
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}
	var maxTime time.Time
	var maxTimeFile string
	for _, fi := range fis {
		path := filepath.Join(dir, fi.Name())
		if fi.IsDir() {
			t, err := run(path)
			if err != nil {
				return nil, err
			}
			if t != nil && maxTime.Before(*t) {
				maxTime = *t
				maxTimeFile = path
			}
		} else if maxTime.Before(fi.ModTime()) {
			maxTime = fi.ModTime()
			maxTimeFile = path
		}
	}
	if !maxTime.IsZero() && stat.ModTime().Unix() != maxTime.Unix() {
		if verbose {
			fmt.Printf("touching %s (<< %s): %s (<= %s)\n", dir, maxTimeFile, maxTime, stat.ModTime())
		}
		if err := os.Chtimes(dir, maxTime, maxTime); err != nil {
			if verbose {
				fmt.Printf("skipping: %s\n", err)
			}
		}
	}
	return &maxTime, nil
}
