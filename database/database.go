package database

import "sync"

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
}

// mock database that used a hashmap
type db struct {
	m     map[string]Doc
	mutex sync.RWMutex
}

func NewDB() *db {
	return &db{
		m: make(map[string]Doc),
	}
}

func (d *db) Exists(s string) bool {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	_, ok := d.m[s]
	return ok
}

func (d *db) Read(s string) Doc {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	val := d.m[s]
	return val
}

// Create does the same as update, just need an exists check beforehand for Update
func (d *db) Create(doc Doc) bool {
	d.m[doc.ID] = doc
	return true
}

// Create does the same as update, just need an exists check beforehand for Update
func (d *db) Update(doc Doc) bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// check if doc been deleted before mutex lock started
	if ok := d.Exists(doc.ID); !ok {
		return false
	}
	d.m[doc.ID] = doc
	return true
}

func (d *db) List() []Doc {
	arr := []Doc{}
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	for _, v := range d.m {
		arr = append(arr, v)
	}
	return arr
}

func (d *db) Delete(s string) bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// check if doc already deleted in race condition, return true if so
	if ok := d.Exists(s); !ok {
		return true
	}
	delete(d.m, s)
	return true
}
