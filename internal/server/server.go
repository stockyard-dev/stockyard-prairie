package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/stockyard-dev/stockyard-prairie/internal/store"
)

type Server struct {
	db      *store.DB
	mux     *http.ServeMux
	limits  Limits
	dataDir string
	pCfg    map[string]json.RawMessage
}

func New(db *store.DB, limits Limits, dataDir string) *Server {
	s := &Server{
		db:      db,
		mux:     http.NewServeMux(),
		limits:  limits,
		dataDir: dataDir,
	}
	s.loadPersonalConfig()

	// Boards
	s.mux.HandleFunc("GET /api/boards", s.listBoards)
	s.mux.HandleFunc("POST /api/boards", s.createBoard)
	s.mux.HandleFunc("GET /api/boards/{id}", s.getBoard)
	s.mux.HandleFunc("PUT /api/boards/{id}", s.updateBoard) // NEW — original had no update endpoint
	s.mux.HandleFunc("DELETE /api/boards/{id}", s.deleteBoard)
	s.mux.HandleFunc("GET /api/boards/{id}/cards", s.listCards)

	// Cards
	s.mux.HandleFunc("POST /api/cards", s.createCard)
	s.mux.HandleFunc("GET /api/cards/{id}", s.getCard)
	s.mux.HandleFunc("PUT /api/cards/{id}", s.updateCard)
	s.mux.HandleFunc("POST /api/cards/{id}/move", s.moveCard)
	s.mux.HandleFunc("DELETE /api/cards/{id}", s.deleteCard)

	// Stats / health
	s.mux.HandleFunc("GET /api/stats", s.stats)
	s.mux.HandleFunc("GET /api/health", s.health)

	// Personalization
	s.mux.HandleFunc("GET /api/config", s.configHandler)

	// Extras (works for both boards and cards via {resource} path param)
	s.mux.HandleFunc("GET /api/extras/{resource}", s.listExtras)
	s.mux.HandleFunc("GET /api/extras/{resource}/{id}", s.getExtras)
	s.mux.HandleFunc("PUT /api/extras/{resource}/{id}", s.putExtras)

	// Tier
	s.mux.HandleFunc("GET /api/tier", func(w http.ResponseWriter, r *http.Request) {
		wj(w, 200, map[string]any{
			"tier":        s.limits.Tier,
			"upgrade_url": "https://stockyard.dev/prairie/",
		})
	})

	// Dashboard
	s.mux.HandleFunc("GET /ui", s.dashboard)
	s.mux.HandleFunc("GET /ui/", s.dashboard)
	s.mux.HandleFunc("GET /", s.root)

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// ─── helpers ──────────────────────────────────────────────────────

func wj(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func we(w http.ResponseWriter, code int, msg string) {
	wj(w, code, map[string]string{"error": msg})
}

func oe[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}

func (s *Server) root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/ui", 302)
}

// ─── personalization ──────────────────────────────────────────────

func (s *Server) loadPersonalConfig() {
	path := filepath.Join(s.dataDir, "config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var cfg map[string]json.RawMessage
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Printf("prairie: warning: could not parse config.json: %v", err)
		return
	}
	s.pCfg = cfg
	log.Printf("prairie: loaded personalization from %s", path)
}

func (s *Server) configHandler(w http.ResponseWriter, r *http.Request) {
	if s.pCfg == nil {
		wj(w, 200, map[string]any{})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.pCfg)
}

// ─── extras ───────────────────────────────────────────────────────

func (s *Server) listExtras(w http.ResponseWriter, r *http.Request) {
	resource := r.PathValue("resource")
	all := s.db.AllExtras(resource)
	out := make(map[string]json.RawMessage, len(all))
	for id, data := range all {
		out[id] = json.RawMessage(data)
	}
	wj(w, 200, out)
}

func (s *Server) getExtras(w http.ResponseWriter, r *http.Request) {
	resource := r.PathValue("resource")
	id := r.PathValue("id")
	data := s.db.GetExtras(resource, id)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))
}

func (s *Server) putExtras(w http.ResponseWriter, r *http.Request) {
	resource := r.PathValue("resource")
	id := r.PathValue("id")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		we(w, 400, "read body")
		return
	}
	var probe map[string]any
	if err := json.Unmarshal(body, &probe); err != nil {
		we(w, 400, "invalid json")
		return
	}
	if err := s.db.SetExtras(resource, id, string(body)); err != nil {
		we(w, 500, "save failed")
		return
	}
	wj(w, 200, map[string]string{"ok": "saved"})
}

// ─── boards ───────────────────────────────────────────────────────

func (s *Server) listBoards(w http.ResponseWriter, r *http.Request) {
	wj(w, 200, map[string]any{"boards": oe(s.db.ListBoards())})
}

func (s *Server) createBoard(w http.ResponseWriter, r *http.Request) {
	var b store.Board
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		we(w, 400, "invalid json")
		return
	}
	if b.Name == "" {
		we(w, 400, "name required")
		return
	}
	if err := s.db.CreateBoard(&b); err != nil {
		we(w, 500, "create failed")
		return
	}
	wj(w, 201, s.db.GetBoard(b.ID))
}

