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
