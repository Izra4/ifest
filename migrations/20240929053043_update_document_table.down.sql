ALTER TABLE documents
ALTER COLUMN number TYPE int USING number::integer;