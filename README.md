## 启用gomod

### Linux

export GO111MODULE=on

export GOPROXY=https://mirrors.aliyun.com/goproxy/

### Window
自行添加环境变量


### 编译

git clone https://github.com/hongyuefan/hrserver.git

cd hrserver/cmd

go build

### 配置文件

如果不需要连接数据库，编辑字段为空：
"db_url":""

### 运行

./cmd -c config.json

### 测试

http://127.0.0.1:8080/v1/get
