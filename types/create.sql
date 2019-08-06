/* 
    run with command :"mysql -u root dkms < create.sql" 
*/
create table if not exists `admin`(
	pid varchar(128) primary key,
	name varchar(128)
);

create table if not exists `user`(
    pid varchar(256),
    username varchar(256),
    password varchar(256),
    phone varchar(20),
    authKey varchar(256),
    expiredTimeStamp int,
    expiredTimeStamp3rd int
);

create table if not exists `expired`(
    pid varchar(256),
    phone varchar(20),
    last  int,
    last3rd int
);

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