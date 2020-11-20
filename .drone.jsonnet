local pipeline = import 'pipeline.libsonnet';
local name = 'docker-backup-database';

[
  pipeline.test,
]
