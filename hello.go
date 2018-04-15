package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	// glog "github.com/golang/glog"
	fdk "github.com/fnproject/fdk-go"
)

// Input represents the required headers and content.
type Input struct {
	XStorageUser, XStoragePass string
	Content                    interface{}
}

func getToken(input *Input) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://gse00013735.storage.oraclecloud.com/auth/v1.0", nil)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	req.Header.Add("X-Storage-User", input.XStorageUser)
	req.Header.Add("X-Storage-Pass", input.XStoragePass)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		return "", err
	}
	return resp.Header.Get("X-Auth-Token"), err
}

func putObject(content interface{}, token string) (*http.Response, string) {
	client := &http.Client{}
	now := strconv.FormatInt(time.Now().UnixNano(), 10)
	url := "https://gse00013735.storage.oraclecloud.com/v1/Storage-gse00013735/fn_container/log" + now + ".json"
	log.Printf("Url %s\n", url)

	b, err := json.Marshal(content)
	if err != nil {
		log.Fatal(err)
	}

	r := bytes.NewReader(b)
	req, err := http.NewRequest("PUT", url, r)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("X-Auth-Token", token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp, url
}

func getObject(url string, token string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	req.Header.Add("X-Auth-Token", token)
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	} else {
		defer resp.Body.Close()
		out, err := os.Create("output.json")
		if err != nil {
			log.Fatal(err)
		} else {
			defer out.Close()
			n, err := io.Copy(out, resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Lenght: %d", n)
		}
	}

	return resp, err
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	input := &Input{}
	json.NewDecoder(in).Decode(input)
	log.Printf("%s, %s, %s", input.XStorageUser, input.XStoragePass, input.Content)
	token, err := getToken(input)
	log.Printf("Token: %s", token)
	if err != nil {
		json.NewEncoder(out).Encode(err.Error())
		log.Fatal(err)
	} else {
		resp, url := putObject(input.Content, token)
		log.Printf("Status: %d, URL: %s", resp.StatusCode, url)
		json.NewEncoder(out).Encode(url)
	}
}

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}
