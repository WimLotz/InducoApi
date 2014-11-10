package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
)

func MarshalObjectToJson(obj interface{}) []byte {
	bytes, err := json.Marshal(obj)
	if err != nil {
		log.Printf("Json Marshalling error: %v\n", err)
		return nil
	}
	return bytes
}

func UnmarshalJsonToObject(data []byte, v interface{}) {
	err := json.Unmarshal(data, &v)
	if err != nil {
		log.Printf("Json Unmarshalling error: %v\n", err)
	}
}

func RandomString(length int) (str string) {
	b := make([]byte, length)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
