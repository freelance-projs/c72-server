create table tx_tag (
  tag_id varchar(128) not null primary key,
  tx_id int unsigned not null,
  status enum('lending', 'washing')
);

-- create table tx_log_dept (
--   id int unsigned not null auto_increment primary key,
--   details json not null,
--   overview json not null,
--   created_at datetime not null default current_timestamp,
--   key idx_created_at (created_at)
-- );
--
-- create table tx_log_company (
--   id int unsigned not null auto_increment primary key,
--   details json not null,
--   overview json not null,
--   created_at datetime not null default current_timestamp,
--   key idx_created_at (created_at)
-- );
