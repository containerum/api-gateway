ALTER TABLE listeners
ALTER COLUMN group_refer TYPE uuid USING group_refer::uuid;
