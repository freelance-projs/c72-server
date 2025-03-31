CREATE TABLE IF NOT EXISTS lending (
  id INT PRIMARY KEY AUTO_INCREMENT,
  department VARCHAR(128) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS lending_tag (
  lending_id INT NOT NULL,
  tag_id VARCHAR(128) NOT NULL,
  status ENUM('lending', 'returned') NOT NULL DEFAULT 'lending',
  PRIMARY KEY (lending_id, tag_id)
);
