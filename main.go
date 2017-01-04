package main

import (
	"net/http"
	"log"
	"time"
	"encoding/json"
	"html/template"
	"os"
	"io"
	b64 "encoding/base64"
	"strconv"
	"path/filepath"
	"strings"
)

// stats is a container for general usage statistics
type Stats struct {
	StartupTime                      time.Time
	FileCount, UploadCount, GetCount int
}

type Config struct {
	Port, Path string
}

type UploadResponse struct {
	Filename, OriginalFilename string
	Uploaded                   time.Time
}

var statistics Stats
var configuration Config
var templates *template.Template
var filesMap map[string]string

func main() {
	initializeConfiguration()
	initializeFileMap()
	initializeStatistics()
	initializeTemplate()

	http.HandleFunc("/stats", StatsHandler)
	http.HandleFunc("/fetch/", Fetch)
	http.HandleFunc("/upload", Upload)
	http.HandleFunc("/", StaticHandler)

	log.Printf("Starting server on port: '%s'", configuration.Port)
	log.Printf("Serving/Storing files from: '%s'", configuration.Path)

	http.ListenAndServe(":" + configuration.Port, nil)
}

func initializeTemplate() {
	var err error
	templates, err = template.ParseFiles("index.html");
	if err != nil {
		log.Fatal("Abort: Failed parsing template")
	}
}

func initializeStatistics() {
	statistics = Stats{StartupTime:time.Now(),
		FileCount:len(filesMap)}
}

func initializeConfiguration() {
	configuration.Path = "tmp"
	configuration.Port = "8080"
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

// StatsHandler writes usage statistics as json to the response.
func StatsHandler(w http.ResponseWriter, r *http.Request) {
	statistics.FileCount = len(filesMap)
	json.NewEncoder(w).Encode(statistics)
}

func Fetch(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Path
	tokens := strings.Split(url, "/")
	dirname := tokens[len(tokens) - 1]

	filename, found := filesMap[dirname]
	if !found {
		msg := "Fetch: Did not find file"
		log.Print(msg)
		http.Error(w, msg, http.StatusNotFound)
		return
	}

	f, err := os.Open(configuration.Path + "/" + dirname + "/" + filename)
	if err != nil {
		log.Printf("Upload: Failed reading file: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Disposition", "attachment; filename=\"" + filename + "\"")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	io.Copy(w, f)

	statistics.GetCount++
}

func Upload(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("uploadfile")

	if err != nil {
		log.Printf("Upload: Failed reading form value: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	data := strconv.Itoa(time.Now().Nanosecond()) + handler.Filename
	directoryName := b64.StdEncoding.EncodeToString([]byte(data))[:15]

	if err := os.Mkdir(configuration.Path + "/" + directoryName, 0755); err != nil {
		log.Printf("Upload: Failed creating directory: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	f, err := os.OpenFile(configuration.Path + "/" + directoryName + "/" + handler.Filename, os.O_WRONLY | os.O_CREATE, 0666)
	if err != nil {
		log.Printf("Upload: Failed writing file: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	io.Copy(f, file)
	filesMap[directoryName] = handler.Filename

	response, _ := json.Marshal(UploadResponse{
		Filename:directoryName,
		OriginalFilename: handler.Filename,
		Uploaded:time.Now(),
	})
	statistics.UploadCount++

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func initializeFileMap() {

	filesMap = make(map[string]string)

	// create targetDir if not exists
	if _, err := os.Stat(configuration.Path); os.IsNotExist(err) {
		os.Mkdir(configuration.Path, 0755)
		log.Print("Init: Created missing base directory")
	}

	visitor := func(path string, f os.FileInfo, err error) error {
		cleanPath := path[len(configuration.Path):]
		if tokens := strings.Split(cleanPath, "/"); len(tokens) == 3 {
			filesMap[tokens[1]] = tokens[2]
		}
		return nil
	}

	if err := filepath.Walk(configuration.Path, visitor); err != nil {
		log.Fatalf("Init: Failed to load files from directory: %s", err)
	}

	log.Printf("Init: Loaded '%d' files", len(filesMap))
}