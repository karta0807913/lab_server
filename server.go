package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type ServerSettings struct {
	PublicKeyPath, PrivateKeyPath string
	ServerAddress                 string
	Storage                       Storage
	Db                            *sql.DB
	Drive                         *GoogleDrive
}

func NewSessionHttpServer(config ServerSettings) (*HttpServer, error) {
	jwt, err := NewJwtHelper(config.PublicKeyPath, config.PrivateKeyPath)
	if err != nil {
		return nil, err
	}
	handler := &SessionServMux{
		jwt:             jwt,
		session_storage: config.Storage,
	}
	server := &HttpServer{
		Server: &http.Server{
			Handler: handler,
			Addr:    config.ServerAddress,
		},
		drive: config.Drive,
		db:    config.Db,
	}
	return server, nil
}

type SessionServMux struct {
	http.ServeMux
	jwt *JwtHelper
}

type SessionData struct {
	session map[string]interface{}
	updated bool
}

func (self *SessionData) Get(key string) interface{} {
	return self.session[key]
}

func (self *SessionData) Set(key string, value interface{}) {
	self.updated = true
	self.session[key] = value
}

func (self *SessionData) Del(key string) {
	_, ok := self.session[key]
	self.updated = self.updated || ok
	delete(self.session, key)
}

type HttpServer struct {
	*http.Server
	drive *GoogleDrive
	db    *sql.DB
}

func (self *SessionServMux) getSession(r *http.ReadRequest) map[stirng]interface{} {
	signature, err := r.Cookie("session")

	if err != nil {
		return make(map[string]interface{})
	}
	claim, err := self.jwt.Verify([]byte(signature.Value))
	if err != nil {
		return make(map[string]interface{})
	}
	sid, ok := id.(string)
	if !ok {
		return make(map[string]interface{})
	}
	sid, ok := id.(string)
	if !ok {
		return make(map[string]interface{})
	}
	var session map[string]interface{}
	err = self.session_storage.Get(id, session)
	if err != nil {
		log.Println("fetch session error %s", err.Error())
		return make(map[string]interface{})
	}
	return session
}

func (self *SessionServMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	sessoin := self.getSession(r)

	wrap := WrapResponseWriter(
		w, self.jwt,
		&SessionData{
			session: session,
			updated: false,
		},
	)

	self.ServeMux.ServeHTTP(wrap, r)
	if !wrap.done {
		w.WriteHeader(404)
		w.Write([]byte("404 not found"))
	}
	if wrap.statusCode == 0 {
		wrap.statusCode = 200
	}
	log.Printf("%s %s %d\n", r.Host, r.URL.Path, wrap.statusCode)
}

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
	done       bool
	Session    *SessionData
	jwt        *JwtHelper
}

func (self *ResponseWriter) Header() http.Header {
	return self.ResponseWriter.Header()
}

func (self *ResponseWriter) Write(b []byte) (int, error) {
	if self.Session.updated && !self.done {
		data, err := self.jwt.Sign(self.Session.session, time.Now())
		if err == nil {
			self.Header().Add("Set-Cookie", string(data))
		} else {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Println(err)
		}
	}
	self.done = true
	return self.ResponseWriter.Write(b)
}

func (self *ResponseWriter) WriteHeader(code int) {
	if self.Session.updated && !self.done {
		data, err := self.jwt.Sign(self.Session.session, time.Now())
		if err == nil {
			http.SetCookie(self, &http.Cookie{
				Name:     "session",
				Value:    string(data),
				Secure:   true,
				HttpOnly: true,
			})
		} else {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Println(err)
		}
	}
	self.done = true
	self.statusCode = code
	self.ResponseWriter.WriteHeader(code)
}

func WrapResponseWriter(origin http.ResponseWriter, jwt *JwtHelper, session *SessionData) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: origin,
		Session:        session,
		jwt:            jwt,
	}
}
