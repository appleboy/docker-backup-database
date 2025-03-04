# docker-backup-database

[![GoDoc](https://godoc.org/github.com/appleboy/docker-backup-database?status.svg)](https://godoc.org/github.com/appleboy/docker-backup-database)
[![codecov](https://codecov.io/gh/appleboy/docker-backup-database/branch/master/graph/badge.svg)](https://codecov.io/gh/appleboy/docker-backup-database)
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/docker-backup-database)](https://goreportcard.com/report/github.com/appleboy/docker-backup-database)
[![Docker Image](https://github.com/appleboy/docker-backup-database/actions/workflows/docker.yml/badge.svg)](https://github.com/appleboy/docker-backup-database/actions/workflows/docker.yml)

Docker image to periodically backup a your database (MySQL, Postgres or MongoDB) to Local Disk or S3 ([AWS S3](https://aws.amazon.com/free/storage/s3) or [Minio](https://min.io/)).

[中文 Youtube 影片](https://www.youtube.com/watch?v=nsiKKSy5fUA)

## Support Database

see the [docker hub page](https://hub.docker.com/repository/docker/appleboy/docker-backup-database).

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

## Docker Image

You can pull the latest image of the project from the Docker Hub Registry.

```sh
docker pull appleboy/docker-backup-database:postgres12
```

Or you can pull the latest image of the project from the GitHub Container Registry.

```sh
docker pull ghcr.io/appleboy/docker-backup-database:postgres12
```

## Usage

First steps: Setup the Minio and Postgres 12 Server using docker-compose command.

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

Second Steps: Backup your database and upload the dump file to S3 storage.

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

The default lifecycle policy is disabled. You can enable it by setting the `STORAGE_DAYS` environment variable. You can change the `STORAGE_DAYS` environment variable to keep the backup files for a different number of days. You also can change the `STORAGE_PATH` environment variable to save the backup files in a different directory.

```yaml
STORAGE_DAYS: 30
STORAGE_PATH: backup_postgres
```

Cron schedule to run periodic backups. See the `TIME_SCHEDULE` and `TIME_LOCATION`

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

Each line of a crontab file represents a job, and looks like this:

```sh
# ┌───────────── minute (0 - 59)
# │ ┌───────────── hour (0 - 23)
# │ │ ┌───────────── day of the month (1 - 31)
# │ │ │ ┌───────────── month (1 - 12)
# │ │ │ │ ┌───────────── day of the week (0 - 6) (Sunday to Saturday;
# │ │ │ │ │                                   7 is also Sunday on some systems)
# │ │ │ │ │
# │ │ │ │ │
# * * * * * <command to execute>
```

A cron expression represents a set of times, using 5 space-separated fields.

| Field name   | Mandatory? | Allowed values  | Allowed special characters |
| ------------ | ---------- | --------------- | -------------------------- |
| Minutes      | Yes        | 0-59            | \* / , -                   |
| Hours        | Yes        | 0-23            | \* / , -                   |
| Day of month | Yes        | 1-31            | \* / , - ?                 |
| Month        | Yes        | 1-12 or JAN-DEC | \* / , -                   |
| Day of week  | Yes        | 0-6 or SUN-SAT  | \* / , - ?                 |

You may use one of several pre-defined schedules in place of a cron expression.

```sh
| Entry                  | Description                                | Equivalent To |
| ---------------------- | ------------------------------------------ | ------------- |
| @yearly (or @annually) | Run once a year, midnight, Jan. 1st        | 0 0 1 1 *     |
| @monthly               | Run once a month, midnight, first of month | 0 0 1 * *     |
| @weekly                | Run once a week, midnight between Sat/Sun  | 0 0 * * 0     |
| @daily (or @midnight)  | Run once a day, midnight                   | 0 0 * * *     |
| @hourly                | Run once an hour, beginning of hour        | 0 * * * *     |
```

### Setup Webhook Notification

You can setup the webhook notification to send the backup status to the slack channel.

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

## Envionment Variables

### Database Section

- DATABASE_DRIVER - support `postgres`, `mysql` or `mongo`. default is `postgres`
- DATABASE_USERNAME - database username
- DATABASE_PASSWORD - database password
- DATABASE_NAME - database name
- DATABASE_HOST - database host
- DATABASE_OPTS - see the `pg_dump`, `mysqldump` or `mongodump` command

### Storage Section

- STORAGE_DRIVER - support `s3` or `disk`. default is `s3`
- ACCESS_KEY_ID - Minio or AWS S3 ACCESS Key ID
- SECRET_ACCESS_KEY - Minio or AWS S3 SECRET ACCESS Key
- STORAGE_ENDPOINT - S3 Endpoint. default is `s3.amazonaws.com`
- STORAGE_BUCKET - S3 bucket name
- STORAGE_REGION - S3 Region. default is `ap-northeast-1`
- STORAGE_PATH - backup folder path in bucket. default is `backup` and all dump file will save in `bucket_name/backup` directory
- STORAGE_SSL - default is `false`
- STORAGE_INSECURE_SKIP_VERIFY - default is `false`
- STORAGE_DAYS - The number of days to keep the backup files. default is `7`

### Schedule Section

- TIME_SCHEDULE - You may use one of several pre-defined schedules in place of a cron expression.
- TIME_LOCATION - By default, all interpretation and scheduling is done in the machine's local time zone. You can specify a different time zone on construction.

## File Section

- FILE_PREFIX - Prefix name of file, default is `storage driver` name.
- FILE_SUFFIX - Suffix name of file
- FILE_FORMAT - Format string of file, default is `20060102150405`.

## Webhook Section

- WEBHOOK_URL - Webhook URL
- WEBHOOK_INSECURE - default is `false`
