# .PHONY: createdb dropdb

# createdb:
# 	PGPASSWORD=#Kyabawasir69 createdb --username=postgres --owner=postgres laxmi_bank

# dropdb:
# 	PGPASSWORD=#Kyabawasir69 dropdb --username=postgres simple_bank


.PHONY: postgres createdb dropdb migrateup migratedown sqlc

postgres:
	docker run --name postgres2 -p 5433:5432 -e POSTGRES_PASSWORD=#Kyabawasir69 -d postgres:12-alpine

createdb:
	docker exec -it postgres2 createdb --username=postgres --owner=postgres laxmi_bank

dropdb:
	docker exec -it postgres2 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://postgres:#Kyabawasir69@127.0.0.1:5432/laxmi_bank" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:#Kyabawasir69@127.0.0.1:5432/laxmi_bank" -verbose down

sqlc:
	sqlc generate

