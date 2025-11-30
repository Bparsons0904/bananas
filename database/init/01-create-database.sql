-- Create user if it doesn't exist (ignore error if already exists)
DO $$
BEGIN
   CREATE USER bananas_user WITH PASSWORD 'bananas_pass';
EXCEPTION
   WHEN duplicate_object THEN 
      RAISE NOTICE 'User bananas_user already exists';
END
$$;

-- Create database if it doesn't exist (ignore error if already exists)
DO $$
BEGIN
   CREATE DATABASE bananas_dev OWNER bananas_user;
EXCEPTION
   WHEN duplicate_database THEN 
      RAISE NOTICE 'Database bananas_dev already exists';
END
$$;

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE bananas_dev TO bananas_user;