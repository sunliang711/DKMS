/* 
    run with command :"mysql -u root dkms < create.sql" 
*/
/*admin表里数据是手动添加的*/
create table if not exists `admin`(
	pid varchar(128) primary key,
    phone varchar(20),
	name varchar(128),
    token varchar(256),
    t int,
    n int
);
insert into `admin` values('001','13800000000','ys','yingshi001',3,5);
insert into `admin` values('002','13900000000','bx','baoxian001',3,5);

create table if not exists `user`(
    pid varchar(256),
    username varchar(256),
    password varchar(256),
    phone varchar(20),
    authKey varchar(256),
    beginTimestamp int,
    expiredTimestamp int,
    beginTimestamp3rd int,
    expiredTimestamp3rd int
);

/*
create table if not exists `expired`(
    pid varchar(256),
    phone varchar(20),
    last  int,
    last3rd int
);
*/

create table if not exists `key`(
    pid varchar(256),
    phone varchar(20),
    pk varchar(512),
    sk varchar(512),
    keyType int
);

create table if not exists `keyfrag`(
    pid varchar(256),
    phone varchar(20),
    receiver varchar(256),
    t int,
    n int,
    segment varchar(2048)
);