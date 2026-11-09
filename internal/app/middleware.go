package app

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
	"greenlight.dekutyavin.net/internal/data"
	"greenlight.dekutyavin.net/internal/validator"
)

const (
	msgAuthSkipped   = "authentication skipped"
	msgAuthFailed    = "authentication failed"
	msgAuthSucceeded = "authentication succeeded"

	skipReasonNoHeader        = "no authorization header"
	failReasonMalformedHeader = "malformed authorization header"
	failReasonInvalidToken    = "invalid token"
	failReasonTokenNotFound   = "token not found"
	failReasonServerError     = "server error"

	logFieldReason = "reason"
	logFieldUser   = "user_id"
)

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			pv := recover()

			if pv != nil {
				w.Header().Set("Connection", "close")

				app.serverErrorResponse(w, r, fmt.Errorf("%v", pv))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *Application) rateLimit(next http.Handler) http.Handler {
	if !app.Config.Limiter.Enabled {
		return next
	}

	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := realip.FromRequest(r)

		mu.Lock()

		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(app.Config.Limiter.RPS), app.Config.Limiter.Burst)}
		}

		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			app.rateLimitExceedResponse(w, r)
			return
		}

		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

func (app *Application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")

		logAuthEvent := func(msg string, attrs ...any) {
			attrs = append([]any{"method", r.Method, "uri", r.URL.RequestURI()}, attrs...)
			app.Logger.Debug(msg, attrs...)
		}

		if authorizationHeader == "" {
			logAuthEvent(msgAuthSkipped, logFieldReason, skipReasonNoHeader)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			logAuthEvent(msgAuthFailed, logFieldReason, failReasonMalformedHeader)
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		v := validator.New()

		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			logAuthEvent(msgAuthFailed, logFieldReason, failReasonInvalidToken)
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := app.Models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				logAuthEvent(msgAuthFailed, logFieldReason, failReasonTokenNotFound)
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		r = app.contextSetAuthenticatedUser(r, user)
		logAuthEvent(msgAuthSucceeded, logFieldUser, user.ID)
		next.ServeHTTP(w, r)
	})
}
