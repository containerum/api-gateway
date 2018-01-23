ALTER TABLE ONLY listeners
ALTER COLUMN created_at SET DEFAULT current_timestamp,
ALTER COLUMN updated_at SET DEFAULT current_timestamp
