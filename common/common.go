package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"	
	"strings"
)

func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func JsonPrettyEncode(data interface{}, out io.Writer) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "    ")
	if err := enc.Encode(data); err != nil {
		return err
	}
	return nil
}

func JsonPrettyPrint(j any) {
	var buffer bytes.Buffer

	err := JsonPrettyEncode(j, &buffer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(buffer.String())
}

func MakeResponse(name string, route string) []byte {
	resp := make(map[string]string)
	resp["status"] = "ok"
	resp["message"] = "10-https-server/" + name
	resp["route"] = route
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	return jsonResp
}

func First[T, U any](val T, _ U) T {
	return val
}

func first[T, U any](val T, _ U) T {
	return val
}

func Exists(name string) (bool, error) {	
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func GetRootDir() string{
	current := First(os.Getwd())
	arr := strings.Split(current, "/")

	return strings.Join(arr[:len(arr)-1], "/")	
}