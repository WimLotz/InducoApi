package main

import (
	"encoding/json"
	"fmt"
	"github.com/WimLotz/InducoApi/datastore"
	"github.com/WimLotz/InducoApi/user"
	"github.com/WimLotz/InducoApi/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	host = "localhost"
	port = "4567"
)

var sessionStore = sessions.NewCookieStore([]byte(utils.RandomString(32)))

//type (
//	appError struct {
//		Error   error
//		Message string
//		Code    int
//	}
//	appErrorWrapper func(http.ResponseWriter, *http.Request) *appError
//)
//
//func (fn appErrorWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	if e := fn(w, r); e != nil {
//		log.Println(e.Error)
//		http.Error(w, e.Message, e.Code)
//	}
//}
//
//func makeHandler(fn func(http.ResponseWriter, *http.Request, *sessions.Session) *appError) http.HandlerFunc {
//
//	return func(w http.ResponseWriter, r *http.Request) {
//		session, err := sessionStore.Get(r, "sessionName")
//		if err != nil {
//			log.Printf("Error occured retrieving session: %v\n", err)
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//		}
//
//		if appErr := fn(w, r, session); appErr != nil {
//			log.Println(appErr.Error)
//			http.Error(w, appErr.Message, appErr.Code)
//		}
//	}
//}
//
//func saveUser(w http.ResponseWriter, r *http.Request) *appError {
//
//	body := utils.ReadRequestBody(r.Body)
//
//	u := user.New()
//	utils.UnmarshalJsonToObject(body, &u)
//	u.Id = bson.NewObjectId()
//	u.Save()
//
//	return nil
//
//}

func login(w http.ResponseWriter, r *http.Request) {

	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")

	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}

	if r.Method == "POST" {
		defer r.Body.Close()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Unable to read body", err)
			return
		}

		u := user.New()

		err = json.Unmarshal(body, &u)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Unable to parse body", err)
			return
		}

		userSuppliedPassword := u.Password
		err = u.FetchOnEmail()
		if err != nil {
			w.WriteHeader(401)
			fmt.Fprintf(w, "Unable to find user", err)
			return
		} else {
			isCorrectPassword := u.IsSuppliedPasswordCorrect(userSuppliedPassword)
			fmt.Printf("pass:%v", isCorrectPassword)
			fmt.Printf("user:%v", u)
			if !isCorrectPassword {
				w.WriteHeader(401)
				fmt.Fprintf(w, "Wrong password")
				return
			}

			session, _ := sessionStore.Get(r, "sessionName")

			session.Values["userId"] = bson.ObjectId.Hex(u.Id)
			err = session.Save(r, w)
			if err != nil {
				w.WriteHeader(500)
				fmt.Fprintf(w, "Unable to save session")
				return
			}

			w.WriteHeader(200)
			return
		}
	}
	w.WriteHeader(400)
}

//func fetchUserProfiles(w http.ResponseWriter, r *http.Request, session *sessions.Session) *appError {
//
//	userId := session.Values["userId"]
//
//	if bson.IsObjectIdHex(userId.(string)) {
//		p := profile.New()
//		profiles := p.Fetch(bson.ObjectIdHex(userId.(string)))
//
//		w.Header().Set("Content-Type", "application/json")
//		w.Write(utils.MarshalObjectToJson(profiles))
//
//	} else {
//		return &appError{nil, "Error converting session userId to bson.ObjectId", http.StatusInternalServerError}
//	}
//
//	return nil
//}
//
//func saveProfile(w http.ResponseWriter, r *http.Request, session *sessions.Session) *appError {
//
//	body := utils.ReadRequestBody(r.Body)
//
//	p := profile.New()
//	utils.UnmarshalJsonToObject(body, &p)
//
//	userId := session.Values["userId"]
//
//	if bson.IsObjectIdHex(userId.(string)) {
//		if p.Id == "" {
//			p.Id = bson.NewObjectId()
//		}
//
//		p.UserId = bson.ObjectIdHex(userId.(string))
//		p.Save()
//	} else {
//		return &appError{nil, "Error converting session userId to bson.ObjectId", http.StatusInternalServerError}
//	}
//
//	return nil
//}

func main() {
	r := mux.NewRouter()

	//r.Handle("/saveUser", saveUser)
	r.HandleFunc("/login", login)
	//r.Handle("/saveProfile", makeHandler(saveProfile))
	//r.Handle("/fetchUserProfiles", makeHandler(fetchUserProfiles))
	//r.Handle("/fetchAllProfiles", makeHandler(fetchAllProfiles))

	//r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("."))))
	http.Handle("/", r)

	db := datastore.New()
	db.Connect()

	log.Printf("Server ready and listening on %v:%v", host, port)

	err := http.ListenAndServe(host+":"+port, nil)
	if err != nil {
		log.Fatal("Listen And Serve: ", err)
	}
}
