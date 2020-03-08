package main

import (
	"database/sql"
	"encoding/json"
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
	jwt             *JwtHelper
	session_storage Storage
}

type HttpServer struct {
	*http.Server
	drive *GoogleDrive
	db    *sql.DB
}

func (self *SessionServMux) getSession(r *http.Request) *SessionData {
	signature, err := r.Cookie("session")

	if err != nil {
		return NewSessionData()
	}
	data := []byte(signature.Value)
	claim, err := self.jwt.Verify(data)
	if err != nil {
		return NewSessionData()
	}
	id, ok := claim.Set["sid"]
	if !ok {
		return NewSessionData()
	}
	sid, ok := id.(string)
	if !ok {
		return NewSessionData()
	}
	session, err := self.session_storage.Get(sid)
	if err != nil {
		log.Printf("fetch session error, reason: %s\n", err.Error())
		return NewSessionData()
	}
	return session.(*SessionData)
}

func (self *SessionServMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	session := self.getSession(r)
	wrap := &ResponseWriter{
		ResponseWriter: w,
		Session:        session,
		callback: func(w *ResponseWriter) {
			self.session_storage.Set(session)
			session := map[string]interface{}{
				"sid": w.Session.GetId(),
			}
			data, err := self.jwt.Sign(
				session,
				time.Now().AddDate(0, 0, 1),
			)
			if err == nil {
				http.SetCookie(w, &http.Cookie{
					Name:     "session",
					Value:    string(data),
					HttpOnly: true,
					Secure:   false,
				})
			} else {
				log.SetFlags(log.LstdFlags | log.Lshortfile)
				log.Println(err)
			}
		},
	}

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
	Session    *SessionData
	statusCode int
	done       bool
	callback   func(*ResponseWriter)
}

func (self *ResponseWriter) Header() http.Header {
	return self.ResponseWriter.Header()
}

func (self *ResponseWriter) Write(b []byte) (int, error) {
	if self.Session.updated && !self.done {
		self.callback(self)
	}
	self.done = true
	return self.ResponseWriter.Write(b)
}

func (self *ResponseWriter) WriteHeader(code int) {
	if self.Session.updated && !self.done {
		self.callback(self)
	}
	self.done = true
	self.statusCode = code
	self.ResponseWriter.WriteHeader(code)
}

func (self *ResponseWriter) WriteJson(body interface{}) (int, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return 0, err
	}
	return self.Write(data)
}
