create database testdb;
create user 'test'@'%';
grant all privileges on testdb.* to 'test'@'%';
use testdb;
create table event(
    id int primary key auto_increment,
    dt date, money int, description varchar(32)
);
create table tag(
    id int primary key auto_increment,
    name varchar(16) unique
);
create table event_to_tag(
    event_id int, tag_id int,
    foreign key (event_id) references event(id) on delete cascade,
    foreign key (tag_id) references tag(id) on delete cascade,
    primary key(event_id, tag_id)
);
