package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

//GenerateRandom 生成随机的字符串
func GenerateRandom(size int) string {
	kind := 3
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return string(result)
}

// MakeTimestamp 获取当前时间戳，毫秒
func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

//StructToJSONBytes 把数据转成json的bytes
func StructToJSONBytes(data interface{}) ([]byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return b, nil
}

func LogMessage(message string) {
	fmt.Println(message)
}

func LogStruct(data interface{}) {
	fmt.Printf("%+v\n", data)
}
