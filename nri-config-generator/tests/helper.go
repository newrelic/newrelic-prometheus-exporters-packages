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

func getTemplateVars(exporterPort string) map[string]string {
	return map[string]string{
		"exporterPort": exporterPort,
	}
}

func getConfigGeneratorEnvVars(configFileName string) []string {

	path := filepath.Join(rootDir(), "tests", "testdata", configFileName)
	return []string{
		fmt.Sprintf("CONFIG_PATH=%s", path),
	}
}

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
func buildGeneratorConfig(integration string, integrationVersion string) error {
	if err := copyIntegrationTemplate(integration); err != nil {
		return err
	}
	cmd := &exec.Cmd{
		Path: "/usr/bin/make",
		Args: []string{
			"make",
			"compile",
			fmt.Sprintf("PACKAGE_NAME=%s", integration),
			fmt.Sprintf("VERSION=%s", integrationVersion),
		},
		Dir: rootDir(),
	}
	return cmd.Run()
}

func fetchDefinitions(integration string) error {
	sourceFile := filepath.Join("testdata", fmt.Sprintf("%s-definitions.yml", integration))
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return err
	}

	destinationFile := filepath.Join(rootDir(), "definitions", "definitions.yml")
	return ioutil.WriteFile(destinationFile, input, 0777)
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
	executable := fmt.Sprintf("%s", integration)
	ctx, _ := context.WithTimeout(context.Background(), 5000*time.Millisecond)
	name := filepath.Join(rootDir(), "bin", executable)
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Env = env
	return cmd.Output()
}
