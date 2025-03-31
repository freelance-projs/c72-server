CREATE TABLE IF NOT EXISTS laundry (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(128) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS laundry_tag (
  laundry_id INT NOT NULL,
  tag_id VARCHAR(128) NOT NULL,
  status ENUM('washing', 'returned') NOT NULL DEFAULT 'washing',
  PRIMARY KEY (laundry_id, tag_id)
);
