package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
)

func ReadRequestBody(r io.Reader) []byte {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		log.Printf("Error occured reading from request body: %v\n", err)
	}
	return body
}

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
