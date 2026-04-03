package server
import ("encoding/json";"log";"net/http";"strconv";"github.com/stockyard-dev/stockyard-prairie/internal/store")
type Server struct{db *store.DB;mux *http.ServeMux;limits Limits}
func New(db *store.DB,limits Limits)*Server{s:=&Server{db:db,mux:http.NewServeMux(),limits:limits}
s.mux.HandleFunc("GET /api/boards",s.listBoards);s.mux.HandleFunc("POST /api/boards",s.createBoard);s.mux.HandleFunc("GET /api/boards/{id}",s.getBoard);s.mux.HandleFunc("DELETE /api/boards/{id}",s.deleteBoard)
s.mux.HandleFunc("GET /api/boards/{id}/cards",s.listCards);s.mux.HandleFunc("POST /api/cards",s.createCard);s.mux.HandleFunc("PUT /api/cards/{id}",s.updateCard);s.mux.HandleFunc("POST /api/cards/{id}/move",s.moveCard);s.mux.HandleFunc("DELETE /api/cards/{id}",s.deleteCard)
s.mux.HandleFunc("GET /api/stats",s.stats);s.mux.HandleFunc("GET /api/health",s.health)
s.mux.HandleFunc("GET /ui",s.dashboard);s.mux.HandleFunc("GET /ui/",s.dashboard);s.mux.HandleFunc("GET /",s.root);
s.mux.HandleFunc("GET /api/tier",func(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"tier":s.limits.Tier,"upgrade_url":"https://stockyard.dev/prairie/"})})
return s}
func(s *Server)ServeHTTP(w http.ResponseWriter,r *http.Request){s.mux.ServeHTTP(w,r)}
func wj(w http.ResponseWriter,c int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(c);json.NewEncoder(w).Encode(v)}
func we(w http.ResponseWriter,c int,m string){wj(w,c,map[string]string{"error":m})}
func(s *Server)root(w http.ResponseWriter,r *http.Request){if r.URL.Path!="/"{http.NotFound(w,r);return};http.Redirect(w,r,"/ui",302)}
func(s *Server)listBoards(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"boards":oe(s.db.ListBoards())})}
func(s *Server)createBoard(w http.ResponseWriter,r *http.Request){var b store.Board;json.NewDecoder(r.Body).Decode(&b);if b.Name==""{we(w,400,"name required");return};s.db.CreateBoard(&b);wj(w,201,s.db.GetBoard(b.ID))}
func(s *Server)getBoard(w http.ResponseWriter,r *http.Request){b:=s.db.GetBoard(r.PathValue("id"));if b==nil{we(w,404,"not found");return};wj(w,200,b)}
func(s *Server)deleteBoard(w http.ResponseWriter,r *http.Request){s.db.DeleteBoard(r.PathValue("id"));wj(w,200,map[string]string{"deleted":"ok"})}
func(s *Server)listCards(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"cards":oe(s.db.ListCards(r.PathValue("id")))})}
func(s *Server)createCard(w http.ResponseWriter,r *http.Request){var c store.Card;json.NewDecoder(r.Body).Decode(&c);if c.Title==""||c.BoardID==""{we(w,400,"title and board_id required");return};s.db.CreateCard(&c);wj(w,201,s.db.GetCard(c.ID))}
func(s *Server)updateCard(w http.ResponseWriter,r *http.Request){id:=r.PathValue("id");ex:=s.db.GetCard(id);if ex==nil{we(w,404,"not found");return};var c store.Card;json.NewDecoder(r.Body).Decode(&c);if c.Title==""{c.Title=ex.Title};s.db.UpdateCard(id,&c);wj(w,200,s.db.GetCard(id))}
func(s *Server)moveCard(w http.ResponseWriter,r *http.Request){var req struct{Column string `json:"column"`;Position int `json:"position"`};json.NewDecoder(r.Body).Decode(&req);if req.Column==""{we(w,400,"column required");return};s.db.MoveCard(r.PathValue("id"),req.Column,req.Position);wj(w,200,s.db.GetCard(r.PathValue("id")))}
func(s *Server)deleteCard(w http.ResponseWriter,r *http.Request){s.db.DeleteCard(r.PathValue("id"));wj(w,200,map[string]string{"deleted":"ok"})}
func(s *Server)stats(w http.ResponseWriter,r *http.Request){wj(w,200,s.db.Stats())}
func(s *Server)health(w http.ResponseWriter,r *http.Request){st:=s.db.Stats();wj(w,200,map[string]any{"status":"ok","service":"prairie","boards":st.Boards,"cards":st.Cards})}
func oe[T any](s []T)[]T{if s==nil{return[]T{}};return s}
func init(){log.SetFlags(log.LstdFlags|log.Lshortfile);_=strconv.Atoi}
