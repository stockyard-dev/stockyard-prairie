package store
import ("database/sql";"encoding/json";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Board struct{ID string `json:"id"`;Name string `json:"name"`;Description string `json:"description,omitempty"`;Columns []string `json:"columns"`;CreatedAt string `json:"created_at"`;CardCount int `json:"card_count"`}
type Card struct{ID string `json:"id"`;BoardID string `json:"board_id"`;Column string `json:"column"`;Title string `json:"title"`;Description string `json:"description,omitempty"`;Assignee string `json:"assignee,omitempty"`;Labels string `json:"labels,omitempty"`;Position int `json:"position"`;CreatedAt string `json:"created_at"`;UpdatedAt string `json:"updated_at"`}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"prairie.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
for _,q:=range[]string{
`CREATE TABLE IF NOT EXISTS boards(id TEXT PRIMARY KEY,name TEXT NOT NULL,description TEXT DEFAULT '',columns_json TEXT DEFAULT '["Backlog","Todo","In Progress","Done"]',created_at TEXT DEFAULT(datetime('now')))`,
`CREATE TABLE IF NOT EXISTS cards(id TEXT PRIMARY KEY,board_id TEXT NOT NULL,col TEXT DEFAULT 'Backlog',title TEXT NOT NULL,description TEXT DEFAULT '',assignee TEXT DEFAULT '',labels TEXT DEFAULT '',position INTEGER DEFAULT 0,created_at TEXT DEFAULT(datetime('now')),updated_at TEXT DEFAULT(datetime('now')))`,
`CREATE INDEX IF NOT EXISTS idx_cards_board ON cards(board_id)`,
}{if _,err:=db.Exec(q);err!=nil{return nil,fmt.Errorf("migrate: %w",err)}};return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)CreateBoard(b *Board)error{b.ID=genID();b.CreatedAt=now();if b.Columns==nil{b.Columns=[]string{"Backlog","Todo","In Progress","Done"}}
cj,_:=json.Marshal(b.Columns);_,err:=d.db.Exec(`INSERT INTO boards VALUES(?,?,?,?,?)`,b.ID,b.Name,b.Description,string(cj),b.CreatedAt);return err}
func(d *DB)GetBoard(id string)*Board{var b Board;var cj string;if d.db.QueryRow(`SELECT id,name,description,columns_json,created_at FROM boards WHERE id=?`,id).Scan(&b.ID,&b.Name,&b.Description,&cj,&b.CreatedAt)!=nil{return nil}
json.Unmarshal([]byte(cj),&b.Columns);d.db.QueryRow(`SELECT COUNT(*) FROM cards WHERE board_id=?`,b.ID).Scan(&b.CardCount);return &b}
func(d *DB)ListBoards()[]Board{rows,_:=d.db.Query(`SELECT id,name,description,columns_json,created_at FROM boards ORDER BY name`);if rows==nil{return nil};defer rows.Close()
var o []Board;for rows.Next(){var b Board;var cj string;rows.Scan(&b.ID,&b.Name,&b.Description,&cj,&b.CreatedAt);json.Unmarshal([]byte(cj),&b.Columns);d.db.QueryRow(`SELECT COUNT(*) FROM cards WHERE board_id=?`,b.ID).Scan(&b.CardCount);o=append(o,b)};return o}
func(d *DB)DeleteBoard(id string)error{d.db.Exec(`DELETE FROM cards WHERE board_id=?`,id);_,err:=d.db.Exec(`DELETE FROM boards WHERE id=?`,id);return err}
func(d *DB)CreateCard(c *Card)error{c.ID=genID();c.CreatedAt=now();c.UpdatedAt=c.CreatedAt;if c.Column==""{c.Column="Backlog"}
_,err:=d.db.Exec(`INSERT INTO cards(id,board_id,col,title,description,assignee,labels,position,created_at,updated_at)VALUES(?,?,?,?,?,?,?,?,?,?)`,c.ID,c.BoardID,c.Column,c.Title,c.Description,c.Assignee,c.Labels,c.Position,c.CreatedAt,c.UpdatedAt);return err}
func(d *DB)GetCard(id string)*Card{var c Card;if d.db.QueryRow(`SELECT id,board_id,col,title,description,assignee,labels,position,created_at,updated_at FROM cards WHERE id=?`,id).Scan(&c.ID,&c.BoardID,&c.Column,&c.Title,&c.Description,&c.Assignee,&c.Labels,&c.Position,&c.CreatedAt,&c.UpdatedAt)!=nil{return nil};return &c}
func(d *DB)ListCards(boardID string)[]Card{rows,_:=d.db.Query(`SELECT id,board_id,col,title,description,assignee,labels,position,created_at,updated_at FROM cards WHERE board_id=? ORDER BY position,created_at`,boardID);if rows==nil{return nil};defer rows.Close()
var o []Card;for rows.Next(){var c Card;rows.Scan(&c.ID,&c.BoardID,&c.Column,&c.Title,&c.Description,&c.Assignee,&c.Labels,&c.Position,&c.CreatedAt,&c.UpdatedAt);o=append(o,c)};return o}
func(d *DB)MoveCard(id,column string,position int)error{_,err:=d.db.Exec(`UPDATE cards SET col=?,position=?,updated_at=? WHERE id=?`,column,position,now(),id);return err}
func(d *DB)UpdateCard(id string,c *Card)error{_,err:=d.db.Exec(`UPDATE cards SET title=?,description=?,assignee=?,labels=?,updated_at=? WHERE id=?`,c.Title,c.Description,c.Assignee,c.Labels,now(),id);return err}
func(d *DB)DeleteCard(id string)error{_,err:=d.db.Exec(`DELETE FROM cards WHERE id=?`,id);return err}
type Stats struct{Boards int `json:"boards"`;Cards int `json:"cards"`}
func(d *DB)Stats()Stats{var s Stats;d.db.QueryRow(`SELECT COUNT(*) FROM boards`).Scan(&s.Boards);d.db.QueryRow(`SELECT COUNT(*) FROM cards`).Scan(&s.Cards);return s}
