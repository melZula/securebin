package apiserver

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/sethvargo/go-password/password"

	"github.com/google/uuid"

	"github.com/gorilla/handlers"

	"github.com/melZula/securebin/internal/app/model"

	"github.com/gorilla/mux"
	"github.com/melZula/securebin/internal/app/store"
	"github.com/sirupsen/logrus"
)

const (
	sessionName        = "zulasession"
	ctxKeyData  ctxKey = iota
	ctxKeyRequestID
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password") // cover real error
	errNotAuthenticated         = errors.New("not authenticated")
	errNoContent                = errors.New("no content")
)

type ctxKey int8

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/api/paste", s.handlePostCreate()).Methods("POST")
	s.router.HandleFunc("/api/data", s.handleSessionsCreate()).Methods("POST")
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-request-id", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)
		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		logger.Infof(
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (s *server) handlePostCreate() http.HandlerFunc {
	type request struct {
		Text string `json:"text"`
		Time int64  `json:"time"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		timeout := time.Now().UTC().Add(time.Duration(req.Time * 1000000000))

		image, err := renderImage(req.Text)
		if err != nil {
			s.error(w, r, http.StatusTeapot, err)
			return
		}

		pswd, err := password.Generate(10, 4, 6, false, true)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		u := &model.Data{
			Img:      image,
			Password: pswd,
			Lifetime: timeout.Unix(),
		}
		if err := s.store.Data().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		ID       int    `json:"id"`
		Password string `json:"password"`
	}
	type response struct {
		Img     string `json:"img"`
		PrevReq []int  `json:"times"`
	}
	type histEl struct {
		ID   int   `json:"id"`
		Time int64 `json:"time"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		res := &response{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.Data().Find(req.ID)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		if !u.IsAlive() {
			s.error(w, r, http.StatusNoContent, errNoContent)
			return
		}

		res.Img = u.Img
		res.PrevReq, err = s.store.Data().GetPrevRequests(req.ID)
		if err != nil {
			s.logger.Warning("Can not get previous requests")
		}

		if err := s.store.Data().LogRequest(req.ID, r); err != nil {
			s.logger.Warning("Log of request hadn't recorded")
		}

		var h []histEl
		history, err := r.Cookie("history")
		expiration := time.Now().UTC().Add(365 * 24 * time.Hour)
		if err != nil {
			if err != http.ErrNoCookie {
				s.logger.Warning("Unable to get cookie")
			}
		} else {

			buf, err := base64.StdEncoding.DecodeString(history.Value)
			if err != nil {
				s.logger.Warning("Can't decode b64: ", err)
			}
			if err := json.Unmarshal(buf, &h); err != nil {
				s.logger.Warning("Can't unmarshal cookie: ", err)
			}
		}

		h = append(h, histEl{ID: req.ID, Time: time.Now().UTC().Unix()})
		c, err := json.Marshal(h)
		http.SetCookie(w, &http.Cookie{Name: "history", Value: base64.StdEncoding.EncodeToString(c), Expires: expiration, Path: "/"})

		s.respond(w, r, http.StatusOK, *res)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}

}
