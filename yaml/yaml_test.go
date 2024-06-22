package yaml

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// TestLoadYAML 测试 LoadYAML 函数
func TestLoadYAML(t *testing.T) {
	// 主 YAML 数据
	mainYAML := `
name: main
k2: !include ./test_include.yaml
`

	// 包含的 YAML 数据
	includedYAML := `
user: gogo
i_dont_care: 123
`

	// 模拟文件系统
	fileSystem := map[string][]byte{
		"test_include.yaml": []byte(includedYAML),
	}

	// 模拟文件读取函数
	mockReadFile := func(filename string) ([]byte, error) {
		data, ok := fileSystem[filename]
		if !ok {
			return nil, os.ErrNotExist
		}
		return data, nil
	}

	// 定义解析结果的结构体
	type Config struct {
		Name string `yaml:"name"`
		K2   struct {
			User string `yaml:"user"`
		} `yaml:"k2"`
	}

	var config Config
	if err := LoadYAML([]byte(mainYAML), ".", &config, mockReadFile); err != nil {
		t.Fatalf("Failed to load YAML: %v", err)
	}

	// 验证解析结果
	if config.Name != "main" {
		t.Errorf("Expected name to be 'main', got '%s'", config.Name)
	}
	if config.K2.User != "gogo" {
		t.Errorf("Expected include.key to be 'gogo', got '%s'", config.K2.User)
	}
}

func TestLoadYAMLNestedIncludeUsesIncludedFileDir(t *testing.T) {
	mainYAML := `
name: main
profile: !include configs/profile.yaml
`
	profileYAML := `
user: !include user.yaml
`
	userYAML := `
name: gogo
`

	fileSystem := map[string][]byte{
		"configs/profile.yaml": []byte(profileYAML),
		"configs/user.yaml":    []byte(userYAML),
	}
	mockReadFile := func(filename string) ([]byte, error) {
		data, ok := fileSystem[filename]
		if !ok {
			return nil, os.ErrNotExist
		}
		return data, nil
	}

	type Config struct {
		Name    string `yaml:"name"`
		Profile struct {
			User struct {
				Name string `yaml:"name"`
			} `yaml:"user"`
		} `yaml:"profile"`
	}

	var config Config
	if err := LoadYAML([]byte(mainYAML), ".", &config, mockReadFile); err != nil {
		t.Fatalf("LoadYAML(%q) error = %v, want nil", mainYAML, err)
	}

	if config.Profile.User.Name != "gogo" {
		t.Errorf("LoadYAML(%q) profile.user.name = %q, want %q", mainYAML, config.Profile.User.Name, "gogo")
	}
}

func TestLoadYAMLAbsoluteIncludeBypassesBaseDir(t *testing.T) {
	baseDir := filepath.Join(t.TempDir(), "base")
	includePath := filepath.Join(t.TempDir(), "profile.yaml")
	mainYAML := "profile: !include " + includePath + "\n"
	includedYAML := "user: gogo\n"

	readFilenames := make([]string, 0, 1)
	mockReadFile := func(filename string) ([]byte, error) {
		readFilenames = append(readFilenames, filename)
		if filename != includePath {
			return nil, os.ErrNotExist
		}
		return []byte(includedYAML), nil
	}

	type Config struct {
		Profile struct {
			User string `yaml:"user"`
		} `yaml:"profile"`
	}

	var config Config
	if err := LoadYAML([]byte(mainYAML), baseDir, &config, mockReadFile); err != nil {
		t.Fatalf("LoadYAML(%q, baseDir %q) error = %v, want nil", mainYAML, baseDir, err)
	}

	if config.Profile.User != "gogo" {
		t.Errorf("LoadYAML(%q, baseDir %q) profile.user = %q, want %q", mainYAML, baseDir, config.Profile.User, "gogo")
	}
	if len(readFilenames) != 1 {
		t.Fatalf("LoadYAML(%q, baseDir %q) read filenames count = %d, want %d", mainYAML, baseDir, len(readFilenames), 1)
	}
	if readFilenames[0] != includePath {
		t.Errorf("LoadYAML(%q, baseDir %q) read filename = %q, want %q", mainYAML, baseDir, readFilenames[0], includePath)
	}
}

func TestLoadYAMLReturnsReadErrorForMissingInclude(t *testing.T) {
	mainYAML := `
profile: !include missing.yaml
`

	var readFilenames []string
	mockReadFile := func(filename string) ([]byte, error) {
		readFilenames = append(readFilenames, filename)
		return nil, os.ErrNotExist
	}

	type Config struct {
		Profile struct {
			User string `yaml:"user"`
		} `yaml:"profile"`
	}

	var config Config
	err := LoadYAML([]byte(mainYAML), ".", &config, mockReadFile)
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("LoadYAML(%q) error = %v, want errors.Is(err, %v)", mainYAML, err, os.ErrNotExist)
	}

	if len(readFilenames) != 1 {
		t.Fatalf("LoadYAML(%q) read filenames count = %d, want %d", mainYAML, len(readFilenames), 1)
	}
	if readFilenames[0] != "missing.yaml" {
		t.Errorf("LoadYAML(%q) read filename = %q, want %q", mainYAML, readFilenames[0], "missing.yaml")
	}
}

func TestLoadYAMLIgnoresUnknownFieldsAfterIncludeExpansion(t *testing.T) {
	mainYAML := `
name: main
unknown_root: ignored
profile: !include profile.yaml
`
	profileYAML := `
user: gogo
unknown_profile: ignored
`

	fileSystem := map[string][]byte{
		"profile.yaml": []byte(profileYAML),
	}
	mockReadFile := func(filename string) ([]byte, error) {
		data, ok := fileSystem[filename]
		if !ok {
			return nil, os.ErrNotExist
		}
		return data, nil
	}

	type Config struct {
		Name    string `yaml:"name"`
		Profile struct {
			User string `yaml:"user"`
		} `yaml:"profile"`
	}

	var config Config
	if err := LoadYAML([]byte(mainYAML), ".", &config, mockReadFile); err != nil {
		t.Fatalf("LoadYAML(%q) error = %v, want nil", mainYAML, err)
	}

	if config.Name != "main" {
		t.Errorf("LoadYAML(%q) name = %q, want %q", mainYAML, config.Name, "main")
	}
	if config.Profile.User != "gogo" {
		t.Errorf("LoadYAML(%q) profile.user = %q, want %q", mainYAML, config.Profile.User, "gogo")
	}
}
