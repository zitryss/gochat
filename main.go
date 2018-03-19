package main

import (
	"net/http"
	"text/template"

	"fmt"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/objx"
)

type handler struct {
	tmpl *template.Template
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{
		"Host": req.Host,
	}
	if authCookie, err := req.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	h.tmpl.Execute(resp, data)
}

func main() {
	gomniauth.SetSecurityKey("PUT YOUR AUTH KEY HERE")
	gomniauth.WithProviders(
		github.New("a126fa1fa26ecd304407", "0ff076d3e6ed46ddb1c75ee37d1007906e785f8b", "http://localhost:9000/auth/callback/github"),
	)

	r := newRoom()
	go r.run()

	tmpl1 := template.Must(template.ParseFiles("chat.gohtml"))
	tmpl2 := template.Must(template.ParseFiles("login.gohtml"))
	http.Handle("/chat", MustAuth(&handler{tmpl1}))
	http.Handle("/login", &handler{tmpl2})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	err := http.ListenAndServe("localhost:9000", nil)
	if err != nil {
		fmt.Print(err)
	}
}
