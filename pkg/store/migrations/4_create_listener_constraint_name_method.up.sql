ALTER TABLE "listeners" DROP CONSTRAINT IF EXISTS unique_name_method;
ALTER TABLE "listeners" ADD CONSTRAINT unique_name_method UNIQUE ("name", "method");
