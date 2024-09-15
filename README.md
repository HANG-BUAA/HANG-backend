# HANG-backend
“小航书”项目后端仓库

## 项目本地启动方式

### 1. 启动mysql

暂时使用`mysql`作为rdb，后面可能会换。

本地启动`mysql`数据库服务，创建一个新的数据库，名为`hang`，字符集选择`utf8mb4`

或者用`docker`启动（推荐），下面的代码默认使用`root`用户

~~~bash
docker run --name mysql-container \  # 容器名
  -e MYSQL_ROOT_PASSWORD=my-secret-pw \  # 替换为你要设置的密码
  -e MYSQL_DATABASE=hang \
  -p 3306:3306 \  # 暴露端口，建议直接映射3306
  -d mysql:latest

~~~

### 2. 启动redis

按照网上教程在本地配置redis环境，不使用密码。

或者采用docker部署：

~~~bash
docker run --name redis-server \  # 容器名
  -p 6379:6379 \  # 暴露端口，建议直接映射6397
  -d redis
~~~

### 3. 本地配置go环境

略，要求设置`go mod`模式，并设置合适的`GOPROXY`，如阿里云等

可以用`go mod env`查看`go`全部环境配置

### 4. 配置自己本地的邮箱smtp服务

根据自己使用的邮箱不同，到各自官网上进行配置，记住生成的秘钥，要在下面的配置文件里进行配置。

### 5. 配置文件

创建`./src/config/settings.yml`文件，并修改为自己的配置，样例如下：

~~~yaml
mode:
  develop: true # 调试模式

server:
  port: 8000 # 服务启动的端口号

db:
  mysql:
    username: root  # 用户名
    password: 123456  # 密码
    host: 127.0.0.1  # 服务器地址
    port: 3306  # 暴露端口
    dbname: hang  # 数据库名
    maxIdleConn: 10 # 最多空闲链接数
    maxOpenConn: 10 # 最多打开连接数

  redis:
    host: localhost  # 服务器地址
    port: 6379  # 暴露端口

jwt:
  tokenExpire: 30 # token有效时长（分钟）
  signingKey: hang.com # 签名使用的 key

log:
  MaxSize: 1 # 日志文件最大的尺寸(M)，超限后自动分割
  MaxBackups: 10 # 保留旧文件的最大个数
  MaxAge: 90 # 保留旧文件的最大天数
  Compress: false # 是否开启压缩

smtp:
  user_addr: 18646154381@163.com  # 提供服务的邮箱地址
  smtp_addr: smtp.163.com:25  # smtp 服务器地址+端口
  smtp_host: smtp.163.com  # smtp 服务器主机地址
  smtp_password: your_password  # smtp 秘钥
  expiration: 5  # 过期时间（分钟）
~~~

### 6. 启动服务

在顶层目录下运行：

~~~bash
go mod download
go install github.com/swaggo/swag/cmd/swag@latest

cd ./src
swag init

go run .
~~~

即可启动服务

### 7. 检验

1. http检验：浏览器访问`http:localhost:8000/ping`，若返回`pong`，则说明服务启动成功；
2. swag检验：浏览器访问`http:lcoalhost:8000/swagger/index.html`，出现swagger接口文档说明配置成功；
3. redis检验：按照swagger或apifox里的接口文档给自己的北航邮箱发送验证码，发送成功且后端不报错说明链接成功；（目前我已经准备弃用swagger了，请使用apifox）
4. 数据库检验：运行swagger或apifox里的注册与登录接口，状态码200则说明数据库连接成功，也可连接到数据库进行查看检验。（同上）

