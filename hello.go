package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	fdk "github.com/fnproject/fdk-go"
)

// Token represnts AUTH token and/or error.
type Token struct {
	Token, Error string
}

// Input represents the required headers and content.
type Input struct {
	XStorageUser, XStoragePass, Content string
}

func getToken(input *Input) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://gse00013735.storage.oraclecloud.com/auth/v1.0", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("X-Storage-User", input.XStorageUser)
	req.Header.Add("X-Storage-Pass", input.XStoragePass)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return resp.Header.Get("X-Auth-Token"), nil
}

func putObject(o string, t *Token) (*http.Response, error) {
	client := &http.Client{}
	now := strconv.FormatInt(time.Now().UnixNano(), 10)
	url := "https://gse00013735.storage.oraclecloud.com/v1/Storage-gse00013735/fn_container/log" + now + ".json"
	log.Printf("Url %s\n", url)
	req, err := http.NewRequest("PUT", url, strings.NewReader(o))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("X-Auth-Token", t.Token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp, err
}

// func getObject(name string, t *Token) {
// 	client := &http.Client{}
// 	req, err := http.NewRequest("GET", "https://gse00013735.storage.oraclecloud.com/auth/v1.0", nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	req.Header.Add("X-Storage-User", input.XStorageUser)
// 	req.Header.Add("X-Storage-Pass", input.XStoragePass)
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return resp.Header.Get("X-Auth-Token"), nil
// }

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	input := &Input{}
	json.NewDecoder(in).Decode(input)
	log.Printf("%s, %s, %s\n", input.XStorageUser, input.XStoragePass, input.Content)
	token, err := getToken(input)
	log.Printf("Token: %s\n", token)
	t := &Token{Token: token, Error: ""}
	if err != nil {
		t.Error = err.Error()
		json.NewEncoder(out).Encode(t)
		log.Fatal(err)
	}

	resp, err := putObject(input.Content, t)
	log.Printf("%s %s", resp, err)
}
