create table if not exists tag_company (
  id varchar(128) not null primary key,
  name varchar(128) null
);

create table if not exists tag_department (
  id varchar(128) not null primary key,
  name varchar(128) null
);
