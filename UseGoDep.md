# 介绍

godep是解决包依赖的管理工具，目前最主流的一种，原理是扫描记录版本控制的信息，并在go命令前加壳来做到依赖管理

# 安装

```sh
go get -u -v github.com/tools/godep
```

成功安装后，在`$GOPATH的bin目录`下会有一个godep可执行的二进制文件，后面执行的命令都是用这个，建议这个目录加入到全局环境变量中

# 包管理使用 godep

以下命令都是在工程的根目录运行

## 拉取依赖 restore 用于开发


```sh
godep restore
```

建议开发过程使用这个命令来同步依赖库

如果下载的项目中只有Godeps.json文件，而没有包含第三库则可以使用godep restore这个命令将所有的依赖库下来到`$GOPATH\src`中 用于开发

godep restore执行时，godep会按照`Godeps/Godeps.json`内列表，依次执行`go get -d -v`来下载对应依赖包到GOPATH路径下

> 如果某个原先的依赖包保存路径（GOPATH下的相对路径）与下载url路径不一致，比如kuberbetes在github上路径是github.com/kubernetes，而代码内import则是my.io，则会导致无法下载成功，也就是说godep restore不成功。这种只能手动，比如手动创建$GOPATH/my.io目录，然后git clone

## 检出依赖 save

```sh
godep save
```

- 自动扫描当前目录所属包中import的所有外部依赖库（非系统库）
- 将所有的依赖库下来下来到当前工程中，产生文件 `Godeps\Godeps.json` 文件

godep save能否成功执行需要有两个要素：

- 当前或者需扫描的包均能够编译成功：因此所有依赖包事先都应该已经或go get或手工操作保存当前GOPATH路径下
- 依赖包必须使用了某个代码管理工具（如git，hg）：这是因为godep需要记录revision

这个命令用于编译构建的，三方构建工具需要额外配置构建参数

## godep 编译运行 build

项目用godep管理后，要编译和运行项目的时候再用go run和go build显然就不行

> 因为go命令是直接到GOPATH目录下去找第三方库，而使用godep下载的依赖库放到Godeps/workspace目录下的，但是不影响继续使用依赖GOPATH目录，所以与三方工具本身不冲突

故使用

```sh
godep go build XXX
```

> godep中的go命令，就是将原先的go命令加了一层壳，执行godep go的时候，会将当前项目的workspace目录加入GOPATH变量中

## Godeps目录的作用

godep save时godep把所有依赖包代码从GOPATH路径拷贝到Godeps目录下，并去除代码管理目录。这个用处主要是为了支撑godep go tool的一系列操作，尤其是git clone了代码库下来后，通常直接用godep go install xxx即可完成编译，一定程度上能够缓解golang比较严格的代码路径和包管理带来的烦恼。

而在`使用IDE时`，可以通过把Godeps/_workspace添加到GOPATH实现代码跳转和编译等功能，比较方便

# godep其他命令

```sh
    save     list and copy dependencies into Godeps
    go       run the go tool with saved dependencies
    get      download and install packages with specified dependencies
    path     print GOPATH for dependency code
    restore  check out listed dependency versions in GOPATH
    update   update selected packages or the go version
    diff     shows the diff between current and previously saved set of dependencies
    version  show version info
```