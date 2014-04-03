package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/kurrik/oauth1a"
	"log"
	"net/http"
	"os"
)

var (
	key    = flag.String("key", "", "API key")
	secret = flag.String("secret", "", "API secret")
	addr   = flag.String("addr", ":10000", "bind address")
)

var service *oauth1a.Service
var sessions = make(map[string]*oauth1a.UserConfig)

func rootHandler(w http.ResponseWriter, req *http.Request) {
	var (
		url string
		err error
	)
	httpClient := new(http.Client)
	userConfig := &oauth1a.UserConfig{}
	if err = userConfig.GetRequestToken(service, httpClient); err != nil {
		http.Error(w, "request token error: "+err.Error(), 500)
		return
	}
	if url, err = userConfig.GetAuthorizeURL(service); err != nil {
		http.Error(w, "auth url error: "+err.Error(), 500)
		return
	}

	// Generate session id.
	b := make([]byte, 128)
	if _, err := rand.Read(b); err != nil {
		panic("rand error: " + err.Error())
	}
	sessionID := base64.URLEncoding.EncodeToString(b)

	log.Printf("new session: %v\n", sessionID)
	sessions[sessionID] = userConfig
	http.SetCookie(w, &http.Cookie{Name: "session_id", Value: sessionID, MaxAge: 60, Secure: false, Path: "/"})
	http.Redirect(w, req, url, 302)
}

func callbackHandler(w http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("session_id")
	if err != nil {
		http.Error(w, "cookie error: "+err.Error(), 400)
		return
	}
	sessionID := c.Value
	log.Printf("callback: %v.", sessionID)

	userConfig, ok := sessions[sessionID]
	if !ok {
		http.Error(w, "invalid session: "+err.Error(), 400)
		return
	}

	token, verifier, err := userConfig.ParseAuthorize(req, service)
	if err != nil {
		http.Error(w, "auth parse error: "+err.Error(), 500)
		return
	}

	httpClient := new(http.Client)
	if err = userConfig.GetAccessToken(token, verifier, service, httpClient); err != nil {
		log.Printf("Error getting access token: %v", err)
		http.Error(w, "access token error: "+err.Error(), 500)
		return
	}

	delete(sessions, sessionID)
	http.SetCookie(w, &http.Cookie{Name: "session_id", Value: "", MaxAge: 0, Secure: false, Path: "/"})

	w.Header().Set("Content-Type", "text/plain;charset=utf-8")
	fmt.Fprintf(w, "Screen Name:         %v\n", userConfig.AccessValues.Get("screen_name"))
	fmt.Fprintf(w, "Access Token Key:    %v\n", userConfig.AccessTokenKey)
	fmt.Fprintf(w, "Access Token Secret: %v\n", userConfig.AccessTokenSecret)
}

func main() {
	log.SetFlags(0)
	flag.Parse()

	if *key == "" {
		fmt.Fprintln(os.Stderr, "key required: -key.")
		flag.PrintDefaults()
		os.Exit(1)
	} else if *secret == "" {
		fmt.Fprintln(os.Stderr, "secret required: -secret.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	service = &oauth1a.Service{
		RequestURL:   "https://api.twitter.com/oauth/request_token",
		AuthorizeURL: "https://api.twitter.com/oauth/authorize",
		AccessURL:    "https://api.twitter.com/oauth/access_token",
		ClientConfig: &oauth1a.ClientConfig{
			ConsumerKey:    *key,
			ConsumerSecret: *secret,
			CallbackURL:    fmt.Sprintf("http://localhost%s/callback/", *addr),
		},
		Signer: new(oauth1a.HmacSha1Signer),
	}

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/callback/", callbackHandler)
	log.Printf("Listening on http://localhost%v", *addr)
	log.SetFlags(log.LstdFlags)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
