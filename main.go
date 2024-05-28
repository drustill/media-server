package mediaserver

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ServeMux() {
	mediaDir := "./media" // Directory where your media files are stored
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Println("Request received: ", r.URL.Path)
			filePath := filepath.Join(mediaDir, r.URL.Path)
			log.Printf("Serving file: %s", filePath)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
					http.NotFound(w, r)
					return
			}

			file, err := os.Open(filePath)
			if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
			}
			defer file.Close()

			fileStat, err := file.Stat()
			if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
			}

			fileSize := fileStat.Size()

			w.Header().Set("Content-Type", "video/mp4")
			w.Header().Set("Accept-Ranges", "bytes")

			if r.Header.Get("Range") != "" {
					rangeHeader := r.Header.Get("Range")
					ranges := strings.Split(rangeHeader, "=")[1]
					parts := strings.Split(ranges, "-")
					start, _ := strconv.ParseInt(parts[0], 10, 64)
					end := fileSize - 1
					if len(parts) > 1 && parts[1] != "" {
							end, _ = strconv.ParseInt(parts[1], 10, 64)
					}

					if start > end || start < 0 || end >= fileSize {
							w.Header().Set("Content-Range", "bytes */"+strconv.FormatInt(fileSize, 10))
							http.Error(w, "Requested Range Not Satisfiable", http.StatusRequestedRangeNotSatisfiable)
							return
					}

					w.Header().Set("Content-Range", "bytes "+strconv.FormatInt(start, 10)+"-"+strconv.FormatInt(end, 10)+"/"+strconv.FormatInt(fileSize, 10))
					w.Header().Set("Content-Length", strconv.FormatInt(end-start+1, 10))
					w.WriteHeader(http.StatusPartialContent)

					file.Seek(start, 0)
					http.ServeContent(w, r, filePath, fileStat.ModTime(), file)
			} else {
					w.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))
					http.ServeContent(w, r, filePath, fileStat.ModTime(), file)
			}
	})

	log.Println("Listening on localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Fatalf("Failed to start server: %s", err)
	}
}
