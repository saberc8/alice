该目录用于存放由 swag 工具生成的 swagger 文档文件 (docs.go, swagger.json, swagger.yaml)。

生成命令示例:

```
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g main.go -o docs
```

注意: 生成后的文件应提交到版本库 (或视团队策略决定)。
