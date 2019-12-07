package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const name = "mtimedir"

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		args = []string{"."}
	}
	for _, arg := range args {
		if _, err := run(arg); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", name, err)
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
		fmt.Printf("touching %s (<< %s): %s (<= %s)\n", dir, maxTimeFile, maxTime, stat.ModTime())
		if err := os.Chtimes(dir, maxTime, maxTime); err != nil {
			fmt.Fprintf(os.Stderr, "skipping: %s\n", err)
		}
	}
	return &maxTime, nil
}
