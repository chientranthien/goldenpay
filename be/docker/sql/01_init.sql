create schema if not exists golden_pay_db;
use golden_pay_db;

create table if not exists user_tab (
    id bigint unsigned auto_increment,
    email varchar(256) not null,
    hashed_password varchar(256) not null,
    name nvarchar(256) not null,
    status int not null,
    version int not null,
    ctime bigint not null,
    mtime bigint not null,
    primary key(id),
    unique key(email)
);

create table if not exists wallet_tab (
    id bigint unsigned auto_increment,
    balance bigint not null,
    user_id bigint not null,
    currency varchar(50) not null,
    status int not null,
    version int not null,
    ctime bigint not null,
    mtime bigint not null,
    primary key(id),
    KEY `user_id_idx` (`user_id`)
);

create table if not exists transaction_tab (
    id bigint unsigned auto_increment,
    from_user bigint unsigned not null,
    to_user bigint unsigned not null,
    from_wallet bigint unsigned not null,
    to_wallet bigint unsigned not null,
    amount bigint not null,
    status int not null,
    version int not null,
    ctime bigint not null,
    mtime bigint not null,
    primary key(id),
    KEY `transaction_user_idx` (from_user, to_user, id DESC)
);

create table if not exists topup_tab (
   id bigint unsigned auto_increment,
   user_id bigint unsigned not null,
   wallet_id bigint unsigned not null,
   amount bigint not null,
   status int not null,
   version int not null,
   ctime bigint not null,
   mtime bigint not null,
   primary key(id)
);
