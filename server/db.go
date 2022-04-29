package main

import (
	crc16 "storage/util/crc16"
	"strconv"
	"sync"
)

const length int = 163

type slot_t struct {
	sync.RWMutex
	data []string
}

type DataBase struct {
	slots [length]slot_t
	len   int64
}

func hashFuncInside(key string) uint16 {
	return crc16.Checksum([]byte(key), crc16.IBMTable) % uint16(length)
}

func (db *DataBase) set(key string) {
	pos := hashFuncInside(key)
	db.slots[pos].Lock()
	db.slots[pos].data = append(db.slots[pos].data, key)
	db.len++
	db.slots[pos].Unlock()
}

func (db *DataBase) get(key string) bool {
	pos := hashFuncInside(key)
	find := func(key string, arr []string) bool {
		for i := range arr {
			if arr[i] == key {
				return true
			}
		}
		return false
	}
	db.slots[pos].RLock()
	found := find(key, db.slots[pos].data)
	db.slots[pos].RUnlock()
	return found
}
func (db *DataBase) del(key string) bool {
	index := hashFuncInside(key)
	success := false
	db.slots[index].Lock()
	for i := range db.slots[index].data {
		if db.slots[index].data[i] == key {
			thisLength := len(db.slots[index].data)
			db.slots[index].data[i] = db.slots[index].data[thisLength-1]
			db.slots[index].data = db.slots[index].data[:thisLength-1]
			db.len--
			return true
		}
	}

	db.slots[index].Unlock()
	return success
}

func (db *DataBase) scan() string {
	gap := ""
	res := ""
	for d := range db.slots {
		res += gap + strconv.Itoa(d)
		gap = "\t"
		for i := range db.slots[d].data {
			res += gap + db.slots[d].data[i]
			gap = "\t"
		}
		gap = "\n"
	}
	return res
}

// func main() {
// 	var db DataBase
// 	// set...
// 	fmt.Println("set...")
// 	db.set(fmt.Sprintf("slot_a_%d", 0))
// 	db.set(fmt.Sprintf("slot_a_%d", 1))
// 	db.set(fmt.Sprintf("slot_a_%d", 2))
// 	// get...
// 	fmt.Println("get...")
// 	for i := 0; i <= 2; i++ {
// 		key := fmt.Sprintf("slot_a_%d", i)
// 		found := db.get(key)
// 		fmt.Printf("get %s : %t\n", key, found)
// 	}

// 	// del...
// 	fmt.Println("del...")
// 	for i := 0; i <= 2; i++ {
// 		key := fmt.Sprintf("slot_a_%d", i)
// 		found := db.del(key)
// 		fmt.Printf("get %s : %t\n", key, found)
// 	}
// 	// get...
// 	fmt.Println("get...")
// 	for i := 0; i <= 2; i++ {
// 		key := fmt.Sprintf("slot_a_%d", i)
// 		found := db.get(key)
// 		fmt.Printf("get %s : %t\n", key, found)
// 	}
// }
