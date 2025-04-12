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
