CREATE TABLE `users` (
    id CHAR(36) NOT NULL,
    name varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    password varchar(255) NOT NULL,
    PRIMARY KEY (id)
)ENGINE = InnoDB;