package yaml

// IncludeTag 用于标记需要包含的文件
type IncludeTag struct {
	Path string `yaml:"path"`
}

// UnmarshalYAML 实现自定义的解码逻辑
func (i *IncludeTag) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var path string
	if err := unmarshal(&path); err != nil {
		return err
	}
	i.Path = path
	return nil
}
