up:
	docker compose up -d

down:
	docker compose down --remove-orphans --volumes

run_neo4j:
	docker run \
	--name neo4j \
    --detach \
    --publish=7474:7474 --publish=7687:7687 \
    --env NEO4J_dbms_memory_pagecache_size=4G \
    neo4j:5.22.0

rm_neo4j:
	docker rm -f $(shell docker ps -a -q --filter ancestor=neo4j:5.22.0)

cli:
	docker exec -it neo4j cypher-shell

ui:
	open http://localhost:7474