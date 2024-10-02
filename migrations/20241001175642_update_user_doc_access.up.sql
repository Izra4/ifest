ALTER TABLE user_doc_access
    ADD COLUMN token varchar(255) NOT NULL UNIQUE,
    ADD COLUMN expired_at TIMESTAMPTZ NOT NULL,
    ADD CONSTRAINT user_doc_access_pk PRIMARY KEY (user_id, doc_id, token);