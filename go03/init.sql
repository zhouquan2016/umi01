create table menu
(
    id             int auto_increment
        primary key,
    name           varchar(20) not null comment '菜单名称',
    is_leaf        tinyint(1)  not null comment '是否是叶子节点',
    is_sys_default tinyint(1)  not null comment '系统预留无法删除修改',
    locale         varchar(10) not null comment '语言',
    path           varchar(50) not null comment '访问路径',
    parent_id      int         not null comment '父菜单id',
    constraint menu_path_uindex
        unique (path)
)
    comment '菜单';
INSERT INTO menu (id,name, is_leaf, is_sys_default, locale, path, parent_id, dept) VALUES (1,'系统管理', 0, 1, '', '/admin', 0, '1:');
INSERT INTO menu (id,name, is_leaf, is_sys_default, locale, path, parent_id, dept) VALUES (2,'菜单管理', 1, 1, '', '/admin/menu', '1:2');

create table menu_role
(
    id      int auto_increment
        primary key,
    role_id int not null comment '角色id',
    menu_id int not null comment '菜单id',
    constraint uindex
        unique (role_id, menu_id)
)
    comment '角色菜单';

create table user
(
    id        int auto_increment
        primary key,
    name      varchar(10)  not null comment '姓名',
    password  varchar(32)  not null comment '密码',
    avatar    varchar(200) not null comment '头像地址',
    email     varchar(20)  not null comment '邮箱',
    signature varchar(100) not null comment '个性签名',
    title     varchar(20)  not null comment '头衔',
    address   varchar(100) not null comment '地址',
    phone     varchar(11)  not null comment '手机号',
    access    varchar(100) not null,
    constraint user_email_uindex
        unique (email),
    constraint user_phone_uindex
        unique (phone)
)
    comment '用户';


