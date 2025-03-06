# docker-backup-database

[![GoDoc](https://godoc.org/github.com/appleboy/docker-backup-database?status.svg)](https://godoc.org/github.com/appleboy/docker-backup-database)
[![codecov](https://codecov.io/gh/appleboy/docker-backup-database/branch/master/graph/badge.svg)](https://codecov.io/gh/appleboy/docker-backup-database)
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/docker-backup-database)](https://goreportcard.com/report/github.com/appleboy/docker-backup-database)
[![Docker Image](https://github.com/appleboy/docker-backup-database/actions/workflows/docker.yml/badge.svg)](https://github.com/appleboy/docker-backup-database/actions/workflows/docker.yml)

[English](README.md) | 繁體中文 | [简体中文](README_zh-CN.md)

Docker 映像檔，用於定期備份您的資料庫（MySQL、Postgres 或 MongoDB）到本地磁碟或 S3（[AWS S3](https://aws.amazon.com/free/storage/s3) 或 [Minio](https://min.io/)）。

[中文 Youtube 影片](https://www.youtube.com/watch?v=nsiKKSy5fUA)

## 支援的資料庫

請參閱 [docker hub 頁面](https://hub.docker.com/repository/docker/appleboy/docker-backup-database)。

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

## Docker 映像檔

您可以從 Docker Hub Registry 拉取此專案的最新映像檔。

```sh
docker pull appleboy/docker-backup-database:postgres12
```

或者，您可以從 GitHub Container Registry 拉取此專案的最新映像檔。

```sh
docker pull ghcr.io/appleboy/docker-backup-database:postgres12
```

## 使用方法

第一步：使用 docker-compose 指令設置 Minio 和 Postgres 12 伺服器。

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

第二步：備份您的資料庫並將轉儲檔案上傳到 S3 儲存。

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

預設的生命週期策略是禁用的。您可以通過設置 `STORAGE_DAYS` 環境變數來啟用它。您可以更改 `STORAGE_DAYS` 環境變數來保留備份檔案的天數。您也可以更改 `STORAGE_PATH` 環境變數來將備份檔案保存到不同的目錄。

```yaml
STORAGE_DAYS: 30
STORAGE_PATH: backup_postgres
```

使用 Cron 排程來定期備份。請參閱 `TIME_SCHEDULE` 和 `TIME_LOCATION`

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

每行 crontab 檔案代表一個工作，格式如下：

```sh
# ┌───────────── 分鐘 (0 - 59)
# │ ┌───────────── 小時 (0 - 23)
# │ │ ┌───────────── 每月的第幾天 (1 - 31)
# │ │ │ ┌───────────── 月份 (1 - 12)
# │ │ │ │ ┌───────────── 每週的第幾天 (0 - 6) (星期日到星期六;
# │ │ │ │ │                                   7 在某些系統中也是星期日)
# │ │ │ │ │
# │ │ │ │ │
# * * * * * <要執行的命令>
```

Cron 表達式表示一組時間，使用 5 個空格分隔的欄位。

| 欄位名稱     | 必填？ | 允許的值        | 允許的特殊字符 |
| ------------ | ------ | --------------- | -------------- |
| 分鐘         | 是     | 0-59            | \* / , -       |
| 小時         | 是     | 0-23            | \* / , -       |
| 每月的第幾天 | 是     | 1-31            | \* / , - ?     |
| 月份         | 是     | 1-12 或 JAN-DEC | \* / , -       |
| 每週的第幾天 | 是     | 0-6 或 SUN-SAT  | \* / , - ?     |

您可以使用幾個預定義的排程之一來代替 Cron 表達式。

```sh
| 項目                  | 描述                                | 等同於 |
| ---------------------- | ------------------------------------------ | ------------- |
| @yearly (或 @annually) | 每年運行一次，午夜，1 月 1 日        | 0 0 1 1 *     |
| @monthly               | 每月運行一次，午夜，每月的第一天 | 0 0 1 * *     |
| @weekly                | 每週運行一次，午夜，週六/週日之間  | 0 0 * * 0     |
| @daily (或 @midnight)  | 每天運行一次，午夜                   | 0 0 * * *     |
| @hourly                | 每小時運行一次，整點        | 0 * * * *     |
```

### 設置 Webhook 通知

您可以設置 Webhook 通知將備份狀態發送到 Slack 頻道。

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

## 環境變數

### 資料庫部分

- DATABASE_DRIVER - 支援 `postgres`、`mysql` 或 `mongo`。預設為 `postgres`
- DATABASE_USERNAME - 資料庫使用者名稱
- DATABASE_PASSWORD - 資料庫密碼
- DATABASE_NAME - 資料庫名稱
- DATABASE_HOST - 資料庫主機
- DATABASE_OPTS - 請參閱 `pg_dump`、`mysqldump` 或 `mongodump` 命令

### 儲存部分

- STORAGE_DRIVER - 支援 `s3` 或 `disk`。預設為 `s3`
- ACCESS_KEY_ID - Minio 或 AWS S3 ACCESS Key ID
- SECRET_ACCESS_KEY - Minio 或 AWS S3 SECRET ACCESS Key
- STORAGE_ENDPOINT - S3 端點。預設為 `s3.amazonaws.com`
- STORAGE_BUCKET - S3 桶名稱
- STORAGE_REGION - S3 區域。預設為 `ap-northeast-1`
- STORAGE_PATH - 桶中的備份資料夾路徑。預設為 `backup`，所有轉儲檔案將保存在 `bucket_name/backup` 目錄中
- STORAGE_SSL - 預設為 `false`
- STORAGE_INSECURE_SKIP_VERIFY - 預設為 `false`
- STORAGE_DAYS - 保留備份檔案的天數。預設為 `7`

### 排程部分

- TIME_SCHEDULE - 您可以使用幾個預定義的排程之一來代替 Cron 表達式。
- TIME_LOCATION - 預設情況下，所有解釋和排程都在機器的本地時區進行。您可以在構建時指定不同的時區。

## 檔案部分

- FILE_PREFIX - 檔案的前綴名稱，預設為 `storage driver` 名稱。
- FILE_SUFFIX - 檔案的後綴名稱
- FILE_FORMAT - 檔案的格式字串，預設為 `20060102150405`。

## Webhook 部分

- WEBHOOK_URL - Webhook URL
- WEBHOOK_INSECURE - 預設為 `false`
