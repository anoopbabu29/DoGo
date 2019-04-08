package main

import (
	"context"
	"fmt"
	//"net/url"
	"time"
	//"google.golang.org/genproto/googleapis/devtools/remoteworkers/v1test2"
	//"io"
	//"log"
	"github.com/PuerkitoBio/goquery"
	//"google.golang.org/api/customsearch/v1"
	//"github.com/gocolly/colly"
	//"reflect"
	//"io/ioutil"
	//"bytes"
	"github.com/EdmundMartin/gosearcher"

	"io/ioutil"
	"net/http"
	"os"

	"cloud.google.com/go/bigquery"
	"encoding/json"
	"google.golang.org/api/iterator"
	"regexp"
	"strconv"
	"strings"
)

type testData struct {
	orgName     string `bigquery:"orgName"`
	webName     string `bigquery:"webName"`
	wordCont    string `bigquery:"wordCont"`
	donType     string `bigquery:"donType"`
	yearsHosted int64  `bigquery:"yearsHosted"`
	city        string `bigquery:"city"`
	state       string `bigquery:"state"`
	assets      int64  `bigquery:"assets"`
}

var url2 = ""

// var url1 = ""
// var url2 = ""
// var i = 1
var urls map[string]bool
var words map[string]int
var testEntry testData

// var donTYPE = "Arts, Culture and Humanities"

func processElement(index int, element *goquery.Selection) {
	// See if the href attribute exists on the element
	href, exists := element.Attr("href")
	var fo bool
	_, found := urls[url2]
	if exists {
		_, found = urls[url2+href]
	} else {
		found = true
	}
	fo = found
	if len(href) == 0 {
		return
	}
	if exists && href[0] == '/' && !fo {

		//u, _ := url.ParseRequestURI(url2)
		urls[url2+href] = true

		resp, err1 := http.Get(url2 + href)

		if err1 != nil {
			fmt.Println(err1)
		}

		doc, err3 := goquery.NewDocumentFromReader(resp.Body)
		if err3 != nil {
			fmt.Println(err3)
		}

		doc.Find("p").Each(processWords)
		doc.Find("span").Each(processWords)
		doc.Find("h1").Each(processWords)
		doc.Find("h2").Each(processWords)
		doc.Find("h3").Each(processWords)
		doc.Find("h4").Each(processWords)
		doc.Find("h5").Each(processWords)
		doc.Find("h6").Each(processWords)
		doc.Find("li").Each(processWords)
		doc.Find("b").Each(processWords)
		doc.Find("i").Each(processWords)
		doc.Find("em").Each(processWords)
		doc.Find("u").Each(processWords)
		doc.Find("a").Each(processWords)
		doc.Find("del").Each(processWords)
		doc.Find("td").Each(processWords)
		doc.Find("thead").Each(processWords)

		doc.Find("a").Each(processElement)

	}
}

