go build -o ./deploy/bin/client ./client/main.go
go build -o ./deploy/bin/storageInside ./server/storage_inside.go ./server/db.go
# go build -o ./bin/storageLevelDB ./server/storage_LevelDB.go 