func (s *Server) getBoard(w http.ResponseWriter, r *http.Request) {
	b := s.db.GetBoard(r.PathValue("id"))
	if b == nil {
		we(w, 404, "not found")
		return
	}
	wj(w, 200, b)
}

// updateBoard accepts a partial board payload. The original prairie had
// no update endpoint at all — boards were create+delete-only. Now boards
// can be renamed, re-described, and have their column list rewritten.
func (s *Server) updateBoard(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ex := s.db.GetBoard(id)
	if ex == nil {
		we(w, 404, "not found")
		return
	}

	var raw map[string]json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		we(w, 400, "invalid json")
		return
	}

	patch := *ex
	if v, ok := raw["name"]; ok {
		var s string
		json.Unmarshal(v, &s)
		if s != "" {
			patch.Name = s
		}
	}
	if v, ok := raw["description"]; ok {
		json.Unmarshal(v, &patch.Description)
	}
	if v, ok := raw["columns"]; ok {
		var cols []string
		if err := json.Unmarshal(v, &cols); err == nil && len(cols) > 0 {
			patch.Columns = cols
		}
	}

	if err := s.db.UpdateBoard(id, &patch); err != nil {
		we(w, 500, "update failed")
		return
	}
	wj(w, 200, s.db.GetBoard(id))
}

func (s *Server) deleteBoard(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	// Cascade: delete all card extras for this board's cards
	cardIDs := s.db.BoardCardIDs(id)
	for _, cid := range cardIDs {
		s.db.DeleteExtras("cards", cid)
	}
	if err := s.db.DeleteBoard(id); err != nil {
		we(w, 500, err.Error())
		return
	}
	s.db.DeleteExtras("boards", id)
	wj(w, 200, map[string]string{"deleted": "ok"})
}

// ─── cards ────────────────────────────────────────────────────────

func (s *Server) listCards(w http.ResponseWriter, r *http.Request) {
	wj(w, 200, map[string]any{"cards": oe(s.db.ListCards(r.PathValue("id")))})
}

func (s *Server) createCard(w http.ResponseWriter, r *http.Request) {
	if s.limits.MaxItems > 0 {
		var n int
		s.db.Stats() // (Stats does its own counts; we don't need this)
		_ = n
	}
	var c store.Card
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		we(w, 400, "invalid json")
		return
	}
	if c.Title == "" || c.BoardID == "" {
		we(w, 400, "title and board_id required")
		return
	}
	if err := s.db.CreateCard(&c); err != nil {
		we(w, 500, "create failed")
		return
	}
	wj(w, 201, s.db.GetCard(c.ID))
}

func (s *Server) getCard(w http.ResponseWriter, r *http.Request) {
	c := s.db.GetCard(r.PathValue("id"))
	if c == nil {
		we(w, 404, "not found")
		return
	}
	wj(w, 200, c)
}

// updateCard accepts a partial card payload. The original updateCard
// only preserved Title and silently nuked description, assignee, labels,
// and column on every partial PUT — same severity bug as ledger.
func (s *Server) updateCard(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ex := s.db.GetCard(id)
	if ex == nil {
		we(w, 404, "not found")
		return
	}

	var raw map[string]json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		we(w, 400, "invalid json")
		return
	}

	patch := *ex
	if v, ok := raw["title"]; ok {
		var s string
		json.Unmarshal(v, &s)
		if s != "" {
			patch.Title = s
		}
	}
	if v, ok := raw["description"]; ok {
		json.Unmarshal(v, &patch.Description)
	}
	if v, ok := raw["assignee"]; ok {
		json.Unmarshal(v, &patch.Assignee)
	}
	if v, ok := raw["labels"]; ok {
		json.Unmarshal(v, &patch.Labels)
	}
	if v, ok := raw["column"]; ok {
		var s string
		json.Unmarshal(v, &s)
		if s != "" {
			patch.Column = s
		}
	}

	if err := s.db.UpdateCard(id, &patch); err != nil {
		we(w, 500, "update failed")
		return
	}
	wj(w, 200, s.db.GetCard(id))
}

func (s *Server) moveCard(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Column   string `json:"column"`
		Position int    `json:"position"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		we(w, 400, "invalid json")
		return
	}
	if req.Column == "" {
		we(w, 400, "column required")
		return
	}
	id := r.PathValue("id")
	if err := s.db.MoveCard(id, req.Column, req.Position); err != nil {
		we(w, 500, "move failed")
		return
	}
	wj(w, 200, s.db.GetCard(id))
}

func (s *Server) deleteCard(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	s.db.DeleteCard(id)
	s.db.DeleteExtras("cards", id)
	wj(w, 200, map[string]string{"deleted": "ok"})
}

func (s *Server) stats(w http.ResponseWriter, r *http.Request) {
	wj(w, 200, s.db.Stats())
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	st := s.db.Stats()
	wj(w, 200, map[string]any{
		"status":  "ok",
		"service": "prairie",
		"boards":  st.Boards,
		"cards":   st.Cards,
	})
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
