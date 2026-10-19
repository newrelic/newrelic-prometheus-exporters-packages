package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	templateFileMask = 0600
)

func getTemplateVars(exporterPort string) map[string]string {
	return map[string]string{
		"exporterPort": exporterPort,
	}
}

func getConfigGeneratorEnvVars(configFileName string) []string {

	path := filepath.Join(rootDir(), "integration-tests", "testdata", configFileName)
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

func copyTemplate(fileName string) error {
	sourcePath := filepath.Join("testdata", "integration_template", fileName)
	bytesRead, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}
	targetPath := filepath.Join(rootDir(), "templates", fileName)
	return os.WriteFile(targetPath, bytesRead, templateFileMask)
}

func buildGeneratorConfig(integration string, integrationVersion string) error {
	integrationTemplateFileName := fmt.Sprintf("%s.json.tmpl", integration)
	prometheusTemplateFileName := fmt.Sprintf("%s.prometheus.json.tmpl", integration)

	if err := copyTemplate(integrationTemplateFileName); err != nil {
		return err
	}

	if err := copyTemplate(prometheusTemplateFileName); err != nil {
		return err
	}

	return nil
}

// buildGenerator compiles the config generator into binDir, so its executable folder is deterministic
// (`go run` would use a throwaway build cache folder).
func buildGenerator(binDir string) (string, error) {
	genName := "nri-config-generator"
	if runtime.GOOS == "windows" {
		genName += ".exe"
	}

	genPath := filepath.Join(binDir, genName)
	cmd := exec.Command(
		"go",
		"build",
		"-ldflags",
		fmt.Sprintf("-X main.integration=%s -X main.integrationVersion=%s", testIntegration, testIntegrationVersion),
		"-o", genPath,
		"../main.go",
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("building generator: %w: %s", err, out)
	}

	return genPath, nil
}

func callGeneratorConfig(generatorBin string, args []string, env []string) ([]byte, error) {
	cmd := exec.Command(generatorBin, args...)

	cmd.Env = append(cmd.Environ(), env...)

	return cmd.Output()
}
