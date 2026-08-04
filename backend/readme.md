# Go fall-chat

## 技术

- 语言 go 1.25.0
- web框架 fiber v3
- orm框架 gorm
- 数据库 mysql8.0.40
- 驱动 "gorm.io/driver/mysql"
- 密码加密 "golang.org/x/crypto/bcrypt"
- 请求参数验证器："github.com/go-playground/validator/v10"

## 数据库

**库名：fall_chat**
**表**

- `machine_id_pool`
  id bigint auto_increment primary key
  ip varchar(45),
  status tinyint default 1 comment '1:占用, 0: 释放',
  last_heartbeat datetime
  unique key uk_id(id)

- `users`
  id int64 雪花Id
  public_id varchar(36) UUID
  name varchar(50)
  email varchar(100)
  password varchar(16)
