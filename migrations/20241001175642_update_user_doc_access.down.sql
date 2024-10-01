ALTER TABLE user_doc_access
    DROP CONSTRAINT user_doc_access_pk,
    DROP COLUMN token,
    DROP COLUMN expires_at;
