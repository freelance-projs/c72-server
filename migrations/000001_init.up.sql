CREATE TABLE IF NOT EXISTS tags (
  id VARCHAR(128) PRIMARY KEY,
  name VARCHAR(255) null,
  is_scanned BOOLEAN NOT NULL DEFAULT FALSE,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_name (name),
  KEY idx_created_at (created_at)
);

CREATE TABLE IF NOT EXISTS tag_scan_histories (
  id int unsigned primary key auto_increment,
  tag_id VARCHAR(128) not null,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
