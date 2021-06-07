package generator

import (
	"fmt"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const nixExportsBinPath = "/usr/local/prometheus-exporters/bin"

var (
	reTrailingCommaObject = regexp.MustCompile(`,(\s)*}`)
	reTrailingCommaList   = regexp.MustCompile(`,(\s)*]`)
)

func removeTrailingCommas(content string) string {
	content = reTrailingCommaObject.ReplaceAllString(content, `}`)
	return reTrailingCommaList.ReplaceAllString(content, `]`)
}

func compactTextInOneLine(content string) string {
	return strings.ReplaceAll(content, "\n", "")
}

func prometheusExportersBinPath(name string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join("C:\\Program Files\\Prometheus-exporters\\bin", fmt.Sprintf("%s-exporter.exe", name))
	}
	return filepath.Join(nixExportsBinPath, fmt.Sprintf("%s-exporter", name))
}
