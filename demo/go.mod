module demo

go 1.17

// replace example.com/util => ./util/

// replace example.com/util/crc16 => ./util/crc16/crc16

replace example.com/kvstore => ./kvstore

require (
	example.com/kvstore v0.0.0-00010101000000-000000000000
	github.com/go-redis/redis/v8 v8.11.4
	google.golang.org/grpc v1.44.0
)

require (
	example.com/util/crc16 v1.2.5 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20210428140749-89ef3d95e781 // indirect
	golang.org/x/sys v0.0.0-20210423082822-04245dca01da // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)

replace (
	example.com/util/crc16 v1.2.3 => ./util/crc16/crc16
	example.com/util/crc16 v1.2.5 => ./util/crc16
)
