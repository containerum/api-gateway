ALTER TABLE "listeners"
ADD COLUMN roles text[] DEFAULT '{}' NOT NULL;
