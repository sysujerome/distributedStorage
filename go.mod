module storage

go 1.17

require (
	github.com/syndtr/goleveldb v1.0.0
	google.golang.org/grpc v1.44.0
	storage/kvstore v0.0.0-00010101000000-000000000000
	storage/util/common v0.0.0-00010101000000-000000000000
	storage/util/crc16 v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang/protobuf v1.5.0 // indirect
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd // indirect
	golang.org/x/text v0.3.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)

replace storage/kvstore => ./proto

replace storage/util/common => ./util/common

replace storage/util/crc16 => ./util/crc16
