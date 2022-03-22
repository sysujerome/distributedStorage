package main

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"reflect"
)

func main() {
	db, err := leveldb.OpenFile("./foo.db", nil)
	fmt.Println(reflect.TypeOf(db))
	fmt.Println(reflect.TypeOf(err))
	if err != nil {
		log.Fatal("Yikes!")
	}
	defer db.Close()

	err = db.Put([]byte("fizz"), []byte("buzz"), nil)
	err = db.Put([]byte("fizz2"), []byte("buzz2"), nil)
	err = db.Put([]byte("fizz3"), []byte("buzz3"), nil)

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Printf("key: %s | value: %s\n", key, value)
	}

	fmt.Println("\n")

	for ok := iter.Seek([]byte("fizz2")); ok; ok = iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Printf("key: %s | value: %s\n", key, value)
	}

	fmt.Println("\n")

	for ok := iter.First(); ok; ok = iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Printf("key: %s | value: %s\n", key, value)
	}

	iter.Release()
	err = iter.Error()


	data, err := db.Get([]byte("fizz2"), nil)
	fmt.Println(string(data))
	fmt.Println(err)
	data, err = db.Get([]byte("dasd"), nil)
	if err == nil {
		return
	}
	fmt.Println(string(data))
	fmt.Println(err)
}