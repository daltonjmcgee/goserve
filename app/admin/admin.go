package admin

import (
	"encoding/json"
	"fmt"
	"goserve/config"
	"goserve/helpers"
	"goserve/httpErrorHandler"
	"net/http"
)

var adminPath string = "./admin/public/"
var conf map[string]string = config.ReturnConfig("config.dev.json")

func handleGet(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	file, err := helpers.LoadFile(adminPath + path[1:] + ".html")

	if err != nil {
		httpErrorHandler.Handle404(w)
	} else {
		fmt.Fprintf(w, file)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpErrorHandler.Handle405(w, r.Method)
	} else {
		r.ParseForm()
		email := r.PostForm.Get("email")
		password := r.PostForm.Get("password")
		jsonBytes, err := helpers.LoadFile(conf["databasePath"])

		if err != nil {
			fmt.Fprint(w, "Database could not be found.")
			return
		}

		jsonMap := map[string][]interface{}{}
		json.Unmarshal([]byte(jsonBytes), &jsonMap)

		for _, val := range jsonMap["users"] {
			creds, ok := val.(map[string]interface{})
			if !ok {
				fmt.Fprintf(w, "type map[string]interface{} required; got %T", val)
			}
			if creds["email"] == email && creds["password"] == password {
				fmt.Fprint(w, "Login Successful")
			} else {
				fmt.Fprint(w, "Username or Password could not be verified.")
			}
		}
	}
	return
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Println(path)
	if path != "/admin/static/" {
		http.ServeFile(w, r, adminPath+path)
	} else {
		httpErrorHandler.Handle404(w)
	}
}

func AdminPanel() {
	http.HandleFunc("/admin", handleGet)
	http.HandleFunc("/admin/", handleGet)
	http.HandleFunc("/admin/login", handleLogin)
	http.HandleFunc("/admin/static/", handleStatic)
}
