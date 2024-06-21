package yaml

import (
	"os"
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
