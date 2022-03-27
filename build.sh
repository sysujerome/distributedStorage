go build -o ./bin/client ./client/main.go
go build -o ./bin/storeInside ./server/storage_inside.go
go build -o ./bin/storeLevelDB ./server/storage_LevelDB.go
