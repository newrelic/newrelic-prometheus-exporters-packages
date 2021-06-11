package test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func rootDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(dir)
}

func copyIntegrationTemplate(integration string) error {
	fileName := fmt.Sprintf("%s.json.tmpl", integration)
	sourcePath := filepath.Join("testdata", "integration_template", fileName)
	bytesRead, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return err
	}
	targetPath := filepath.Join(rootDir(), "templates", fileName)
	return ioutil.WriteFile(targetPath, bytesRead, 0755)
}
func buildGeneratorConfig(integration string, defaultExporterPort string) error {
	if err := copyIntegrationTemplate(integration); err != nil {
		return err
	}
	cmd := &exec.Cmd{
		Path: "/usr/bin/make",
		Args: []string{
			"make",
			"compile",
			fmt.Sprintf("INTEGRATION_NAME=%s", integration),
		},
		Dir: rootDir(),
	}
	return cmd.Run()
}

func clean() error {
	cmd := &exec.Cmd{
		Path: "/usr/bin/make",
		Args: []string{
			"make",
			"clean",
		},
		Dir: rootDir(),
	}
	return cmd.Run()
}

func callGeneratorConfig(integration string, args []string, env []string) ([]byte, error) {
	executable := fmt.Sprintf("nri-%s", integration)
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
	cmd := exec.CommandContext(ctx, filepath.Join(rootDir(), "bin", executable), args...)
	cmd.Env = env
	return cmd.Output()
}
