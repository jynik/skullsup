// SPDX License Identifier: MIT
package server

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/jynik/skullsup/go/src/logger"
	"github.com/jynik/skullsup/go/src/network"
)

type Server struct {
	cfg *Config
	q   *MessageQueue
	log *logger.Logger

	impl *http.Server

	queueFullLogged  bool
	queueEmptyLogged bool
}

type handlerContext struct {
	w      http.ResponseWriter
	r      *http.Request
	user   *network.User
	queue  string
	source string
}

func (s *Server) errorForbidden(w http.ResponseWriter, source string, reason string) {
	http.Error(w, "ph'nglui mglw'nafh Cthulhu R'lyeh wgah'nagl fhtagn", http.StatusForbidden)

	if reason != "" {
		reason = ": " + reason
	}

	s.log.Error("Auth failure from %s%s\n", source, reason)
}

func (s *Server) handleQueueWrite(ctx *handlerContext) {
	if !ctx.user.CanWrite(ctx.queue) {
		s.errorForbidden(ctx.w, ctx.source, "User does not have write permissions for requested queue.")
		return
	}

	if ctx.r.ContentLength < 1 {
		s.log.Error("No content received from %s\n", ctx.source)
		http.Error(ctx.w, "Empty queue message received", 400)
		return
	} else if ctx.r.ContentLength > 16384 {
		s.log.Error("Excessively large reuqest (%d) from %s\n", ctx.r.ContentLength, ctx.source)
		http.Error(ctx.w, "You're not worthy of such a grand request.", 413)
		return
	}

	defer ctx.r.Body.Close()
	body, err := ioutil.ReadAll(ctx.r.Body)
	if err != nil {
		s.log.Error("Failed to read body in request from %s: %s\n", ctx.source, err.Error())
		http.Error(ctx.w, "I've descended into maddness.", 500)
		return
	}

	var msg network.Message
	if err := json.Unmarshal(body, &msg); err != nil {
		s.log.Error("Received invalid message from %s\n", ctx.source)
		http.Error(ctx.w, "Your ramblings are incomrehensible!", 400)
		return
	}

	if err := s.q.Enqueue(ctx.queue, msg); err != nil {
		// Avoid filling logs with duplicate back-to-back error
		if err.Error() == network.ErrorQueueFull {
			if !s.queueFullLogged {
				s.log.Error("Queue full. Dropping enqueue request(s) from: %s\n", ctx.source)
				s.queueFullLogged = true
			} else {
				s.log.Debug("Queue full. Dropping enqueue request(s) from: %s\n", ctx.source)
			}
		} else {
			s.log.Error("Enqueue from %s failed: %s\n", ctx.source, err)
			s.queueFullLogged = false
		}
		http.Error(ctx.w, err.Error(), 666)
		return
	}
	s.queueFullLogged = false
	s.log.Debug("Enqueued message from %s: %s\n", ctx.source, msg.String())
}

func (s *Server) handleQueueRead(ctx *handlerContext) {
	if !ctx.user.CanRead(ctx.queue) {
		s.errorForbidden(ctx.w, ctx.source, "User does not have reader access to requested queue.")
		return
	}

	msg, err := s.q.Dequeue(ctx.queue)
	if err != nil {
		if err.Error() == network.ErrorQueueEmpty {
			if !s.queueEmptyLogged {
				s.log.Error("Queue empty. Dropping dequeue request(s) from: %s\n", ctx.source)
				s.queueEmptyLogged = true
			} else {
				s.log.Debug("Queue empty. Dropping dequeue request(s) from: %s\n", ctx.source)
			}
		} else {
			s.log.Error("Dequeue for %s failed: %s\n", ctx.source, err)
			s.queueEmptyLogged = false
		}

		http.Error(ctx.w, err.Error(), 666)
		return
	}
	s.queueEmptyLogged = false
	s.log.Debug("Dequeued message for %s: %s\n", ctx.source, msg.String())

	if body, err := json.Marshal(msg); err != nil {
		s.log.Error("Failed to marshal message for %s: %s\n", ctx.source, err)
		s.log.Error(" Message was: %s\n", msg.String())
		http.Error(ctx.w, "I've descended into maddness.", 500)
	} else if _, err := ctx.w.Write(body); err != nil {
		s.log.Error("Failed to write message to %s: %s\n", ctx.source, err)
		http.Error(ctx.w, "I have faltered, oh Dark One.", 500)
	}
}

func (s *Server) handleQueueRequest(w http.ResponseWriter, r *http.Request) {
	var user *network.User
	var err error

	for _, chain := range r.TLS.VerifiedChains {
		// We expect to see client -> our CA
		chainLen := len(chain)
		if chainLen != 2 {
			msg := "Expected verified chain of length 2. Skipping chain of length %d\n"
			s.log.Debug(msg, chainLen)
			continue
		}

		user, err = s.cfg.LookupUser(chain[0])
		if err != nil {
			s.errorForbidden(w, r.RemoteAddr, err.Error())
			return
		} else if user == nil {
			serial := hex.EncodeToString(chain[0].SerialNumber.Bytes())
			s.log.Debug("No user for certificate %s (from %s)\n", serial, r.RemoteAddr)
		}
	}

	if user == nil {
		s.errorForbidden(w, r.RemoteAddr, "No user associated with provided certificate.")
		return
	}

	source := fmt.Sprintf("%s@%s <%s>", user.Name, r.RemoteAddr, user.CertSerial)
	s.log.Debug("Handling request from: %s\n", source)

	queue := network.QueueFromURL(r.URL.Path)
	if queue == "" {
		// 403 instead of 404 to avoid enumerating other users' queues
		// ... because why not? :shrug:
		s.errorForbidden(w, source, "User requested non-existent queue.")
		s.log.Error("Request for non-existent from: %s\n", source)
		return
	}

	ctx := handlerContext{w: w, r: r, user: user, queue: queue, source: source}

	switch r.Method {
	case http.MethodPut:
		s.handleQueueWrite(&ctx)

	case http.MethodGet:
		s.handleQueueRead(&ctx)

	default:
		http.Error(w, "Shub-Niggurath!", http.StatusMethodNotAllowed)
		s.log.Error("%s@%s sent an invalid method.\n", user.Name, r.RemoteAddr)
	}
}

func New(filename string) (*Server, error) {
	var err error
	server := new(Server)

	server.cfg, err = loadConfig(filename)
	if err != nil {
		return nil, err
	}

	//  TODO Make configurable
	server.q = NewMessageQueue(10, 16)

	server.log, err = logger.New(server.cfg.LogPath, server.cfg.LogLevel)
	if err != nil {
		return nil, err
	}

	// Load CA certificate that is used to verify client certs
	caCert, err := ioutil.ReadFile(server.cfg.CAPath)
	if err != nil {
		return nil, err
	}

	server.log.Debug("Loaded CA Certificate: %s\n", server.cfg.CAPath)

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	server.impl = &http.Server{
		Addr:      fmt.Sprintf(":%d", int(server.cfg.Port)),
		TLSConfig: tlsConfig,
	}

	return server, nil
}

func (s *Server) Run() error {
	http.HandleFunc("/", s.handleQueueRequest)

	addr := s.cfg.Address + ":" + strconv.Itoa(int(s.cfg.Port))
	s.log.Info("Starting SkullsUp! server on " + addr + "\n")

	return s.impl.ListenAndServeTLS(s.cfg.CertPath, s.cfg.KeyPath)
}
