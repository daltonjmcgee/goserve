package admin

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"goserve/config"
	"goserve/helpers"
	"goserve/httpErrorHandler"
	"html/template"
	"net/http"
	"strings"

	"golang.org/x/crypto/argon2"
)

var adminPath string = "./admin/public/"
var conf map[string]string = config.ReturnConfig("config.dev.json")

func getDb(w http.ResponseWriter, conf map[string]string) (error, map[string][]interface{}) {
	jsonBytes, err := helpers.LoadFile(conf["databasePath"])

	if err != nil {
		httpErrorHandler.Handle500(w)
		return err, nil
	}

	jsonMap := map[string][]interface{}{}
	json.Unmarshal([]byte(jsonBytes), &jsonMap)

	return nil, jsonMap
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	dbErr, jsonMap := getDb(w, conf)

	if dbErr != nil {
		return
	}

	t, err := template.ParseFiles(adminPath + path[1:] + ".html")

	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			httpErrorHandler.Handle404(w)
		} else {
			httpErrorHandler.Handle500(w)
		}
	} else {
		t.Execute(w, jsonMap)
	}
	return
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpErrorHandler.Handle405(w, r.Method)
		return
	}

	r.ParseForm()
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	dbErr, jsonMap := getDb(w, conf)

	if dbErr != nil {
		return
	}

	// Loop over all users and look for match to login info
	for _, val := range jsonMap["users"] {
		creds, ok := val.(map[string]interface{})

		if !ok {
			fmt.Fprintf(w, "type map[string]interface{} required; got %T", val)
		}

		// check to see if email exists first, then check to see if password is correct
		if creds["email"] == email {
			// hash given password with stored salt and convert to hex string to compare against database
			hashedPwd := hex.EncodeToString(argon2.IDKey([]byte(password), []byte(creds["salt"].(string)), 1, 64*1024, 4, 32))
			if creds["password"] == hashedPwd {
				http.Redirect(w, r, "/admin/panel", http.StatusFound)
				return
			}
		}
	}

	fmt.Fprint(w, "Username or Password could not be verified.")
	return
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpErrorHandler.Handle405(w, r.Method)
		return
	}

	// jsonMap := getDb(w)

	r.ParseForm()
	// var request []string
	for key, val := range r.PostForm {
		fmt.Println("\n", key, strings.Join(val, "\n"))
	}
	fmt.Fprintf(w, "DONE")
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
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
	http.HandleFunc("/admin/update", handleUpdate)
	http.HandleFunc("/admin/static/", handleStatic)
}
