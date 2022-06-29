SELECT datname, oid 
FROM pg_database 
WHERE datname IN ($1)