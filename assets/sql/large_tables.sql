SELECT table_name as name
FROM (
    SELECT table_name,pg_total_relation_size(table_name) AS total_size
	FROM (
	    SELECT (table_schema || '.' || table_name) AS table_name 
        FROM information_schema.tables) AS all_tables
 	    ORDER BY total_size DESC) AS pretty_sizes 
WHERE total_size > 4294967296;