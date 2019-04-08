package main1

import (
	"context"
	"fmt"
	"time"
	//"google.golang.org/genproto/googleapis/devtools/remoteworkers/v1test2"
	//"io"
	//"log"
	"github.com/PuerkitoBio/goquery"
	//"github.com/gocolly/colly"
	//"reflect"
	//"io/ioutil"
	//"bytes"
	"net/http"
	"os"

	"cloud.google.com/go/bigquery"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	//"google.golang.org/api/iterator"
)

type testData struct {
	orgName     string `bigquery:"orgName"`
	webName     string `bigquery:"webName"`
	wordCont    string `bigquery:"wordcont"`
	donType     string `bigquery:"donType"`
	yearsHosted int64  `bigquery:"yearsHosted"`
	city        string `bigquery:"city"`
	state       string `bigquery:"state"`
	assets      string `bigquery:"assets"`
}

var url = "https://www.guidestar.org/nonprofit-directory/arts-culture-humanities/humanities-historical-societies/"
var url1 = ""
var url2 = ""
var i = 1
var urls map[string]bool
var words map[string]int
var testEntry testData
var donTYPE = "Arts, Culture and Humanities"

func processElement(index int, element *goquery.Selection) {
	// See if the href attribute exists on the element
	href, exists := element.Attr("href")
	_, found := urls[url2+href]
	if exists && href[0] == '/' && !found {
		fmt.Println(i, ": ", href)
		i = i + 1
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
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")

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
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "DoGo-80cbcac25e42.json")
	os.Setenv("GOOGLE_CLOUD_PROJECT", "dogo-236814")

	proj := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if proj == "" {
		fmt.Println("GOOGLE_CLOUD_PROJECT environment variable must be set.")
		os.Exit(1)
	}
	urls = make(map[string]bool)

	words = make(map[string]int)

	fmt.Println("hi4")
	for i := 1; i <= 196; i++ {
		newURL := url + strconv.Itoa(i) + ".aspx"
		fmt.Println(newURL)

		client := &http.Client{}

		req, _ := http.NewRequest("GET", newURL, nil)
		req.Header.Add("cookie", "ASP.NET_SessionId=hhndyssqmmw0v52u0i32xuxf; _vwo_uuid_v2=D2A75294C8F0DF5FC5CE467D65FC35EE9|47c2e0b7d413768f712e6783c88b2b00; _gcl_au=1.1.1715526934.1553560592; _ga=GA1.2.221230262.1553560592; D_IID=68A66195-893B-3616-8E69-D701F0176FA5; D_UID=2441ADC6-A032-3C7E-8E78-8DA7B2035261; D_ZID=1769BC59-1C98-38D1-9768-655A5AEB482D; D_ZUID=1EDC497A-F791-3545-AE81-61561106B124; D_HID=EFBD2318-2C84-35CD-A7F9-B023A816218A; D_SID=72.186.196.35:x9qDzAZKE9p6cNlFaYYcdyxFRdNFxFp1gtJ+1BEFZAM; _fbp=fb.1.1553560592021.1891645858; hubspotutk=38618180477c1583ff2a25b7ba71f806; __hssrc=1; ajs_group_id=null; _gid=GA1.2.373003171.1554563162; intercom-id-na1ksg18=b4dc7aa7-bb5a-4924-b835-77102dda13a2; __atssc=google%3B2; today=1; .gifAuth=D9D0085771E9DE807521543389F3EAEEDA2224C58CA74CF8B4D48B6324E3F74C060E407F57357B901786E1BBC0213DB09BC773BB339CE21CAA14FD3AB1B74AF2467F3C48FC0DC80FB2233DBA5B4305F5DB49453BA0D1821E0540382319069FCCBEFA42E7D755F2221BCA37AF; NOPCOMMERCE.AUTH=F75621D21F43C86924D8CE6708DE7C052807A465915DE55A9CD334C3F8860A4D9583DECDF456FF70B6E8102B5EC59946EC08DD33E1331BD2B7B3DB33B4EB79868C776687AAFF28A5255669CB97403AD4510DEA38EFB6BC417F6E148E4400AFD11B4D425D901D8C95F088986F804E9D5F0093646F80153345CA2604FB3581F8FAE4B3B047; Nop.customer=9e681a22-22d4-4596-9967-47b2bae53db1; ajs_user_id=5314007; ajs_anonymous_id=%2233504b22-c15b-4215-9c10-d5774c00e721%22; __hstc=126119634.38618180477c1583ff2a25b7ba71f806.1553560592137.1554605070786.1554611439121.8")
		req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36")
		resp, err1 := client.Do(req)

		if err1 != nil {
			fmt.Println(err1)
		}

		doc, err3 := goquery.NewDocumentFromReader(resp.Body)
		time.Sleep(5 * time.Second)
		if err3 != nil {
			fmt.Println(err3)
		}
		for j := 1; j <= 3; j++ {
			el := doc.Find("div")
			fmt.Println(el.Text())
			time.Sleep(2 * time.Second)

			doc.Find(".column-33l a").Each(func(index int, element *goquery.Selection) {
				time.Sleep(5 * time.Second)

				fmt.Println("hi2")
				//href, exists := element.Attr("href")
				if true {
					client := &http.Client{}
					req, _ := http.NewRequest("GET", "https://www.guidestar.org/profile/26-2291682", nil)
					req.Header.Add("cookie", "ASP.NET_SessionId=hhndyssqmmw0v52u0i32xuxf; _vwo_uuid_v2=D2A75294C8F0DF5FC5CE467D65FC35EE9|47c2e0b7d413768f712e6783c88b2b00; _gcl_au=1.1.1715526934.1553560592; _ga=GA1.2.221230262.1553560592; D_IID=68A66195-893B-3616-8E69-D701F0176FA5; D_UID=2441ADC6-A032-3C7E-8E78-8DA7B2035261; D_ZID=1769BC59-1C98-38D1-9768-655A5AEB482D; D_ZUID=1EDC497A-F791-3545-AE81-61561106B124; D_HID=EFBD2318-2C84-35CD-A7F9-B023A816218A; D_SID=72.186.196.35:x9qDzAZKE9p6cNlFaYYcdyxFRdNFxFp1gtJ+1BEFZAM; _fbp=fb.1.1553560592021.1891645858; hubspotutk=38618180477c1583ff2a25b7ba71f806; __hssrc=1; ajs_group_id=null; _gid=GA1.2.373003171.1554563162; intercom-id-na1ksg18=b4dc7aa7-bb5a-4924-b835-77102dda13a2; __atssc=google%3B2; today=1; .gifAuth=D9D0085771E9DE807521543389F3EAEEDA2224C58CA74CF8B4D48B6324E3F74C060E407F57357B901786E1BBC0213DB09BC773BB339CE21CAA14FD3AB1B74AF2467F3C48FC0DC80FB2233DBA5B4305F5DB49453BA0D1821E0540382319069FCCBEFA42E7D755F2221BCA37AF; NOPCOMMERCE.AUTH=F75621D21F43C86924D8CE6708DE7C052807A465915DE55A9CD334C3F8860A4D9583DECDF456FF70B6E8102B5EC59946EC08DD33E1331BD2B7B3DB33B4EB79868C776687AAFF28A5255669CB97403AD4510DEA38EFB6BC417F6E148E4400AFD11B4D425D901D8C95F088986F804E9D5F0093646F80153345CA2604FB3581F8FAE4B3B047; Nop.customer=9e681a22-22d4-4596-9967-47b2bae53db1; ajs_user_id=5314007; ajs_anonymous_id=%2233504b22-c15b-4215-9c10-d5774c00e721%22; __hstc=126119634.38618180477c1583ff2a25b7ba71f806.1554617749645.1554617749645.1554617749645.1; mp_65f8e55e260a1f3b2f557db5556fcdfc_mixpanel=%7B%22distinct_id%22%3A%205314007%2C%22%24device_id%22%3A%20%22169f6704d5533e-0ac86b66c26b13-36637902-13c680-169f6704d5698f%22%2C%22mp_lib%22%3A%20%22Segment%3A%20web%22%2C%22%24initial_referrer%22%3A%20%22https%3A%2F%2Fwww.guidestar.org%2Fnonprofit-directory%2Farts-culture-humanities%2Fhumanities-historical-societies%2F1.aspx%22%2C%22%24initial_referring_domain%22%3A%20%22www.guidestar.org%22%2C%22%24user_id%22%3A%205314007%2C%22mp_name_tag%22%3A%20%22anoop.bab%40gmail.com%22%2C%22id%22%3A%205314007%2C%22%24created%22%3A%20%222019-04-06T15%3A08%3A10.000Z%22%2C%22%24email%22%3A%20%22anoop.bab%40gmail.com%22%2C%22%24first_name%22%3A%20%22Anoop%22%2C%22%24last_name%22%3A%20%22Babu%22%2C%22%24name%22%3A%20%22Anoop%20Babu%22%7D; __atuvc=8%7C15; __atuvs=5ca99594d084dc75007; mp_5d9e4f46acaba87f5966b8c0d2e47e6e_mixpanel=%7B%22distinct_id%22%3A%20%22anoop.bab%40gmail.com%22%2C%22%24device_id%22%3A%20%22169f670a1f77c8-0aabfcbf924c37-36637902-13c680-169f670a1f88d9%22%2C%22%24initial_referrer%22%3A%20%22https%3A%2F%2Fwww.guidestar.org%2Fnonprofit-directory%2Farts-culture-humanities%2Fhumanities-historical-societies%2F1.aspx%22%2C%22%24initial_referring_domain%22%3A%20%22www.guidestar.org%22%2C%22__mps%22%3A%20%7B%7D%2C%22__mpso%22%3A%20%7B%7D%2C%22__mpus%22%3A%20%7B%7D%2C%22__mpa%22%3A%20%7B%7D%2C%22__mpu%22%3A%20%7B%7D%2C%22__mpr%22%3A%20%5B%5D%2C%22__mpap%22%3A%20%5B%5D%2C%22%24email%22%3A%20%22anoop.bab%40gmail.com%22%2C%22%24first_name%22%3A%20%22Anoop%22%2C%22%24last_name%22%3A%20%22Babu%22%2C%22Market%20Segment%22%3A%20%22Other%22%2C%22Market%20Subsegment%22%3A%20%22Other%22%2C%22GX%20NPO%20Manager%22%3A%20%22No%22%2C%22Premium%22%3A%20%22No%22%2C%22Charity%20Check%22%3A%20%22No%22%2C%22Premium%20Pro%22%3A%20%22No%22%2C%22%24user_id%22%3A%20%22anoop.bab%40gmail.com%22%7D; _gat_UA-946060-8=1; __hssc=126119634.17.1554617749647; OptanonConsent=landingPath=NotLandingPage&datestamp=Sun+Apr+07+2019+02%3A22%3A26+GMT-0400+(Eastern+Daylight+Time)&version=3.6.18&AwaitingReconsent=false&groups=1%3A1%2C2%3A1%2C3%3A1%2C4%3A1%2C101%3A1%2C102%3A1%2C103%3A1%2C104%3A1%2C105%3A1%2C106%3A1%2C107%3A1%2C108%3A1%2C109%3A1")

					req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36")
					resp, err1 := client.Do(req)
					//resp, err1 := http.Get(href)
					if err1 != nil {
						fmt.Println(err1)
						return

					} else {
						doc, err3 := goquery.NewDocumentFromReader(resp.Body)
						if err3 != nil {
							fmt.Println(err3)
						}

						testEntry.orgName = element.Text()
						testEntry.donType = donTYPE

						doc.Find(".location").Each(func(index int, element *goquery.Selection) {
							fmt.Println(element.Text())
							time.Sleep(5 * time.Second)
							result := element.Text()
							reg, _ := regexp.Compile("[^a-zA-Z0-9]+")

							words2 := strings.Fields(result)
							if len(words2) >= 2 {
								testEntry.city = reg.ReplaceAllString(words2[0], "")
								testEntry.state = words2[1]
							}
						})

						doc.Find(".website").Each(func(index int, element *goquery.Selection) {
							fmt.Println(element.Text())
							time.Sleep(5 * time.Second)
							fmt.Println("hi1")
							href, exists := element.Attr("href")
							fmt.Println(href)
							if exists {

								testEntry.webName = href
								url2 = href
								urlSim := url2 + "/"

								for k := range urls {
									delete(urls, k)
								}
								for k := range words {
									delete(urls, k)
								}

								urls[urlSim] = true
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

								ctx := context.Background()

								client, err := bigquery.NewClient(ctx, proj)

								fmt.Println(client)
								if err != nil {
									fmt.Println("client")
								}
								query := client.Query(`INSERT INTO  DoGoOrgInfo.webInfo (orgName, webName, wordCont, donType, city, state) 
														VALUES(` + testEntry.orgName + `, ` + testEntry.webName + `,` + testEntry.wordCont + `,` + testEntry.city + `,` + testEntry.state + `)`)
								_, err2 := query.Run(ctx)
								if err2 != nil {
									fmt.Println(err2)
								}
							}
						})
					}
				}

			})

		}

		// n1 := bytes.IndexByte(wordsJSON, 0)
		// n2 := bytes.IndexByte(urlsJSON, 0)
		// s1 := string(wordsJSON[:n1])
		// s2 := string(urlsJSON[:n2])
		// fmt.Println(s1)
		// fmt.Println(s2)
		// query := client.Query(`INSERT INTO  DoGoOrgInfo.webInfo (orgName, webName, wordCont, donType, yearsHosted, city, state, assets) VALUES('test', 'test', '{\'hi\':2}', 'hi', 2, 'Tampa', 'FL', 0)`)

		// _, err2 := query.Run(ctx)

		// query2 := client.Query("Select * FROM DoGoOrgInfo.webInfo")

		// job, _ := query2.Run(ctx)

		// if err2 != nil {
		// 	return
		// }
		// it, err := job.Read(ctx)
		// for {
		// 	var row []bigquery.Value
		// 	err := it.Next(&row)
		// 	if err == iterator.Done {
		// 		break
		// 	}
		// 	if err != nil {
		// 		return
		// 	}
		// 	fmt.Println(row)
		// }
		// fmt.Println()
	}
}
