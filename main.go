package main

import (
    "net/http"
    "fmt"
    "log"
    "io/ioutil"
    "os"
    "encoding/json"
    "strings"
    "io"
)

type File struct {
    Path string
    Content string
}

const root = "root/"

func main(){
    http.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir("UI"))))
    http.Handle("/explorer/", http.StripPrefix("/explorer/", http.FileServer(http.Dir("root"))))

    http.HandleFunc("/file/", api)
    http.HandleFunc("/copy/", copyApi)
    http.HandleFunc("/move/", moveApi)
    http.HandleFunc("/dir/", dirApi)

    log.Println("Listening in port 8000...")
    http.ListenAndServe(":8000", nil)
}

func api(w http.ResponseWriter, r *http.Request){
    dir := r.URL.Path[len("/file/"):]
    switch r.Method {
    case "OPTIONS":
        w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
        w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")
        w.Header().Set("Access-Control-Allow-Headers", "Acccept, Content-Type, Content-Length")
        log.Println("preflight response...done")
    case "GET":
        file, err := ioutil.ReadFile(root + dir)
        if err != nil {
            http.Error(w, "{\"message\": \"Something went wrong\"}", http.StatusBadRequest)
        }else{
            fmt.Fprintf(w, "{\"file\": \"" + dir + "\", \"content\": \"" + string(file) + "\"}")
        }
    case "POST":
        var file File

        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&file)
        if err != nil {
            http.Error(w, "{\"message\": \"Something went wrong\"}", http.StatusInternalServerError)
        }else{
            newFile, err := os.Create(root + dir)
            if err != nil{
                http.Error(w, "{\"message\": \"Cannot create a file\"}", http.StatusInternalServerError)
            }else{
                _, err := newFile.WriteString(file.Content)
                if err != nil{
                    http.Error(w, "{\"message\": \"Cannot write into the file\"}", http.StatusInternalServerError)
                    os.Remove(root + dir)
                }else{
                    http.Error(w, "{\"message\": \"Successfully created " + dir + "\"}", http.StatusOK)
                }
            }
            defer newFile.Close()
        }
    case "PUT":
        var file File

        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&file)
        if err != nil {
            http.Error(w, "{\"message\": \"Something went wrong\"}", http.StatusInternalServerError)
        }else{
            f, err := os.Open(root + dir)
            if err != nil {
                http.Error(w, "{\"message\": \"File does not exist\"}", http.StatusInternalServerError)
            }else{
                _, err := f.Write([]byte(file.Content))
                if err != nil {
                    log.Println(err)
                    http.Error(w, "{\"message\": \"Cannot edit the file\"}", http.StatusInternalServerError)
                }else{
                    http.Error(w, "{\"message\": \"Modifications are saved\"}", http.StatusOK)
                }
            }
        }
    case "DELETE":
        err := os.Remove(root + dir)
        if err != nil {
            http.Error(w, "{\"message\": \"Cannot delete the file\"}", http.StatusInternalServerError)
        }else{
            http.Error(w, "{\"message\": \"success\"}", http.StatusOK)
        }
    default:
        http.Error(w, "{\"message\": \"Invalid Method\"}", http.StatusBadRequest)
    }
}

func copyApi(w http.ResponseWriter, r *http.Request){
    dir := r.URL.Path[len("/copy/"):]
    if r.Method == "POST" {
        if _, err := os.Stat(root + dir); err == nil{
            file, err := os.Open(root + dir)
            if err != nil {
                http.Error(w, "{\"message\": \"Something went wrong\"}", http.StatusInternalServerError)
            }else{
                tmp := strings.Split(dir, "/")
                filename := tmp[len(tmp) - 1]
                var newfile File

                decoder := json.NewDecoder(r.Body)
                err := decoder.Decode(&newfile)
                if err != nil {
                    http.Error(w, "{\"message\": \"Something went wrong\"}", http.StatusInternalServerError)
                }else{
                    f, _ := os.Create(root + newfile.Path + filename)
                    _, err := io.Copy(f, file)
                    defer f.Close()
                    if err != nil {
                        http.Error(w, "{\"message\": \"Something went wrong\"}", http.StatusInternalServerError)
                    }else{
                        http.Error(w, "{\"message\": \"Directory created uccessfully\"}", http.StatusOK)
                    }
                }
                defer file.Close()
            }
        }else{
            http.Error(w, "{\"message\": \"File does not exist\"}", http.StatusInternalServerError)
        }
    }else{
        http.Error(w, "{\"message\": \"Invalid Method\"}", http.StatusBadRequest)
    }
}

func moveApi(w http.ResponseWriter, r *http.Request){
    dir := r.URL.Path[len("/move/"):]
    if r.Method == "POST" {
        tmp := strings.Split(dir, "/")
        filename := tmp[len(tmp) - 1]
        var newfile File

        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&newfile)
        if err != nil {
            http.Error(w, "{\"message\": \"Something went wrong\"}", http.StatusInternalServerError)
        }else{
            err := os.Rename(dir, newfile.Path + filename)
            if err != nil {
                http.Error(w, "{\"message\": \"Something went wrong\"}", http.StatusInternalServerError)
            }else{
                http.Error(w, "{\"message\": \"Successfully moved the file\"}", http.StatusOK)
            }
        }
    }
}

func dirApi(w http.ResponseWriter, r *http.Request){
    dir := r.URL.Path[len("/dir/"):]
    if _, err := os.Stat(root + dir); os.IsNotExist(err){
        dir = root + dir
        os.MkdirAll(dir, os.ModePerm)
        http.Error(w, "{\"message\": \"Directory created successfully\"}", http.StatusOK)
    }else{
        http.Error(w, "{\"message\": \"Directory Already exist\"}", http.StatusInternalServerError)
    }
}
