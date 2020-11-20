local pipeline = import 'pipeline.libsonnet';
local name = 'docker-backup-database';

[
  pipeline.test,
  pipeline.build(name, 'linux', 'amd64', 'mysql', '5.6'),
  pipeline.build(name, 'linux', 'amd64', 'mysql', '5.7'),
  pipeline.build(name, 'linux', 'amd64', 'mysql', '8'),
  pipeline.build(name, 'linux', 'amd64', 'postgres', '9'),
  pipeline.build(name, 'linux', 'amd64', 'postgres', '10'),
  pipeline.build(name, 'linux', 'amd64', 'postgres', '11'),
  pipeline.build(name, 'linux', 'amd64', 'postgres', '12'),
  pipeline.build(name, 'linux', 'amd64', 'postgres', '13'),
]
