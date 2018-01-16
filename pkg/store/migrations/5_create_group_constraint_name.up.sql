ALTER TABLE "groups" DROP CONSTRAINT IF EXISTS unique_name;
ALTER TABLE "groups" ADD CONSTRAINT unique_name UNIQUE ("name");
