package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/torfjor/go-vipps/login"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"time"
)

func main() {
	provider, err := login.NewProvider(context.Background(), &login.ProviderConfig{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		IssuerURL:    login.IssuerURLTesting,
		RedirectURL:  "http://localhost:3000/redirect",
		Scopes: []string{
			login.ScopeName,
			login.ScopeEmail,
			login.ScopeBirthDate,
			login.ScopeAdress,
			login.ScopePhoneNumber,
		},
	})
	if err != nil {
		fmt.Printf("NewProvider: %v", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, provider.AuthCodeURL("0xdeadbeef"), http.StatusFound)
	})
	mux.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		claims, err := provider.ExchangeCodeForClaims(r.Context(), r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bytes, err := json.MarshalIndent(claims, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(bytes)
	})

	s := http.Server{
		Addr:         ":3000",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      mux,
	}

	var g errgroup.Group
	g.Go(func() error {
		return s.ListenAndServe()
	})
	fmt.Printf("server listening")
	err = g.Wait()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
