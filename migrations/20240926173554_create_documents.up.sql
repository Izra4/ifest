CREATE TABLE documents (
    id uuid NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL,
    user_id char(36) NOT NULL,
    type varchar(255) NOT NULL,
    status status_enum NOT NULL,
    number int NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);