package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"
)

type ErrorDef struct {
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

type Response struct {
	ErrorCode string `json:"error_code"`
	Key       string `json:"key"`
}

func main() {
	// configure the songs directory name and port
	const m3u8Dir = "vedios_m3u8"
	const port = 8080
	const md5Key = "Taipei_Pass_Stream"

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {

		// header  md5 token
		reqToken := r.Header.Get("TP-STREAM")
		if len(reqToken) <= 0 {
			// w.Write([]byte("not input token!!"))
			errorDef := ErrorDef{"-1", "not input token!!"}
			js, err := json.Marshal(errorDef)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}
		currentDate := time.Now().Local().Format("2006-01-02")
		h := md5.New()
		h.Write([]byte(currentDate + md5Key))
		md5String := hex.EncodeToString(h.Sum(nil))

		if md5String != reqToken {
			// fmt.Fprintf(w, "token no match!!")
			errorDef := ErrorDef{"-1", "token no match!!"}
			js, err := json.Marshal(errorDef)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}

		fmt.Fprintf(w, "md5: %s, token : %s", md5String, reqToken)
	})

	http.HandleFunc("/kell_ffmpeg", Handler)
	http.HandleFunc("/ping1", func(w http.ResponseWriter, r *http.Request) {
		cmdArguments := []string{"-i", "http://111.235.240.69:5552/Z01CC-001BFE06CC1E?key=?", "-c", "copy",
			"-f", "hls", "-hls_time", "0.1", "-hls_list_size", "100", "./vedios_m3u8/test.m3u8"}
		cmd := exec.Command("ffmpeg", cmdArguments...)

		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("command output: %d", cmd.Process.Pid)
		fmt.Fprintf(w, "command output: %d", cmd.Process.Pid)
	})

	// add a handler for the song files
	http.Handle("/", addHeaders(http.FileServer(http.Dir(m3u8Dir))))
	fmt.Printf("Starting server on %v\n", port)
	log.Printf("Serving %s on HTTP port: %v\n", m3u8Dir, port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

//addHeaders will act as middleware to give us CORS support
func addHeaders(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {

	ids, ok := r.URL.Query()["id"]

	if !ok || len(ids[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		fmt.Fprintf(w, "Url Param 'key' is missing")
		return
	}

	id := ids[0]
	u := &Response{
		ErrorCode: "0",
		Key:       string(id),
	}
	cmd := exec.Command("kill", string(id))
	cmd.Start()

	b, err := json.Marshal(u)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
