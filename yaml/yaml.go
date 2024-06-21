package yaml

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// FileReader 定义文件读取函数类型
type FileReader func(filename string) ([]byte, error)

// LoadYAMLFile 读取并解析 YAML 文件
// filename: 要读取的 YAML 文件的路径
// out: 用于存储解析结果的结构体指针
func LoadYAMLFile(filename string, out any) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// 使用默认的文件读取函数 os.ReadFile
	return LoadYAML(data, filepath.Dir(filename), out, os.ReadFile)
}

// LoadYAML 读取并解析 YAML 数据
// data: 要解析的 YAML 数据
// baseDir: 基础目录，用于解析相对路径
// out: 用于存储解析结果的结构体指针
// readFile: 文件读取函数
func LoadYAML(data []byte, baseDir string, out any, readFile FileReader) error {
	// 自定义解码器
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	for {
		var node yaml.Node
		if err := decoder.Decode(&node); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if err := processNode(&node, baseDir, readFile); err != nil {
			return err
		}

		if err := node.Decode(out); err != nil {
			return err
		}
	}

	return nil
}

// processNode 处理 YAML 节点，支持 !include 标签
// node: 当前处理的 YAML 节点
// baseDir: 基础目录，用于解析相对路径
// readFile: 文件读取函数
// processNode 处理 YAML 节点，支持 !include 标签
func processNode(node *yaml.Node, baseDir string, readFile FileReader) error {
	if node.Kind == yaml.ScalarNode && node.Tag == "!include" {
		includePath := node.Value
		if !filepath.IsAbs(includePath) {
			includePath = filepath.Join(baseDir, includePath)
		}

		data, err := readFile(includePath)

		if err != nil {
			return err
		}

		var includedNode yaml.Node
		if err = yaml.Unmarshal(data, &includedNode); err != nil {
			return err
		}

		*node = includedNode
	}

	for _, child := range node.Content {
		if err := processNode(child, baseDir, readFile); err != nil {
			return err
		}
	}

	return nil
}
