CREATE TABLE plugins (
  name TEXT PRIMARY KEY NOT NULL DEFAULT '',
  source TEXT NOT NULL DEFAULT '',
  uri TEXT NOT NULL DEFAULT '',
  install_command TEXT NOT NULL DEFAULT '',
  update_command TEXT NOT NULL DEFAULT '',
  command TEXT NOT NULL
);

CREATE TABLE accounts (
  id SERIAL PRIMARY KEY,
  plugin TEXT NOT NULL REFERENCES plugins (name),
  name TEXT NOT NULL DEFAULT '',
  options TEXT NOT NULL DEFAULT ''
);

CREATE TABLE data (
  id SERIAL PRIMARY KEY,
  remote_id TEXT NOT NULL DEFAULT '',
  plugin TEXT NOT NULL REFERENCES plugins (name),
  resource_name TEXT NOT NULL DEFAULT '',
  uri TEXT NOT NULL DEFAULT '',
  metadata TEXT NOT NULL DEFAULT '{}'
);
