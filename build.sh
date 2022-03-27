go build -o ./bin/client ./client/main.go
go build -o ./bin/storageInside ./server/storage_inside.go
go build -o ./bin/storageLevelDB ./server/storage_LevelDB.go
