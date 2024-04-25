package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"test/db/database"
)

// create new server
func NewTaskServer(d database.DB) *taskServer {
	return &taskServer{
		db: d,
	}
}

// internal struct used in Server interfaceß
type taskServer struct {
	server *http.ServeMux
	db     database.DB
}

// Initialise the mux server and task server. Mux server runs on localhost:8080
// adds path handlers to the mux server, adds inputed database to the task server
func Start(db database.DB) {
	server := http.NewServeMux()
	ts := NewTaskServer(db)
	server.HandleFunc("GET /", ts.homeHandler)
	server.HandleFunc("GET /read/{id}", ts.readHandler)
	server.HandleFunc("POST /create", ts.createHandler)
	server.HandleFunc("POST /update", ts.updateHandler)
	server.HandleFunc("DELETE /delete/{id}", ts.deleteHandler)
	server.HandleFunc("GET /list", ts.getHandler)

	go db.Listen()

	fmt.Println("Server is running on port 8080...")
	err := http.ListenAndServe("localhost:8080", server)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}

}

func (ts *taskServer) homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Something went wrong.....or very right")
}

func (ts *taskServer) readHandler(w http.ResponseWriter, r *http.Request) {
	if ok := ts.db.Exists(r.PathValue("id")); !ok {
		http.Error(w, "Data entry does not exist", 404)
		return
	}

	msg := database.Message{
		Action:   "read",
		Document: database.Doc{ID: r.PathValue("id")},
		Resp:     make(chan database.Response),
	}

	resp := ts.db.Send(msg)
	if !resp.Result {
		http.Error(w, "Error reading from database", http.StatusInternalServerError)
		return
	}

	jsonBody, err := json.Marshal(resp.Documents[0])
	if err != nil {
		http.Error(w, "Error reading database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBody)

}

func (ts *taskServer) createHandler(w http.ResponseWriter, r *http.Request) {
	var body database.Doc
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if ok := ts.db.Exists(body.ID); ok {
		http.Error(w, "Data entry already exists", 409)
		return
	}

	msg := database.Message{
		Action:   "create",
		Document: body,
		Resp:     make(chan database.Response),
	}

	resp := ts.db.Send(msg)
	fmt.Println("I have found the result")
	if !resp.Result {
		http.Error(w, "Error writing to database", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Added document to database:", body.ID)
}

func (ts *taskServer) updateHandler(w http.ResponseWriter, r *http.Request) {
	var body database.Doc
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if ok := ts.db.Exists(body.ID); !ok {
		http.Error(w, "Data entry does not exist", 404)
		return
	}

	msg := database.Message{
		Action:   "update",
		Document: body,
		Resp:     make(chan database.Response),
	}

	resp := ts.db.Send(msg)
	if !resp.Result {
		http.Error(w, "Error writing to database", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Added document to database:", body.ID)
}

func (ts *taskServer) deleteHandler(w http.ResponseWriter, r *http.Request) {

	if ok := ts.db.Exists(r.PathValue("id")); !ok {
		http.Error(w, "Data entry does not exist", 404)
		return
	}

	msg := database.Message{
		Action:   "delete",
		Document: database.Doc{ID: r.PathValue("id")},
		Resp:     make(chan database.Response),
	}

	resp := ts.db.Send(msg)

	if !resp.Result {
		http.Error(w, "Error writing to database", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Deleted document from database:", r.PathValue("id"))

}

func (ts *taskServer) getHandler(w http.ResponseWriter, r *http.Request) {

	msg := database.Message{
		Action: "list",
		Resp:   make(chan database.Response),
	}

	resp := ts.db.Send(msg)
	if !resp.Result {
		http.Error(w, "Error reading from database", http.StatusInternalServerError)
		return
	}

	jsonBody, err := json.Marshal(resp.Documents)
	if err != nil {
		http.Error(w, "Error reading database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBody)
}
