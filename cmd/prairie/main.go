package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-prairie/internal/server";"github.com/stockyard-dev/stockyard-prairie/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9110"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./prairie-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("prairie: %v",err)};defer db.Close();srv:=server.New(db,server.DefaultLimits())
fmt.Printf("\n  Prairie — Self-hosted kanban project board\n  ─────────────────────────────────\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n  Data:       %s\n  ─────────────────────────────────\n  Questions? hello@stockyard.dev\n\n",port,port,dataDir)
log.Printf("prairie: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
