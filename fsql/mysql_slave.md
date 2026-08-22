# 启动主从

      [mysqld]
      log-bin=mysql-bin
      binlog-format=ROW
      server-id=1 # 主库ID

# 配置从

      CREATE USER 'gossip*boy'@'%' IDENTIFIED BY 'secret123';
      GRANT REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO 'gossip_boy'@'%';
      FLUSH PRIVILEGES;

# 看用户密码 协议

select user, host, plugin from mysql.user where user = 'fsql'
