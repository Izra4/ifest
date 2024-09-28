CREATE TABLE user_doc_access(
    user_id CHAR(36) NOT NULL,
    doc_id uuid NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (doc_id) REFERENCES documents(id)
);