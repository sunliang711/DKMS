/* 
    run with command :"mysql -u root dkms < create.sql" 
*/
/*admin表里数据是手动添加的*/
drop table if exists `admin`;
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

drop table if exists `user`;
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
insert into `user` values('001','ys','1234','13800000000','1234',0,0,0,0);
insert into `user` values('002','bx','1234','13900000000','1234',0,0,0,0);
/*
create table if not exists `expired`(
    pid varchar(256),
    phone varchar(20),
    last  int,
    last3rd int
);
*/

drop table if exists `key`;
create table if not exists `key`(
    pid varchar(256),
    phone varchar(20),
    pk varchar(512),
    sk varchar(512),
    keyType int
);
insert into `key` values ('001','13800000000','02ffbb293a6140923103359869be9e6a86ab7daf5905b60926ea8b5d0f995764a2','83a930d4f3688adf6b6a06141dd3816b646acc17c39db3c1da1bfcd3fff124c6',1);
insert into `key` values ('001','13800000000','0337e24a7bd69448664c42818880566c2cc71dc6dfe889bd1447a01fc1453c6b1a','f8d92530fdd2a66f4f26208696bb0f7050b728beee5d7075d1bce1ac959aabf6',2);
insert into `key` values ('002','13900000000','0246745e64c4f3fc6a4f8a80bad29e7c2a4748f8d0138eb5f5954921b593a51535','3baec7f126cf5e8a894f3c0a074a23ac15d6f6afc6c4e5a9b13e32ef2ca715d6',1);
insert into `key` values ('002','13900000000','03e3cf37727b00c5b80e4ff38b0151bb63007ebf6b7605037096b277f70759891b','6f6df6a475651c7715f4e13bed8dda23cbf05c3a969915085294a0f09a958aa6',2);

drop table if exists `keyfrag`;
create table if not exists `keyfrag`(
    pid varchar(256),
    phone varchar(20),
    receiver varchar(256),
    t int,
    n int,
    segment varchar(2048)
);