# HANG-backend
“小航书”项目后端仓库

## 启动方式

### 1. 启动数据库

本地启动`mysql`数据库服务，创建一个新的数据库，表名为`hang`，字符集选择`utf8mb4`

或者用`docker`启动（推荐），下面的代码默认使用`root`用户

~~~bash
docker run --name mysql-container \
  -e MYSQL_ROOT_PASSWORD=my-secret-pw \  # 替换为你要设置的密码
  -e MYSQL_DATABASE=hang \
  -p 3306:3306 \  # 暴露端口，建议直接映射3306
  -d mysql:latest

~~~

### 2. 本地配置go环境

略，要求设置`go mod`模式，并设置合适的`GOPROXY`

可以用`go mod env`查看`go`全部环境配置

### 3. 修改配置文件

打开`./src/config/settings.yml`文件，并修改为自己的配置，样例如下：

~~~yaml
mode:
  develop: true # 调试模式

server:
  port: 8000 # 服务启动的端口号

db:
  mysql:
    username: root  # 数据库用户名
    password: 123456  # 数据库密码，替换为自己的！
    host: 127.0.0.1  # 数据库地址，本地启动不必修改
    port: 3306  # 暴露的端口
    dbname: hang  # 数据库名称
    maxIdleConn: 10 # 最多空闲链接数
    maxOpenConn: 10 # 最多打开连接数

jwt:
  tokenExpire: 30 # token有效时长（分钟）
  signingKey: hang.com # 签名使用的 key

log:
  MaxSize: 1 # 日志文件最大的尺寸(M)，超限后自动分割
  MaxBackups: 10 # 保留旧文件的最大个数
  MaxAge: 90 # 保留旧文件的最大天数
  Compress: false # 是否开启压缩

~~~

### 4. 启动服务

在顶层目录下运行：

~~~bash
go mod download
go install github.com/swaggo/swag/cmd/swag@latest

cd ./src
swag init

go run .
~~~

即可启动服务

### 5. 检验

1. http检验：浏览器访问`http:localhost:8000/ping`，若返回`pong`，则说明服务启动成功；
2. swag检验：浏览器访问`http:lcoalhost:8000/swagger/index.html`，出现swagger接口文档说明配置成功；
3. 数据库检验：运行apifox里的注册接口，状态码200则说明数据库连接成功，也可连接到数据库进行查看检验。

