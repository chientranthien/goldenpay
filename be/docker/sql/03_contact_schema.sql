create table if not exists contact_tab (
    id bigint unsigned auto_increment,
    user_id bigint unsigned not null,
    contact_user_id bigint unsigned not null,
    status int not null,
    version int not null,
    ctime bigint not null,
    mtime bigint not null,
    primary key(id),
    unique key unid_idx_user_id_contact_user_id (user_id, contact_user_id)
);