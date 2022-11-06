package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	proxyList = []string{
		"20.24.43.214:8123",
		"157.245.207.186:8080",
		"68.183.185.62:80",
		"182.72.234.138:3127",
		"145.40.121.73:3128",
		"200.12.133.6:8080",
		"193.164.131.202:7890",
		"85.183.140.35:8080",
		"89.208.219.121:8080",
		"51.68.207.81:80",
		"134.209.96.9:443",
		"159.203.84.241:3128",
	}
)

func main() {

	content, err := ioutil.ReadFile("monthes/april.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	jsonMap := make(map[string][]string)

	err = json.Unmarshal([]byte(content), &jsonMap)
	if err != nil {
		panic(err)
	}

	var events []string
	for k, _ := range jsonMap {
		events = append(events, k)
	}

	//ch := make(chan string)

	for _, key := range events {
		for _, uri := range jsonMap[key] {
			splitUrl := strings.Split(uri, "/")
			downloadImg(uri, splitUrl[len(splitUrl)-1])
		}
	}

	//for {
	//	select {
	//	case message := <-ch:
	//		fmt.Println(message)
	//	}
	//}
}

func downloadImg(img_url string, fileName string) {

	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	randomProxy := proxyList[rand.Intn(len(proxyList))]

	var proxyPath = "http://" + randomProxy

	//creating the proxyURL
	proxyURL, err := url.Parse(proxyPath)

	if err != nil {
		log.Println(err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	client := http.Client{
		Transport: transport,
		//Timeout:   time.Second * 10,
	}

	req, err := http.NewRequest("GET", img_url, nil)

	if err != nil {
		log.Fatal("Не сформировался http запрос, ошибка - ", err)
	}

	//req.Header = http.Header{
	//	"user-agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36"},
	//}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Не выполнился http запрос, ошибка - ", err, proxyPath)
	}

	if err != nil {
		fmt.Println("A error occurred!")
		return
	}

	defer res.Body.Close()

	//body, err := io.ReadAll(res.Body)

	filePath := fmt.Sprintf("images/%s", fileName)

	output, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error while creating", filePath, "-", err)
		return
	}
	defer output.Close()

	n, err := io.Copy(output, res.Body)
	if err != nil {
		fmt.Println("Error while downloading", img_url, "-", err)
		return
	}

	fmt.Sprintf("%d bytes downloaded from %s", n, img_url)
}
