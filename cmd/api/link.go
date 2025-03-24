package main

import (
	"encoding/json"
	"net/http"
	// "strings"
	"time"

	"github.com/devaartana/e01-oprec-rpl/internal/store"
	"github.com/go-chi/chi/v5"
)

func (app *application) SlugHandler(w http.ResponseWriter, r *http.Request) {

	slug := chi.URLParam(r, "slug")

	app.logger.Infow("slug", "slug", slug)
	link, err := app.store.Links.GetBySlug(r.Context(), slug)
	if err != nil {
		if err == store.ErrNotFound {	
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return 
	}

	if link.Expired_date.Before(time.Now()) {
		http.Error(w, "Link is expired", http.StatusGone)
		return
	}

	http.Redirect(w, r, link.OriginalUrl, http.StatusMovedPermanently)
}

type CreateLinkPayload struct {
	Slug        string `json:"slug"`
	OriginalUrl string `json:"original_url"`
}

func (app *application) CreateLinkHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateLinkPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	link := &store.Link{
		Slug:         payload.Slug,
		OriginalUrl:  payload.OriginalUrl,
		Created_at:   time.Now(),
		Expired_date: time.Now().Add(time.Hour * 24 * 30),
	}

	user := r.Context().Value(userCtx).(*store.User)

	if err := app.store.Links.Create(r.Context(), user.Email, link); err != nil {
		if err == store.ErrDuplicateSlug {
			http.Error(w, "Slug is already exist", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Link is created"))
}

func (app *application) GetAllLinksHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userCtx).(*store.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	links, err := app.store.Links.GetAll(r.Context(), user.Email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(links)
}

func (app *application) DeleteLinkHandler(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	link, err := app.store.Links.GetBySlug(r.Context(), slug)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if link.Expired_date.Before(time.Now()) {
		http.Error(w, "Link is expired", http.StatusGone)
		return
	}

	user := r.Context().Value(userCtx).(*store.User)
	if err := app.store.Links.DeleteBySlug(r.Context(), user.Email, slug); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Link is deleted"))
}

type UpdateLinkPayload struct {
	Slug        string `json:"slug"`
	OriginalUrl string `json:"original_url"`
}

func (app *application) UpdateLinkHandler(w http.ResponseWriter, r *http.Request) {
	
	var payload UpdateLinkPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	link, err := app.store.Links.GetBySlug(r.Context(), payload.Slug)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if link.Expired_date.Before(time.Now()) {
		http.Error(w, "Link is expired", http.StatusGone)
		return
	}

	link.Slug = payload.Slug
	link.OriginalUrl = payload.OriginalUrl

	user := r.Context().Value(userCtx).(*store.User)
	if err := app.store.Links.UpdateBySlug(r.Context(), user.Email, link); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Link is updated"))
}

func (app *application) RefreshExpiredDateHandler(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	link, err := app.store.Links.GetBySlug(r.Context(), slug)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	link.Expired_date = time.Now().Add(time.Hour * 24 * 30)

	user := r.Context().Value(userCtx).(*store.User)
	if err := app.store.Links.UpdateBySlug(r.Context(), user.Email, link); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Link is updated"))
}
