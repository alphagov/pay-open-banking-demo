LOCAL_DATABASE_URL=postgres://postgres:@localhost:5432/?sslmode=disable

start_postgres_docker:
	docker run -d -p 5432:5432 --name openbanking -e POSTGRES_PASSWORD= postgres:9.5

stop_postgres_docker:
	docker stop openbanking
	docker rm openbanking

start_local: DATABASE_URL=$(LOCAL_DATABASE_URL) ./stephens-openbanking-test