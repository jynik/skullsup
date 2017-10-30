// SPDX License Identifier: MIT
package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"../../skullsup"

	c "../common"
	"../common/logger"
)

const maxLogDataLen = 256

type Server struct {
	config Config
	mq     *MessageQueue

	log   *logger.Logger
	skull *skullsup.Skull
}

type Config struct {
	Port uint16 // Port number to listen on

	ListenAddr     string // Address to bind listening socket to
	TlsCertPath    string // Path to TLS certificate
	PrivateKeyPath string // Path to TLS private key

	ClientConfigPath string // Path to client config file

	LogFilePath string // Where log output should be written. May be a Filename, "stdout", or "stderr"
	Verbose     bool   // Enable extra output
	Quiet       bool   // Supress all output
}

func (s *Server) logRequest(r *http.Request, logUser string, msg string, data string) {
	s.log.Printf("[%s@%s] %s: %s\n", logUser, r.RemoteAddr, msg, data)
}

func (s *Server) errorResponse(w http.ResponseWriter, r *http.Request, errMsg string, errCode int, logUser string, logMsg string, err error) {
	s.logRequest(r, logUser, logMsg, err.Error())
	http.Error(w, errMsg, errCode)
}


func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	var isWrite bool
	var msg c.Message

	username, password, ok := r.BasicAuth()
	if !ok {
		s.errorResponse(w, r, "Authentication required", 403, username, "Auth failure", errors.New("Credentials not provided"))
		return
	}

	switch r.Method {
	case http.MethodPut:
		isWrite = true
	case http.MethodGet:
		isWrite = false
	default:
		s.errorResponse(w, r, c.ERR_METHOD, 405, username, "Invalid method", errors.New(r.Method))
		return
	}

	queues := strings.Split(r.Header.Get(c.HEADER_QUEUE), c.HEADER_QUEUE_SEP)

	err := s.Authenticate(queues, username, password, isWrite)
	if err != nil {
		s.errorResponse(w, r, err.Error(), 403, username, "Auth failure", err)
		return
	}

	if isWrite {
		if r.ContentLength < 1 {
			s.errorResponse(w, r, c.ERR_INVAL, 666, username, "Unknown Content-Length", fmt.Errorf("%d", r.ContentLength))
			return
		}

		body := make([]byte, r.ContentLength)
		for _, queue := range queues {
			if _, err := r.Body.Read(body); err != nil && err != io.EOF {
				s.errorResponse(w, r, c.ERR_INVAL, 666, username, "Failed to read body", err)
			} else if err := json.Unmarshal(body, &msg); err != nil {
				s.errorResponse(w, r, c.ERR_INVAL, 666, username, "Failed to unpack command", err)
			} else if err := s.mq.Enqueue(queue, msg); err != nil {
				s.errorResponse(w, r, err.Error(), 666, username, "Failed to enqueue", err)
			}

			logMsg := msg.String()
			if len(logMsg) > maxLogDataLen {
				logMsg = logMsg[:maxLogDataLen-3] + "..."
			}
			s.logRequest(r, username, "Enqueued", logMsg)
		}
	} else {
		for _, queue := range queues {
			if msg, err := s.mq.Dequeue(queue); err != nil {
				s.errorResponse(w, r, err.Error(), 666, username, "Failed to dequeue", err)
				return
			} else if body, err := json.Marshal(msg); err != nil {
				s.errorResponse(w, r, c.ERR_FAILURE, 666, username, "Failed to marshal message", err)
			} else if _, err := w.Write(body); err != nil {
				s.errorResponse(w, r, c.ERR_FAILURE, 666, username, "Failed to write body", err)
			} else {
				// Return after successfully writing a dequeued a result
				logMsg := msg.String()
				if len(logMsg) > maxLogDataLen {
					logMsg = logMsg[:maxLogDataLen-3] + "..."
				}
				s.logRequest(r, username, "Dequeued", logMsg)
				return
			}
		}

		// Return an empty object if there was nothing to dequeue
		_, err := w.Write([]byte("{}"))
		if err != nil {
			s.errorResponse(w, r, c.ERR_FAILURE, 666, username, "Failed to write empty body", err)
		}
	}
}

func (s *Server) Run(config Config) error {
	var err error

	s.mq = NewMessageQueue(10, 10)
	s.config = config

	if !s.config.Quiet {
		s.log, err = logger.New(s.config.LogFilePath, s.config.Verbose)
	} else {
		s.log, err = logger.New("", false)
	}

	if err != nil {
		return err
	}

	http.HandleFunc(c.ENDPOINT, s.handleRequest)

	addr := s.config.ListenAddr + ":" + strconv.Itoa(int(s.config.Port))
	s.log.Println("Starting SkullsUp! server on " + addr + "...")
	return http.ListenAndServeTLS(addr, s.config.TlsCertPath, s.config.PrivateKeyPath, nil)
}
