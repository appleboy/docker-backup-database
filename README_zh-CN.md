# docker-backup-database

[![GoDoc](https://godoc.org/github.com/appleboy/docker-backup-database?status.svg)](https://godoc.org/github.com/appleboy/docker-backup-database)
[![codecov](https://codecov.io/gh/appleboy/docker-backup-database/branch/master/graph/badge.svg)](https://codecov.io/gh/appleboy/docker-backup-database)
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/docker-backup-database)](https://goreportcard.com/report/github.com/appleboy/docker-backup-database)
[![Docker Image](https://github.com/appleboy/docker-backup-database/actions/workflows/docker.yml/badge.svg)](https://github.com/appleboy/docker-backup-database/actions/workflows/docker.yml)

[English](README.md) | [繁體中文](README_zh-TW.md) | 简体中文

Docker 镜像定期备份您的数据库（MySQL、Postgres 或 MongoDB）到本地磁盘或 S3（[AWS S3](https://aws.amazon.com/free/storage/s3) 或 [Minio](https://min.io/)）。

[中文 Youtube 影片](https://www.youtube.com/watch?v=nsiKKSy5fUA)

## 支持的数据库

请参阅 [docker hub 页面](https://hub.docker.com/repository/docker/appleboy/docker-backup-database)。

- Postgres (9, 10, 11, 12, 13, 14, 15, 16, 17)
  - 9: appleboy/docker-backup-database:postgres9
  - 10: appleboy/docker-backup-database:postgres10
  - 11: appleboy/docker-backup-database:postgres11
  - 12: appleboy/docker-backup-database:postgres12
  - 13: appleboy/docker-backup-database:postgres13
  - 14: appleboy/docker-backup-database:postgres14
  - 15: appleboy/docker-backup-database:postgres15
  - 16: appleboy/docker-backup-database:postgres16
  - 17: appleboy/docker-backup-database:postgres17
- MySQL (8, 9)
  - 8: appleboy/docker-backup-database:mysql8
  - 9: appleboy/docker-backup-database:mysql9
- Mongo (4.4)
  - 4.4: appleboy/docker-backup-database:mongo4.4

## Docker 镜像

您可以从 Docker Hub 注册表中拉取该项目的最新镜像。

```sh
docker pull appleboy/docker-backup-database:postgres12
```

或者，您可以从 GitHub 容器注册表中拉取该项目的最新镜像。

```sh
docker pull ghcr.io/appleboy/docker-backup-database:postgres12
```

## 使用方法

第一步：使用 docker-compose 命令设置 Minio 和 Postgres 12 服务器。

```yaml
services:
  minio:
    image: quay.io/minio/minio
    restart: always
    volumes:
      - data1-1:/data1
    ports:
      - 9000:9000
      - 9001:9001
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    command: server /data --console-address ":9001"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 30s
      timeout: 20s
      retries: 3

  postgres:
    image: postgres:12
    restart: always
    volumes:
      - pg-data:/var/lib/postgresql/data
    logging:
      options:
        max-size: "100k"
        max-file: "3"
    environment:
      POSTGRES_USER: db
      POSTGRES_DB: db
      POSTGRES_PASSWORD: db
```

第二步：备份您的数据库并将转储文件上传到 S3 存储。

```yaml
backup_postgres:
  image: appleboy/docker-backup-database:postgres12
  logging:
    options:
      max-size: "100k"
      max-file: "3"
  environment:
    STORAGE_DRIVER: s3
    STORAGE_ENDPOINT: minio:9000
    STORAGE_BUCKET: test
    STORAGE_REGION: ap-northeast-1
    STORAGE_PATH: backup_postgres
    STORAGE_SSL: "false"
    STORAGE_INSECURE_SKIP_VERIFY: "false"
    ACCESS_KEY_ID: minioadmin
    SECRET_ACCESS_KEY: minioadmin

    DATABASE_DRIVER: postgres
    DATABASE_HOST: postgres:5432
    DATABASE_USERNAME: db
    DATABASE_PASSWORD: db
    DATABASE_NAME: db
    DATABASE_OPTS:
```

默认的生命周期策略是禁用的。您可以通过设置 `STORAGE_DAYS` 环境变量来启用它。您可以更改 `STORAGE_DAYS` 环境变量以保持备份文件的天数。您还可以更改 `STORAGE_PATH` 环境变量以将备份文件保存在不同的目录中。

```yaml
STORAGE_DAYS: 30
STORAGE_PATH: backup_postgres
```

使用 Cron 调度定期备份。请参阅 `TIME_SCHEDULE` 和 `TIME_LOCATION`

```yaml
backup_mysql:
  image: appleboy/docker-backup-database:mysql8
  logging:
    options:
      max-size: "100k"
      max-file: "3"
  environment:
    STORAGE_DRIVER: s3
    STORAGE_ENDPOINT: minio:9000
    STORAGE_BUCKET: test
    STORAGE_REGION: ap-northeast-1
    STORAGE_PATH: backup_mysql
    STORAGE_SSL: "false"
    STORAGE_INSECURE_SKIP_VERIFY: "false"
    ACCESS_KEY_ID: 1234567890
    SECRET_ACCESS_KEY: 1234567890

    DATABASE_DRIVER: mysql
    DATABASE_HOST: mysql:3306
    DATABASE_USERNAME: root
    DATABASE_PASSWORD: db
    DATABASE_NAME: db
    DATABASE_OPTS:

    TIME_SCHEDULE: "@daily"
    TIME_LOCATION: Asia/Taipei
```

crontab 文件的每一行代表一个作业，如下所示：

```sh
# ┌───────────── 分钟 (0 - 59)
# │ ┌───────────── 小时 (0 - 23)
# │ │ ┌───────────── 月中的某天 (1 - 31)
# │ │ │ ┌───────────── 月 (1 - 12)
# │ │ │ │ ┌───────────── 星期几 (0 - 6) (星期天到星期六;
# │ │ │ │ │                                   在某些系统中，7 也是星期天)
# │ │ │ │ │
# │ │ │ │ │
# * * * * * <要执行的命令>
```

cron 表达式表示一组时间，使用 5 个空格分隔的字段。

| 字段名称   | 必需的？ | 允许的值        | 允许的特殊字符 |
| ---------- | -------- | --------------- | -------------- |
| 分钟       | 是       | 0-59            | \* / , -       |
| 小时       | 是       | 0-23            | \* / , -       |
| 月中的某天 | 是       | 1-31            | \* / , - ?     |
| 月         | 是       | 1-12 或 JAN-DEC | \* / , -       |
| 星期几     | 是       | 0-6 或 SUN-SAT  | \* / , - ?     |

您可以使用几个预定义的调度之一来代替 cron 表达式。

```sh
| 条目                  | 描述                                | 等效于 |
| ---------------------- | ------------------------------------------ | ------------- |
| @yearly (或 @annually) | 每年运行一次，午夜，1 月 1 日        | 0 0 1 1 *     |
| @monthly               | 每月运行一次，午夜，月初 | 0 0 1 * *     |
| @weekly                | 每周运行一次，星期六/星期天之间的午夜  | 0 0 * * 0     |
| @daily (或 @midnight)  | 每天运行一次，午夜                   | 0 0 * * *     |
| @hourly                | 每小时运行一次，整点        | 0 * * * *     |
```

### 设置 Webhook 通知

您可以设置 webhook 通知，将备份状态发送到 slack 频道。

```diff
backup_mysql:
  image: appleboy/docker-backup-database:mysql8
  logging:
    options:
      max-size: "100k"
      max-file: "3"
  environment:
    STORAGE_DRIVER: s3
    STORAGE_ENDPOINT: minio:9000
    STORAGE_BUCKET: test
    STORAGE_REGION: ap-northeast-1
    STORAGE_PATH: backup_mysql
    STORAGE_SSL: "false"
    STORAGE_INSECURE_SKIP_VERIFY: "false"
    ACCESS_KEY_ID: 1234567890
    SECRET_ACCESS_KEY: 1234567890

    DATABASE_DRIVER: mysql
    DATABASE_HOST: mysql:3306
    DATABASE_USERNAME: root
    DATABASE_PASSWORD: db
    DATABASE_NAME: db
    DATABASE_OPTS:

    TIME_SCHEDULE: "@daily"
    TIME_LOCATION: Asia/Taipei

+   WEBHOOK_URL: https://example.com/webhook
+   WEBHOOK_INSECURE: "false"
```

## 环境变量

### 数据库部分

- DATABASE_DRIVER - 支持 `postgres`、`mysql` 或 `mongo`。默认是 `postgres`
- DATABASE_USERNAME - 数据库用户名
- DATABASE_PASSWORD - 数据库密码
- DATABASE_NAME - 数据库名称
- DATABASE_HOST - 数据库主机
- DATABASE_OPTS - 请参阅 `pg_dump`、`mysqldump` 或 `mongodump` 命令

### 存储部分

- STORAGE_DRIVER - 支持 `s3` 或 `disk`。默认是 `s3`
- ACCESS_KEY_ID - Minio 或 AWS S3 访问密钥 ID
- SECRET_ACCESS_KEY - Minio 或 AWS S3 秘密访问密钥
- STORAGE_ENDPOINT - S3 端点。默认是 `s3.amazonaws.com`
- STORAGE_BUCKET - S3 存储桶名称
- STORAGE_REGION - S3 区域。默认是 `ap-northeast-1`
- STORAGE_PATH - 存储桶中的备份文件夹路径。默认是 `backup`，所有转储文件将保存在 `bucket_name/backup` 目录中
- STORAGE_SSL - 默认是 `false`
- STORAGE_INSECURE_SKIP_VERIFY - 默认是 `false`
- STORAGE_DAYS - 保留备份文件的天数。默认是 `7`

### 调度部分

- TIME_SCHEDULE - 您可以使用几个预定义的调度之一来代替 cron 表达式。
- TIME_LOCATION - 默认情况下，所有解释和调度都在机器的本地时区进行。您可以在构建时指定不同的时区。

## 文件部分

- FILE_PREFIX - 文件的前缀名称，默认是 `存储驱动` 名称。
- FILE_SUFFIX - 文件的后缀名称
- FILE_FORMAT - 文件的格式字符串，默认是 `20060102150405`。

## Webhook 部分

- WEBHOOK_URL - Webhook URL
- WEBHOOK_INSECURE - 默认是 `false`
