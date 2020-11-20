# docker-backup-database

[![GoDoc](https://godoc.org/github.com/appleboy/docker-backup-database?status.svg)](https://godoc.org/github.com/appleboy/docker-backup-database)
[![Build Status](https://cloud.drone.io/api/badges/appleboy/docker-backup-database/status.svg)](https://cloud.drone.io/appleboy/docker-backup-database)
[![codecov](https://codecov.io/gh/appleboy/docker-backup-database/branch/master/graph/badge.svg)](https://codecov.io/gh/appleboy/docker-backup-database)
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/docker-backup-database)](https://goreportcard.com/report/github.com/appleboy/docker-backup-database)
[![Docker Pulls](https://img.shields.io/docker/pulls/appleboy/docker-backup-database.svg)](https://hub.docker.com/r/appleboy/docker-backup-database/)

Docker image to periodically backup a your database (MySQL or Postgres) to Local Disk or S3 ([AWS S3](https://aws.amazon.com/free/storage/s3) or [Minio](https://min.io/)).

## Importance

This project still **In Progress** now.

## Support Database

* Postgres
* MySQL

## Usage

First steps: Setup the Minio and Postgres 12 Server using docker-compose command.
