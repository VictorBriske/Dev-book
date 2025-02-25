CREATE DATABASE IF NOT EXISTS devbook;
USE devbook/

DROP TABLE IF EXISTS users;
CREATE TABLE users(
    id int auto_increment primary key,
    nome varchar(50) not null,
    nick varchar(50) unique,
    email varchar(50) not null unique,
    senha varchar(100) not null unique,
    criadoem timestamp default  current_timestamp()
) ENGINE=INNOOB;