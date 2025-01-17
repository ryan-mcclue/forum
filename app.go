package forum

import (
	"forum/cfg"

	"go.hasen.dev/vbeam"
	"go.hasen.dev/vbolt"
)

var dbInfo vbolt.Info

type User struct {
	Id int
	Username string
	Email string
	IsAdmin bool
}

// buckets (map)
// indexes (bidirectional multimap: many-to-manys?) 
// collections (hierarchy of keys?)

// both serialisation/deserialisation
func PackUser(self *User, buf *vpack.Buffer) {
	vpack.Version(1, buf)
	vpack.Int(&self.Id, buf)
	vpack.String(&self.Username, buf)
	vpack.String(&self.Email, buf)
	vpack.Bool(&self.IsAdmin, buf)
}

// fixed int key
var UsersBkt = vbolt.Bucket(&dbInfo, "users", vpack.FInt, PackUser)
var PasswdBkt = vbolt.Bucket(&dbInfo, "passwd", vpack.FInt, vpack.ByteSlice)
var UsernameBkt = vbolt.Bucket(&dbInfo, "username", vpack.StringZ, vpack.Int)

// NOTE: no pagination
func fetchUsers(tx *vbolt.Tx) (users []User) {
	vbolt.IterateAll(tx, UsersBkt, func(key int, value User) bool {
		generic.Append(&users, value)
		return true
	})
	return
}

// local server from filesystem, dev from RAM
func MakeApplication() *vbeam.Application {
	// will shutdown existing server if running on same port
	vbeam.RunBackServer(cfg.Backport)

	db := vbolt.Open(cfg.DBPath)
	vbolt.WithWriteTx(db, func(tx *vbolt.Tx) {
		vbolt.TxRawBucket(tx, "proc")
		vbolt.EnsureBuckets(tx, &dbInfo)
		tx.Commit()
	})

	var app = vbeam.NewApplication("HandcraftedForum", db)

	// set up RPC
	vbeam.RegisterProc(app, AddUser)
	vbeam.RegisterProc(app, ListUsers)

	return app
}

// REST is a network architecture that enforces a hypermedia constraint, i.e. returns HTML and so is self-describing.
// This is in contrast to RPC (e.g. JSON), where you would have to know how to interpret fields, i.e. have out-of-band knowledge (server and client are coupled)


//var usernames []string
var usernames = make([]string, 0)

type AddUserRequest struct {
	Username string
}

type UserListResponse struct {
	AllUsernames []string
}

func AddUser(ctx *vbeam.Context, req AddUserRequest) (resp UserListResponse, err error) {
	usernames = append(usernames, req.Username)
	resp.AllUsernames = usernames
	return
}

type EmptyRequest struct {}

func ListUsers(ctx *vbeam.Context, req EmptyRequest) (resp UserListResponse, err error) {
	resp.AllUsernames = usernames
	return
}