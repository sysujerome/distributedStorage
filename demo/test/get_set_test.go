package main

import (
	"rand"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestSetGet(t *testing.T) {
	conns := make([]*grpc.ClientConn, 5)
	for i := 0; i < 5; i++ {
		var err error
		port := 50050 + i
		conns[i], err = grpc.Dial("192.168.1.128:"+string(port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		check(err)
	}
	for i, conn := range conns {
		c := pb.NewStorageClient(conn)

		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
		defer cancel()

	}
}