func processWords(index int, element *goquery.Selection) {
	// See if the href attribute exists on the element
	remove := element.Children().Text()
	text := element.Text()
	result := strings.Replace(text, remove, "", -1)
	reg, _ := regexp.Compile("[^a-zA-Z]+")

	words2 := strings.Fields(result)
	for _, word := range words2 {
		word = strings.ToLower(word)
		word := reg.ReplaceAllString(word, "")
		_, found := words[word]
		if word != "" {
			if !found {
				words[word] = 1
			} else {
				words[word] = words[word] + 1
			}
		}
	}
}
func main() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "../Special/DoGo-80cbcac25e42.json")
	os.Setenv("GOOGLE_CLOUD_PROJECT", "dogo-236814")

	proj := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if proj == "" {
		fmt.Println("GOOGLE_CLOUD_PROJECT environment variable must be set.")
		os.Exit(1)
	}

	urls = make(map[string]bool)

	words = make(map[string]int)

	data, err := ioutil.ReadFile("../data-download-pub78.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, proj)

	fmt.Println(client)
	if err != nil {
		fmt.Println("client")
	}
	//ins := client.DatasetInProject("dogo-236814", "DoGoOrgInfo").Table("webInfo").Uploader()

	dataStr := string(data)

	temp := strings.Split(dataStr, "\n")

	for l, line := range temp {
		if l < 900 {
			continue
		}

		linCont := strings.Split(line, "|")
		i := 0
		for _, temp2 := range linCont {
			if i == 1 {
				testEntry.orgName = temp2
			} else if i == 2 {
				testEntry.city = temp2
			} else if i == 3 {
				testEntry.state = temp2
			}
			i++

		}

		searchStr := testEntry.orgName
		reg, _ := regexp.Compile("[^a-zA-Z0-9+ ]+")
		searchStr = reg.ReplaceAllString(searchStr, "")
		searchStr = strings.Replace(searchStr, " ", "+", -1)

		fmt.Println(searchStr)
		var s string

		resp2, err8 := gosearcher.BingScrape(testEntry.orgName, "com", nil, 1, 1, 1)
		fmt.Println(resp2)
		time.Sleep(3 * time.Second)
		if err8 != nil {
			continue
		}
		for _, res := range resp2 {
			s = res.ResultURL
			break
		}

		if strings.Contains(s, "facebook") {
			continue
		}

		testEntry.webName = s
		url2 = s

		urlSim := url2
		for k := range urls {
			delete(urls, k)
		}
		for k := range words {
			delete(words, k)
		}

		urls[urlSim] = true
		resp, err5 := http.Get(s)
		if err5 != nil {
			continue
		}
		doc, err3 := goquery.NewDocumentFromReader(resp.Body)
		if err3 != nil {
			fmt.Println(err3)
		}

		doc.Find("p").Each(processWords)
		doc.Find("span").Each(processWords)
		doc.Find("h1").Each(processWords)
		doc.Find("h2").Each(processWords)
		doc.Find("h3").Each(processWords)
		doc.Find("h4").Each(processWords)
		doc.Find("h5").Each(processWords)
		doc.Find("h6").Each(processWords)
		doc.Find("li").Each(processWords)
		doc.Find("b").Each(processWords)
		doc.Find("i").Each(processWords)
		doc.Find("em").Each(processWords)
		doc.Find("u").Each(processWords)
		doc.Find("a").Each(processWords)
		doc.Find("del").Each(processWords)
		doc.Find("td").Each(processWords)
		doc.Find("thead").Each(processWords)
		doc.Find("a").Each(processElement)

		wordsJSON, _ := json.Marshal(words)

		// fmt.Println(string(wordsJSON))
		// fmt.Println(string(urlsJSON))
		testEntry.wordCont = string(wordsJSON)

		fmt.Println(testEntry)

		resp, e := http.Get("http://whois.domaintools.com/" + testEntry.webName)
		if e != nil {
			fmt.Println("h12")
			testEntry.yearsHosted = 0
			continue
		} else {
			doc, _ := goquery.NewDocumentFromReader(resp.Body)
			doc.Find(".table > tbody > tr:last-child > td:nth-child(2)").Each(func(i int, element *goquery.Selection) {
				fmt.Println("h13")
				text := element.Text()
				reg, _ := regexp.Compile("[^0-9]+")
				text = reg.ReplaceAllString(text, "")
				years, _ := strconv.Atoi(text)
				testEntry.yearsHosted = int64(years)
			})
		}

		resp, er := http.Get("http://www.nonprofitfacts.com/" + testEntry.state + "/" + searchStr)
		if er != nil {
			fmt.Println("h14")
			testEntry.assets = 0
			continue
		} else {
			doc, _ := goquery.NewDocumentFromReader(resp.Body)

			doc.Find("#generalInfo > tbody > tr:nth-child(16) > td:nth-child(2)").Each(func(i int, element *goquery.Selection) {
				fmt.Println("h15")
				text := element.Text()
				reg, _ := regexp.Compile("[^0-9]+")
				text = reg.ReplaceAllString(text, "")
				years, _ := strconv.Atoi(text)
				testEntry.assets = int64(years)
			})
		}

		fmt.Println(l)
		if l > 2000 {
			break
		}

		testEntry.donType = ""
		fmt.Println(testEntry)

		// entries := []*testData{
		// 	&testEntry,
		// }

		//ins.Put(ctx, entries)

		query2 := client.Query(`INSERT INTO DoGoOrgInfo.webInfo (orgName, webName, wordCont, yearsHosted, donType, city, state, assets)
		VALUES('` + testEntry.orgName + `', '` + testEntry.webName + `','` + testEntry.wordCont + `',` + strconv.Itoa(int(testEntry.yearsHosted)) + `,'` + testEntry.donType + `','` + testEntry.city + `','` + testEntry.state + `',` + strconv.Itoa(int(testEntry.assets)) + `);`)
		job, err2 := query2.Run(ctx)
		if err2 != nil {
			fmt.Println(err2)
		}
		stat, _ := job.Status(ctx)
		fmt.Println(stat)

	}

	// query := client.Query(`INSERT INTO  DoGoOrgInfo.webInfo (orgName, webName, wordCont, donType, city, state)
	// 						VALUES(` + testEntry.orgName + `, ` + testEntry.webName + `,` + testEntry.wordCont + `,` + testEntry.city + `,` + testEntry.state + `)`)
	// _, err2 := query.Run(ctx)
	// if err2 != nil {
	// 	fmt.Println(err2)

	// n1 := bytes.IndexByte(wordsJSON, 0)
	// n2 := bytes.IndexByte(urlsJSON, 0)
	// s1 := string(wordsJSON[:n1])
	// s2 := string(urlsJSON[:n2])
	// fmt.Println(s1)
	// fmt.Println(s2)
	// query := client.Query(`INSERT INTO  DoGoOrgInfo.webInfo (orgName, webName, wordCont, donType, yearsHosted, city, state, assets) VALUES('test', 'test', '{\'hi\':2}', 'hi', 2, 'Tampa', 'FL', 0)`)

	// _, err2 := query.Run(ctx)

	query2 := client.Query("Select * FROM DoGoOrgInfo.webInfo")

	job, _ := query2.Run(ctx)

	it, err := job.Read(ctx)
	for {
		var row []bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return
		}
		fmt.Println(row)
	}
	fmt.Println()

}
