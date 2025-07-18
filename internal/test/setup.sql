create database testdb;
create user 'test'@'%';
grant all privileges on testdb.* to 'test'@'%';
use testdb;
create table event(
    id bigint primary key auto_increment,
    dt date, money int, description varchar(32)
);
create table tag(
    id bigint primary key auto_increment,
    name varchar(16) unique
);
create table event_to_tag(
    event_id bigint, tag_id bigint,
    foreign key (event_id) references event(id) on delete cascade,
    foreign key (tag_id) references tag(id) on delete cascade,
    primary key(event_id, tag_id)
);

create table pattern(
    id bigint primary key auto_increment,
    key_string varchar(32) unique
);

create table pattern_to_tag(
    pattern_id bigint, tag_id bigint,
    foreign key (pattern_id) references pattern(id) on delete cascade,
    foreign key (tag_id) references tag(id) on delete cascade,
    primary key(pattern_id, tag_id)
);
