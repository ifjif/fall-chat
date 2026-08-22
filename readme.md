# fall-chat 网页聊天

## 使用

      git clone https://github.com/ifjif/fall-chat.git

       在当前目录下创建 .data目录 .logs目录 .store目录
       创建 .env文件
        添加：
          MYSQL_ROOT_PASSWORD=xxxx
          MYSQL_PASSWORD=xxxx
          STMP_HOST=xxx
          STMP_PORT=xxx
          STMP_USER=xxx
          STMP_PASSWORD=xxx

      启动
      docker compose up --build

      前端服务: http://localhost:5173
      后端服务：http://localhost:8080

      minio服务：http://localhost:9001
          user:admin pass:12345678
            Bucketname      = "testbucket"
            EndPoint        = "minio:9000"
            AccessKeyID     = "xxx"
            SecretAccessKey = "xxx"
            UseSSL          = false

      停止
      docker compose down

## 容器

- [docker](https://www.docker.com/)

## 编排

- docker compose

## 后端技术

- 语言 go 1.25.0
- web框架 fiber v3
- orm框架 gorm
- 软删除 gorm.io/plugin/soft_delete
- 数据库 mysql8.0.40
- 驱动 "gorm.io/driver/mysql"
- 密码加密 "golang.org/x/crypto/bcrypt"
- 请求参数验证器："github.com/go-playground/validator/v10"
- JWT:"github.com/golang-jwt/jwt/v5"
- 邮件发送："gopkg.in/gomail.v2"
- 终端颜色: "github.com/fatih/color"
- 存储：github.com/minio/minio-go/v7

## 前端技术

- 框架:"vue": "^3.5.40",
- 网络请求："axios": "^1.19.0",
- UI框架："fall-ui": "^0.2.4",
- 状态管理："pinia": "^4.0.2",
- 文件编码："spark-md5": "^3.0.2",
- 路由:"vue-router": "^5.2.0"
