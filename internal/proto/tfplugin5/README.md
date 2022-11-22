# tfproto5

## Regenerate protos

```bash
protoc --proto_path=. --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative tfproto5.proto
```
