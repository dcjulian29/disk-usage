package main

import (
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type directoryEntity struct {
	Path  string
	Bytes int64
}

func main() {
	root := "."

	entries, err := os.ReadDir(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading directory:", err)
		os.Exit(1)
	}

	var results []directoryEntity

	rootSize := getSize(root)
	results = append(results, directoryEntity{Path: root, Bytes: rootSize})

	for _, entry := range entries {
		fullPath := filepath.Join(root, entry.Name())
		size := int64(0)
		if entry.IsDir() {
			size = getSize(fullPath)
		} else {
			info, err := entry.Info()
			if err == nil {
				size = info.Size()
			}
		}

		results = append(results, directoryEntity{Path: fullPath, Bytes: size})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Bytes > results[j].Bytes
	})

	for _, entry := range results {
		fmt.Printf("%s %-8s\t%s\n", drawBar(entry.Bytes, results[0].Bytes), humanReadable(entry.Bytes), entry.Path)
	}
}

func drawBar(value int64, max int64) string {
	const barWidth = 20

	if max == 0 {
		return ""
	}

	logValue := math.Log10(float64(value) + 1)
	logMax := math.Log10(float64(max) + 1)
	ratio := logValue / logMax
	filled := int(ratio * float64(barWidth))

	if filled > barWidth {
		filled = barWidth
	}

	return strings.Repeat("▄", filled) + strings.Repeat(" ", barWidth-filled)
}

func getSize(path string) int64 {
	var total int64

	err := filepath.WalkDir(path, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Ignore permission denied or other errors
		}

		if !d.IsDir() {
			info, err := d.Info()
			if err == nil {
				total += info.Size()
			}
		}

		return nil
	})

	if err != nil {
		// Fail silently
	}

	return total
}

func humanReadable(size int64) string {
	const unit = 1024

	if size < unit {
		return fmt.Sprintf("%dB", size)
	}

	div, exp := int64(unit), 0

	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f%cB", float64(size)/float64(div), "KMGTPE"[exp])
}
