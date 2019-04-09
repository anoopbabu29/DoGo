/*
 * This file provided by Facebook is for non-commercial testing and evaluation
 * purposes only. Facebook reserves all rights not expressly granted.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * FACEBOOK BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
 * ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
 * WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	//"google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	"cloud.google.com/go/bigquery"
	"context"
	"google.golang.org/api/iterator"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"sync"
	"time"
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

type comment struct {
	ID     int64  `json:"id"`
	Author string `json:"author"`
	Text   string `json:"text"`
}

const dataFile = "./comments.json"

var commentMutex = new(sync.Mutex)

// Handle comments
func handleComments(w http.ResponseWriter, r *http.Request) {
	// Since multiple requests could come in at once, ensure we have a lock
	// around all file operations
	commentMutex.Lock()
	defer commentMutex.Unlock()

	// Stat the file, so we can find its current permissions
	fi, err := os.Stat(dataFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to stat the data file (%s): %s", dataFile, err), http.StatusInternalServerError)
		return
	}

	// Read the comments from the file.
	commentData, err := ioutil.ReadFile(dataFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read the data file (%s): %s", dataFile, err), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "POST":
		// Decode the JSON data
		var comments []comment
		if err := json.Unmarshal(commentData, &comments); err != nil {
			http.Error(w, fmt.Sprintf("Unable to Unmarshal comments from data file (%s): %s", dataFile, err), http.StatusInternalServerError)
			return
		}

		// Add a new comment to the in memory slice of comments
		comments = append(comments, comment{ID: time.Now().UnixNano() / 1000000, Author: r.FormValue("author"), Text: r.FormValue("text")})

		// Marshal the comments to indented json.
		commentData, err = json.MarshalIndent(comments, "", "    ")
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to marshal comments to json: %s", err), http.StatusInternalServerError)
			return
		}

		// Write out the comments to the file, preserving permissions
		err := ioutil.WriteFile(dataFile, commentData, fi.Mode())
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to write comments to data file (%s): %s", dataFile, err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		io.Copy(w, bytes.NewReader(commentData))

	case "GET":
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// stream the contents of the file to the response
		io.Copy(w, bytes.NewReader(commentData))

	default:
		// Don't know the method, so error
		http.Error(w, fmt.Sprintf("Unsupported method: %s", r.Method), http.StatusMethodNotAllowed)
	}
}

func getQuery(w http.ResponseWriter, r *http.Request) {
	// fi, err := os.Stat(dataFile)
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("Unable to stat the data file (%s): %s", dataFile, err), http.StatusInternalServerError)
	// 	return
	// }

	// Read the comments from the file.
	fmt.Println(r)

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		fmt.Println(r.Form)
		fmt.Println(r.FormValue("q"))

		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./Special/DoGo-80cbcac25e42.json")
		os.Setenv("GOOGLE_CLOUD_PROJECT", "dogo-236814")

		proj := os.Getenv("GOOGLE_CLOUD_PROJECT")
		if proj == "" {
			fmt.Println("GOOGLE_CLOUD_PROJECT environment variable must be set.")
			os.Exit(1)
		}
		ctx := context.Background()
		client, err9 := bigquery.NewClient(ctx, proj)
		if err9 != nil {
			fmt.Println(err9)
		}
		fmt.Println(ctx)

		// var comments string
		// if err := json.Unmarshal(commentData, &comments); err != nil {
		// 	http.Error(w, fmt.Sprintf("Unable to Unmarshal comments from data file (%s): %s", dataFile, err), http.StatusInternalServerError)
		// 	return
		// }
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		query2 := client.Query("Select * FROM DoGoOrgInfo.webInfo WHERE wordCont LIKE '%" + r.FormValue("q") + "%';")

		job, err5 := query2.Run(ctx)
		if err5 != nil {
			fmt.Println(err5)
		}

		it, err6 := job.Read(ctx)
		if err6 != nil {
			fmt.Println(err6)
		}

		count := 1
		var arrStr1 []string
		var arrStr2 []string

		var mapRet map[string][]string

		mapRet = make(map[string][]string)

		for {
			var row []bigquery.Value
			err := it.Next(&row)
			if err == iterator.Done {
				break
			}
			if err != nil {
				return
			}
			fmt.Println(row[0])
			//var org string = row[0].(string)

			arrStr1 = append(arrStr1, row[0].(string))
			arrStr2 = append(arrStr2, row[1].(string))

			count++

			if count > 10 {
				break
			}

		}
		mapRet["orgName"] = arrStr1
		mapRet["webName"] = arrStr2

		data, err3 := json.Marshal(mapRet)
		if err3 != nil {
			fmt.Println(err3)
			return
		}

		io.Copy(w, bytes.NewReader(data))

	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	buildHandler := http.FileServer(http.Dir("./dogoweb/build"))
	http.Handle("/", buildHandler)
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./dogoweb/build/static")))
	http.Handle("/static/", staticHandler)

	http.HandleFunc("/api/comments", handleComments)
	http.HandleFunc("/api/getQuery", getQuery)
	//http.Handle("/src", http.FileServer(http.Dir("./src")))
	//http.Handle("/", http.FileServer(http.Dir("./public")))
	log.Println("Server started: http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
