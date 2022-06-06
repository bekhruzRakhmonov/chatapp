package chat

// import (
// 	"flag"
// 	"log"
// 	"net/http"

// )

// var addr = flag.String("addr", ":8080", "http service address")

// func serveHome(w http.ResponseWriter, r *http.Request) {
// 	log.Println(r.URL)
// 	if r.URL.Path != "/" {
// 		http.Error(w, "Not found", http.StatusNotFound)
// 		return
// 	}
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	http.ServeFile(w, r, "home.html")
// }

// func main() {
// 	flag.Parse()
// 	hub := newHub()
// 	go hub.run()
// 	http.HandleFunc("/", serveHome)
// 	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
// 		log.Println("Method:",r.Method)
// 		username :=  r.URL.Query().Get("username")
// 		log.Println("Username:",username)
// 		serveWs(hub, w, r)
// 	})

// 	err := http.ListenAndServe(*addr, nil)
// 	if err != nil {
// 		log.Fatal("ListenAndServe: ", err)
// 	}
// }

// C:\Program Files\Go\src\net\http