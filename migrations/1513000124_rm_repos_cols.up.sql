ALTER TABLE repo
  DROP COLUMN vcs,
  DROP COLUMN http_clone_url,
  DROP COLUMN ssh_clone_url,
  DROP COLUMN homepage_url,
  DROP COLUMN default_branch,
  DROP COLUMN deprecated,
  DROP COLUMN mirror,
  DROP COLUMN vcs_synced_at;
