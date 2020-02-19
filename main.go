package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/psilva261/timsort"
)

const (
	commentsURL = "https://jsonplaceholder.typicode.com/comments"
	postsURL    = "https://jsonplaceholder.typicode.com/posts"
)

// Post hold posts response
type Post struct {
	UserID int    `json:"userId,omitempty"`
	ID     int    `json:"id,omitempty"`
	Title  string `json:"title,omitempty"`
	Body   string `json:"body,omitempty"`
}

// KV holds comments counts
type KV struct {
	Key   int
	Value int
}

// Comments hold comments response
type Comments struct {
	PostID int    `json:"postId,omitempty"`
	ID     int    `json:"id,omitempty"`
	Email  string `json:"email,omitempty"`
	Body   string `json:"body,omitempty"`
	Title  string `json:"title,omitempty"`
}

// Result holds response to user
type Result struct {
	PostID       int    `json:"postId,omitempty"`
	Title        string `json:"title,omitempty"`
	Body         string `json:"body,omitempty"`
	CommentCount int    `json:"commentcount,omitempty"`
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/top/{maxRes}", getTopX)

	log.Println("Running..")
	http.ListenAndServe(":8080", r)
}

func getTopX(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	maxRes := vars["maxRes"]

	maxResInt, err := strconv.Atoi(maxRes)

	log.Printf("MAX RES: %d", maxResInt)

	if err != nil {
		log.Fatalln("Unable to parse max result")
	}

	// get all comments

	cs := getAllComments()

	// count comment for each post
	postCommentsMap := make(map[int]int)
	for _, c := range cs {
		postCommentsMap[c.PostID]++
	}

	postCommentsMap[10] = 9
	postCommentsMap[20] = 12
	postCommentsMap[30] = 20

	// start := time.Now()
	// ss := mergeSort(postCommentsMap)
	// log.Println("merge sort took ", time.Since(start))
	// for _, kv := range ss[0:maxResInt] {
	// 	log.Printf("POSTID %d => %d comments\n", kv.(KV).Key, kv.(KV).Value)
	// }
	start := time.Now()
	sb := builtinSort(postCommentsMap)
	log.Println("buitin sort took ", time.Since(start))
	var results []Result
	for _, kv := range sb[0:maxResInt] {
		post := getPostDetails(kv.Key, kv.Value)
		results = append(results, post)
		log.Println(post)
	}

	b, err := json.Marshal(results)

	w.Write(b)

}

// Get posts details

func getPostDetails(id int, count int) Result {

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%d", postsURL, id), nil)

	if err != nil {
		log.Println("Request failed")

	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("client do")
		log.Println(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	var ps Post
	err = json.Unmarshal(body, &ps)
	if err != nil {
		log.Println(err)
	}

	var rs Result
	rs.PostID = ps.ID
	rs.Body = ps.Body
	rs.Title = ps.Title
	rs.CommentCount = count
	return rs

}

func getAllComments() []Comments {
	// get all comments

	client := http.Client{
		Timeout: 20 * time.Second,
	}
	req, err := http.NewRequest("GET", commentsURL, nil)

	if err != nil {
		log.Println("Request failed")

	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("client do")
		log.Println(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	var cs []Comments
	err = json.Unmarshal(body, &cs)
	if err != nil {
		log.Println(err)
	}

	return cs
}

// Quicksort

func builtinSort(m map[int]int) []KV {
	var ss []KV
	for k, v := range m {
		ss = append(ss, KV{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	return ss
}

// Mergesort

func mergeSort(m map[int]int) []interface{} {
	ss := make([]interface{}, len(m))
	i := 0
	for k, v := range m {
		ss[i] = KV{k, v}
		i++
	}

	timsort.Sort(ss, func(a, b interface{}) bool {
		return a.(KV).Value > b.(KV).Value
	})

	return ss
}
