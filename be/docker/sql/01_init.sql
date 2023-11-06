create schema if not exists golden_pay_db;
use golden_pay_db;

create table if not exists user_tab (
    id bigint unsigned auto_increment,
    username varchar(100) not null,
    hashed_password varchar(256) not null,
    status int not null,
    version int not null,
    ctime bigint not null,
    mtime bigint not null,
    primary key(id)
);

create table if not exists wallet_tab (
    id bigint unsigned auto_increment,
    balance bigint not null,
    currency varchar(50) not null,
    status int not null,
    version int not null,
    ctime bigint not null,
    mtime bigint not null,
    primary key(id)
);

create table if not exists transaction_tab (
    id bigint unsigned auto_increment,
    from_user bigint unsigned not null,
    to_user bigint unsigned not null,
    amount bigint not null,
    status int not null,
    version int not null,
    ctime bigint not null,
    mtime bigint not null,
    primary key(id)
);