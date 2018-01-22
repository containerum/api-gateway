ALTER TABLE "listeners"
ALTER COLUMN group_refer TYPE text USING group_refer::text;
