CREATE TABLE IF NOT EXISTS tag (
  id VARCHAR(128) PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_name (name),
  KEY idx_created_at (created_at)
);

CREATE TABLE IF NOT EXISTS department (name VARCHAR(255) PRIMARY KEY);

CREATE TABLE IF NOT EXISTS `tag_name` (`name` VARCHAR(255) NOT NULL PRIMARY KEY);

CREATE TABLE IF NOT EXISTS setting (
  id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  tx_log_sheet_id VARCHAR(128) NOT NULL,
  report_sheet_id VARCHAR(128) NOT NULL
);

CREATE TABLE IF NOT EXISTS company (name VARCHAR(255) PRIMARY KEY);

create table tx_tag (
  tag_id varchar(128) not null primary key,
  tx_id int unsigned not null,
  status enum('lending', 'washing')
);

create table if not exists tx_log_department (
  id int unsigned,
  department varchar(128) not null,
  action enum('lending', 'returned') not null,
  actor varchar(25) not null default 'admin',
  tag_name varchar(128) not null,
  lending int unsigned,
  returned int unsigned,
  created_at datetime not null default current_timestamp,
  key idx_department (department, created_at),
  key idx_tag_name (tag_name, created_at),
  key idx_created_at (created_at),
  key idx_id (id, tag_name),
  key idx_id_created (id, created_at),
  key idx_dt (department, tag_name),
  key idx_cdt (created_at, department, tag_name)
);

-- lending_stat
create table if not exists lending_stat (
  id int unsigned,
  tag_name varchar(128) not null,
  department varchar(128) not null,
  lending int unsigned,
  returned int unsigned,
  created_at datetime not null default current_timestamp,
  primary key (id, tag_name),
  key idx_department (department, tag_name),
  key idx_tag_name (tag_name)
);

create table if not exists tx_log_company (
  id int unsigned,
  company varchar(128) not null,
  action enum('washing', 'returned') not null,
  actor varchar(25) not null default 'admin',
  tag_name varchar(128) not null,
  washing int unsigned,
  returned int unsigned,
  created_at datetime not null default current_timestamp,
  key idx_company (company, created_at),
  key idx_tag_name (tag_name, created_at)
);

create table if not exists washing_stat (
  id int unsigned,
  tag_name varchar(128) not null,
  company varchar(128) not null,
  washing int unsigned,
  returned int unsigned,
  created_at datetime not null default current_timestamp,
  primary key (id, tag_name),
  key idx_company (company),
  key idx_tag_name (tag_name)
);
