# pks
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