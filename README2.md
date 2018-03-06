# work-together

内部系统 - 接口服务

# 工程说明

- run

```sh
bee run -gendoc=true
```

- build version

golang 1.6.2+
Beego     : 1.7.+
Bee     : 1.5.+

- 包结构

```sh
$GO_PATH/src/kuaifa.com/kuaifa/work-together
mkdir -p $GO_PATH/src/kuaifa.com/kuaifa/
```

- 依赖库

使用godep 轻度管理

# 包管理使用 godep

```sh
go get -u -v github.com/tools/godep
```

在`$GOPATH`的`bin/`目录下会有一个`godep`可执行的二进制文件
将`$GOPATH/bin/`加入环境变量

依赖管理使用的文件目录是 `Godeps` **不要手动修改里面的内容**

## 拉取最新依赖

```sh
godep restore
```

## 增加或者更新依赖

```sh
godep save
```


具体使用文档见 [godep 快速指南](UseGoDep.md)