# GOIMV 
GOIMV 是用于快速搭建go言语的刷题站点。
## 功能介绍
- 基于codemirror、gocode 和gofmt 开发的代码编辑器。支持golang语法的高亮、格式化、和代码补全。
- 基于golang 开发的后端服务，用于支撑前端展示、题目添加、下发
- 基于docker 开发的代码运行环境，将运行环境进行隔离。防止提交代码宿主机影响，便于最资源隔离、控制以及回收
## 工程总体架构
![avatar] (https://img.hacpai.com/file/2020/01/c0e73ba31b960ab9098d07b1dac0f47-2e1e9ec7.png)
## 代码目录结构
```shell
----- goenv
|       |---- autocommplete.go 自动补全相关的代码
|       |---- goenv.go   golang运行环境，负责编译、运行和测试相关功能
|
|---- problem
|       |---- problem.go 负责解析题目，题目列表下发
|
|---- reponse
|        |---- response.go http数据下发的统一包装类
|
|---- views
|       |---- index.html   前端展示页面
|       |---- view.go       前端渲染相关代码
|          
|---- static
|       |---- codemirror 开源js包，用于支持web编辑器
|       |---- js       前端相关的js代码
|
|---- problems
        |--- {$题目名称}
        |        |---- content 题目描述，支持markdown格式
        ……       |---- template 题目模板，用于编写好与展示的代码
                 |---- problem_test.go 单测用例，用于检验题目是否正常
```

## 安装
### 环境部署
相关依赖，确保你已经安装了 golang环境，docker 环境， gocode 和gofmt相关组件
1、golang docker container下载
```golang
docker pull golang
```
2、代码下载
```shell
cd $GOPATH
git clone ....
go run main.go
```

### 题库添加
所有的题目放在工程的problems下，当你想新增一个题目，需要做如下几件事
1. 在problems目录下创建你的子目录，文件名称是你的题目名称
2. 添加题目描述文件，命名为content
3. 添加答题模板文件，命名为template
4. 添加单测文件，写入你单测用例，命名为 problem_test.go

### TODO
1. 身为后端程序，前端页面并不美观，需要调整
2. 没有用户体系，无法提供答题历史记录问题
3. 部分代码逻辑可以抽出成配置，入webserver 端口号。超时时间等
4. 开发时间较短，工程目录设计，代码结构稍显混乱，如果有项目有人关注，后续重构。
