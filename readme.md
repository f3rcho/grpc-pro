# pkg
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

```bash
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Compiler
```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/student.proto
```

## On mac
```bash
 brew install protobuf
 protoc --version
```

## postgress
```bash
go get github.com/lib/pq
```
## Grpc
```bash
go get google.golang.org/grpc
```

## Docker
```bash
	docker build . -t db-grpc
  docker run -p 54321:5432 db-grpc
```
docker-compose ps
docker-compose logs -f postgres
docker-compose exec postgres bash
Instead of using the docker container ip, just use the name of the service when connection to pgadmin. Host name/addres=postgress

docker ps
docker inspect


# Metodología:

1. Crear el archivo .proto primero.
2. Compilarlo para generar los paquetes de Go.
3. Implementar el protobuffer a nivel de servidor, a nivel de base de datos y a nivel de interacción con gRPC.