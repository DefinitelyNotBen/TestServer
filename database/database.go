package database

import "fmt"

type Doc struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

// interface wrapper for database
// added a read/write mutex, only going to use write lock for now, but leaving door open for read lock
type DB interface {
	Create(Doc) bool
	Update(Doc) bool
	Delete(string) bool
	Read(string) Doc
	List() []Doc
	Exists(string) bool
	Send(Message) Response
	Listen()
}

type Message struct {
	Action   string
	Document Doc
	Resp     chan Response
}

// mock database that used a hashmap
type db struct {
	m map[string]Doc
	c chan Message
}

// Response to pass to the server based on result of request
type Response struct {
	Result      bool
	Description string
	Documents   []Doc
}

func NewDB() *db {
	return &db{
		m: make(map[string]Doc),
		c: make(chan Message),
	}
}

// Send Message to server, receive response
func (d *db) Send(m Message) Response {
	fmt.Println("Sending message to the database. Action:", m.Action)
	d.c <- m

	resp := <-m.Resp
	close(m.Resp)
	return resp
}

// listen to server, messages sent will tell which action needs to be taken
func (d *db) Listen() {
	fmt.Println("Database is now listening to the server")
	for {
		select {
		case msg := <-d.c:
			// create a response struct to give result to server
			r := Response{Result: true}

			switch msg.Action {
			case "exists":
				fmt.Println("Database recieved exists message")
				r.Result = d.Exists(msg.Document.ID)

			case "create":
				fmt.Println("Database recieved create message")
				r.Result = d.Create(msg.Document)

			case "delete":
				fmt.Println("Database recieved delete message")
				r.Result = d.Delete(msg.Document.ID)

			case "read":
				fmt.Println("Database recieved read message")
				document := d.Read(msg.Document.ID)

				// if document ID is empty, the document doesn't exist
				if document.ID != "" {
					r.Documents = append(r.Documents, document)
				} else {
					r.Result = false
				}

			case "list":
				fmt.Println("Database recieved list message")
				docs := d.List()

				fmt.Println(len(docs), " documents found!")
				// if no documents returned then assume there was an error
				if len(docs) != 0 {
					r.Documents = append(r.Documents, docs...)
				} else {
					r.Result = false
				}
				fmt.Printf("%v", r)
			}

			// end by sending response back to server
			msg.Resp <- r
		}

	}
}

func (d *db) Exists(s string) bool {
	_, ok := d.m[s]
	return ok
}

func (d *db) Read(s string) Doc {
	if !d.Exists(s) {
		return Doc{ID: ""}
	}
	val := d.m[s]
	return val
}

func (d *db) Create(doc Doc) bool {
	d.m[doc.ID] = doc
	return true
}

// Create does the same as update, just need an exists check beforehand for Update
func (d *db) Update(doc Doc) bool {
	if !d.Exists(doc.ID) {
		return false
	}

	d.m[doc.ID] = doc
	return true
}

func (d *db) List() []Doc {
	arr := []Doc{}
	for _, v := range d.m {
		arr = append(arr, v)
	}
	return arr
}

func (d *db) Delete(s string) bool {
	if !d.Exists(s) {
		return false
	}

	delete(d.m, s)
	return true
}
