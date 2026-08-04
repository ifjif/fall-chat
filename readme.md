# fall-chat 网页聊天

## 使用

      git clone https://github.com/ifjif/fall-chat

       在当前目录下创建 .data目录 .logs目录
       创建 .env文件
        添加：
          MYSQL_ROOT_PASSWORD=xxxx
          MYSQL_PASSWORD=xxxx

      启动
      docker compose up --build
      前端服务: http://localhost:5173
      后端服务：http://localhost:8080

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
- 数据库 mysql8.0.40
- 驱动 "gorm.io/driver/mysql"
- 密码加密 "golang.org/x/crypto/bcrypt"
- 请求参数验证器："github.com/go-playground/validator/v10"

## 前端技术

- 框架 vue3
- 组件库 fall-ui
