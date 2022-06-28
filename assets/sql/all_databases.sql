SELECT datname, oid 
FROM pg_database 
WHERE NOT datname IN ('postgres', 'template1', 'template0')