package server

// import (
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"os"
// 	"path/filepath"
// 	"strconv"
// 	"strings"

// 	"github.com/gorilla/mux"
// )

// func picHandler(w http.ResponseWriter, r *http.Request) {
// 	if !validation(w, r) {
// 		return
// 	}

// 	_, valid := sessionValid(w, r)
// 	if !valid {
// 		return
// 	}

// 	path := mux.Vars(r)

// 	switch path["key"] {
// 	case "help":
// 		fmt.Fprint(w, PicAPI)
// 	case "upload":
// 		picUpload(w, r)
// 	case "get":
// 		getPic(w, r)
// 	case "changeBg":
// 		changeBg(w, r)
// 	default:
// 		http.NotFound(w, r)
// 	}
// }

// func picUpload(w http.ResponseWriter, r *http.Request) {

// 	r.ParseMultipartForm(32 << 20)
// 	file, handler, err := r.FormFile("uploadfile")
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	defer file.Close()

// 	fmt.Fprint(w, handler.Header)
// 	f, err := os.Create("webpage/img/" + handler.Filename)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	defer f.Close()

// 	io.Copy(f, file)
// }

// func getPic(w http.ResponseWriter, r *http.Request) {
// 	args := r.URL.Query()
// 	pdid := args.Get("pdid")
// 	_, err := strconv.Atoi(pdid)
// 	if err != nil {
// 		fmt.Fprint(w, "pdid is not an integer")
// 		return
// 	}

// 	// find file name without knowing extention
// 	name, _ := filepath.Glob("webpage/img/" + pdid + ".*")
// 	if len(name) > 0 {
// 		splited := strings.Split(name[0], "/")
// 		fmt.Fprint(w, splited[len(splited)-1])
// 	} else {
// 		fmt.Fprint(w, "none.webp")
// 		log.Println("img not found")
// 	}
// }

// func changeBg(w http.ResponseWriter, r *http.Request) {
// 	//r.ParseMultipartForm(32 << 20)

// 	file, _, err := r.FormFile("uploadfile")
// 	if err != nil {
// 		log.Println("at upload:", err)
// 		return
// 	}
// 	defer file.Close()

// 	f, err := os.Create("webpage/img/temp.webp")
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	io.Copy(f, file)
// 	f.Close()

// 	err = os.Remove("webpage/img/bg2.webp")
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	err = os.Rename("webpage/img/temp.webp", "webpage/img/bg2.webp")
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	fmt.Fprint(w, " ok")
// }
