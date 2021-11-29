go-extract-embedfs  
===============
a little tool to extract golang embed files to local path
一个用于将go的embed.FS文件集合提取并存储到本地磁盘的小工具
  
## install 安装
```
go get github.com/del-xiong/go-extract-embedfs
```

## examples

```
package main

import (
	"embed"
)
import extractfs "github.com/del-xiong/go-extract-embedfs"

//go:embed myfiles
var res embed.FS

func main() {
	extractfs.ExtractToPath(&res, `/tmp/`)
}
```
myfiles would be extracted to local /tmp/, the final path is /tmp/myfiles  
myfiles 会提取到本地/tmp/目录，最终路径为 **/tmp/myfiles**

## known issues 目前存在的问题  
some files that contains some special characters wont be embed into embed.FS,this might be a bug of golang  
一些包含特殊字符的文件无法被索引到embed.FS中，这可能是一个golang的bug  
for example:  
例如:  
```
（）
【】
😀😁😂😃😄😅emojis...
```
