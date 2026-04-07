package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct{ db *sql.DB }

// Board is a kanban board with a list of column names. The columns array
// is JSON-encoded in the columns_json column.
type Board struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Columns     []string `json:"columns"`
	CreatedAt   string   `json:"created_at"`
	CardCount   int      `json:"card_count"`
}

// Card is a kanban card belonging to a board. Position determines vertical
// order within a column. Labels is a comma-separated string for simplicity.
type Card struct {
	ID          string `json:"id"`
	BoardID     string `json:"board_id"`
	Column      string `json:"column"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Assignee    string `json:"assignee,omitempty"`
	Labels      string `json:"labels,omitempty"`
	Position    int    `json:"position"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func defaultColumns() []string {
	return []string{"Backlog", "Todo", "In Progress", "Done"}
}

func Open(d string) (*DB, error) {
	if err := os.MkdirAll(d, 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", filepath.Join(d, "prairie.db")+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}
	for _, q := range []string{
		`CREATE TABLE IF NOT EXISTS boards(
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT DEFAULT '',
			columns_json TEXT DEFAULT '["Backlog","Todo","In Progress","Done"]',
			created_at TEXT DEFAULT(datetime('now'))
		)`,
		`CREATE TABLE IF NOT EXISTS cards(
			id TEXT PRIMARY KEY,
			board_id TEXT NOT NULL,
			col TEXT DEFAULT 'Backlog',
			title TEXT NOT NULL,
			description TEXT DEFAULT '',
			assignee TEXT DEFAULT '',
			labels TEXT DEFAULT '',
			position INTEGER DEFAULT 0,
			created_at TEXT DEFAULT(datetime('now')),
			updated_at TEXT DEFAULT(datetime('now'))
		)`,
		`CREATE INDEX IF NOT EXISTS idx_cards_board ON cards(board_id)`,
		`CREATE INDEX IF NOT EXISTS idx_cards_col ON cards(col)`,
		`CREATE TABLE IF NOT EXISTS extras(
			resource TEXT NOT NULL,
			record_id TEXT NOT NULL,
			data TEXT NOT NULL DEFAULT '{}',
			PRIMARY KEY(resource, record_id)
		)`,
	} {
		if _, err := db.Exec(q); err != nil {
			return nil, fmt.Errorf("migrate: %w", err)
		}
	}
	return &DB{db: db}, nil
}

func (d *DB) Close() error { return d.db.Close() }

func genID() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string   { return time.Now().UTC().Format(time.RFC3339) }

// ─── Boards ───────────────────────────────────────────────────────

func (d *DB) CreateBoard(b *Board) error {
	b.ID = genID()
	b.CreatedAt = now()
	if b.Columns == nil || len(b.Columns) == 0 {
		b.Columns = defaultColumns()
	}
	cj, _ := json.Marshal(b.Columns)
	_, err := d.db.Exec(
		`INSERT INTO boards(id, name, description, columns_json, created_at) VALUES(?, ?, ?, ?, ?)`,
		b.ID, b.Name, b.Description, string(cj), b.CreatedAt,
	)
	return err
}

func (d *DB) GetBoard(id string) *Board {
	var b Board
	var cj string
	err := d.db.QueryRow(
		`SELECT id, name, description, columns_json, created_at FROM boards WHERE id=?`, id,
	).Scan(&b.ID, &b.Name, &b.Description, &cj, &b.CreatedAt)
	if err != nil {
		return nil
	}
	json.Unmarshal([]byte(cj), &b.Columns)
	if b.Columns == nil {
		b.Columns = defaultColumns()
	}
	d.db.QueryRow(`SELECT COUNT(*) FROM cards WHERE board_id=?`, b.ID).Scan(&b.CardCount)
	return &b
}

func (d *DB) ListBoards() []Board {
	rows, _ := d.db.Query(`SELECT id, name, description, columns_json, created_at FROM boards ORDER BY name ASC`)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Board
	for rows.Next() {
		var b Board
		var cj string
		rows.Scan(&b.ID, &b.Name, &b.Description, &cj, &b.CreatedAt)
		json.Unmarshal([]byte(cj), &b.Columns)
		if b.Columns == nil {
			b.Columns = defaultColumns()
		}
		d.db.QueryRow(`SELECT COUNT(*) FROM cards WHERE board_id=?`, b.ID).Scan(&b.CardCount)
		o = append(o, b)
	}
	return o
}

// UpdateBoard rewrites name, description, and columns. The original
// implementation had no UpdateBoard at all — boards were create+delete-only.
func (d *DB) UpdateBoard(id string, b *Board) error {
	cj, _ := json.Marshal(b.Columns)
	_, err := d.db.Exec(
		`UPDATE boards SET name=?, description=?, columns_json=? WHERE id=?`,
		b.Name, b.Description, string(cj), id,
	)
	return err
}

// DeleteBoard removes the board and all its cards. Card extras are also
// deleted (caller-driven via DeleteExtras for each card or wholesale).
func (d *DB) DeleteBoard(id string) error {
	d.db.Exec(`DELETE FROM cards WHERE board_id=?`, id)
	_, err := d.db.Exec(`DELETE FROM boards WHERE id=?`, id)
	return err
}

func (d *DB) BoardCardIDs(boardID string) []string {
	rows, _ := d.db.Query(`SELECT id FROM cards WHERE board_id=?`, boardID)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var s string
		rows.Scan(&s)
		out = append(out, s)
	}
	return out
}

// ─── Cards ────────────────────────────────────────────────────────

func (d *DB) CreateCard(c *Card) error {
	c.ID = genID()
	c.CreatedAt = now()
	c.UpdatedAt = c.CreatedAt
	if c.Column == "" {
		c.Column = "Backlog"
	}
	_, err := d.db.Exec(
		`INSERT INTO cards(id, board_id, col, title, description, assignee, labels, position, created_at, updated_at)
		 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.ID, c.BoardID, c.Column, c.Title, c.Description, c.Assignee, c.Labels, c.Position, c.CreatedAt, c.UpdatedAt,
	)
	return err
}

func (d *DB) GetCard(id string) *Card {
	var c Card
	err := d.db.QueryRow(
		`SELECT id, board_id, col, title, description, assignee, labels, position, created_at, updated_at
		 FROM cards WHERE id=?`,
		id,
	).Scan(&c.ID, &c.BoardID, &c.Column, &c.Title, &c.Description, &c.Assignee, &c.Labels, &c.Position, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil
	}
	return &c
}

func (d *DB) ListCards(boardID string) []Card {
	rows, _ := d.db.Query(
		`SELECT id, board_id, col, title, description, assignee, labels, position, created_at, updated_at
		 FROM cards WHERE board_id=?
		 ORDER BY position ASC, created_at ASC`,
		boardID,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Card
	for rows.Next() {
		var c Card
		rows.Scan(&c.ID, &c.BoardID, &c.Column, &c.Title, &c.Description, &c.Assignee, &c.Labels, &c.Position, &c.CreatedAt, &c.UpdatedAt)
		o = append(o, c)
	}
	return o
}

// MoveCard updates only the column and position (and bumps updated_at).
// Used by drag-and-drop without touching other fields.
func (d *DB) MoveCard(id, column string, position int) error {
	_, err := d.db.Exec(
		`UPDATE cards SET col=?, position=?, updated_at=? WHERE id=?`,
		column, position, now(), id,
	)
	return err
}

// UpdateCard rewrites the editable card fields. Caller is responsible for
// preserving omitted fields by passing the merged record.
func (d *DB) UpdateCard(id string, c *Card) error {
	_, err := d.db.Exec(
		`UPDATE cards SET title=?, description=?, assignee=?, labels=?, col=?, updated_at=? WHERE id=?`,
		c.Title, c.Description, c.Assignee, c.Labels, c.Column, now(), id,
	)
	return err
}

func (d *DB) DeleteCard(id string) error {
	_, err := d.db.Exec(`DELETE FROM cards WHERE id=?`, id)
	return err
}

// ─── Stats ────────────────────────────────────────────────────────

type Stats struct {
	Boards     int            `json:"boards"`
	Cards      int            `json:"cards"`
	ByColumn   map[string]int `json:"by_column"`
	ByAssignee map[string]int `json:"by_assignee"`
}

func (d *DB) Stats() Stats {
	var s Stats
	d.db.QueryRow(`SELECT COUNT(*) FROM boards`).Scan(&s.Boards)
	d.db.QueryRow(`SELECT COUNT(*) FROM cards`).Scan(&s.Cards)

	s.ByColumn = map[string]int{}
	if rows, _ := d.db.Query(`SELECT col, COUNT(*) FROM cards GROUP BY col`); rows != nil {
		defer rows.Close()
		for rows.Next() {
			var k string
			var v int
			rows.Scan(&k, &v)
			s.ByColumn[k] = v
		}
	}

	s.ByAssignee = map[string]int{}
	if rows, _ := d.db.Query(`SELECT assignee, COUNT(*) FROM cards WHERE assignee != '' GROUP BY assignee`); rows != nil {
		defer rows.Close()
		for rows.Next() {
			var k string
			var v int
			rows.Scan(&k, &v)
			s.ByAssignee[k] = v
		}
	}

	return s
}

// ─── Extras: generic key-value storage for personalization custom fields ───

func (d *DB) GetExtras(resource, recordID string) string {
	var data string
	err := d.db.QueryRow(
		`SELECT data FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	).Scan(&data)
	if err != nil || data == "" {
		return "{}"
	}
	return data
}

func (d *DB) SetExtras(resource, recordID, data string) error {
	if data == "" {
		data = "{}"
	}
	_, err := d.db.Exec(
		`INSERT INTO extras(resource, record_id, data) VALUES(?, ?, ?)
		 ON CONFLICT(resource, record_id) DO UPDATE SET data=excluded.data`,
		resource, recordID, data,
	)
	return err
}

func (d *DB) DeleteExtras(resource, recordID string) error {
	_, err := d.db.Exec(
		`DELETE FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	)
	return err
}

func (d *DB) AllExtras(resource string) map[string]string {
	out := make(map[string]string)
	rows, _ := d.db.Query(
		`SELECT record_id, data FROM extras WHERE resource=?`,
		resource,
	)
	if rows == nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id, data string
		rows.Scan(&id, &data)
		out[id] = data
	}
	return out
}
