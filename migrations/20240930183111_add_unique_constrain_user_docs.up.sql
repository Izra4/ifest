ALTER TABLE user_doc_access
    ADD CONSTRAINT unique_user_doc UNIQUE (user_id, doc_id);