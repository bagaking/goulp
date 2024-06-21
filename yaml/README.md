# YAML 文件处理库

## 简介

这个库提供了读取和解析 YAML 文件的功能，并支持通过 `!include` 标签包含其他文件。它可以帮助开发者更方便地处理复杂的 YAML 配置文件。

## 功能

- 读取并解析 YAML 文件
- 支持 `!include` 标签，用于包含其他 YAML 文件
- 自定义文件读取函数

## 安装

使用 `go get` 命令安装：

```
go get -u github.com/bagaking/goulp
```

## 使用方法

### 读取并解析 YAML 文件

```go
package main

import (
    "fmt"
    "log"
    "github.com/bagaking/goulp/yaml"
)

type Config struct {
    // 定义你的配置结构体
}

func main() {
    var config Config
    err := yaml.LoadYAMLFile("config.yaml", &config)
    if err != nil {
        log.Fatalf("error: %v", err)
    }
    fmt.Printf("Parsed config: %+v\n", config)
}
```

### 支持 `!include` 标签

在你的 YAML 文件中，可以使用 `!include` 标签来包含其他文件：

```yaml
database:
  host: localhost
  port: 5432
  credentials: !include credentials.yaml
```

## 示例

假设有以下两个文件：

**config.yaml**

```yaml
database:
  host: localhost
  port: 5432
  credentials: !include credentials.yaml
```

**credentials.yaml**

```yaml
username: admin
password: secret
```

使用 `LoadYAMLFile` 函数读取 `config.yaml` 文件后，`credentials.yaml` 文件的内容会被包含进来。

## 贡献

欢迎贡献代码！请提交 Pull Request 或报告 Issue。

## 许可证

该项目使用 MIT 许可证。详情请参阅 [LICENSE](./LICENSE) 文件。