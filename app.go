package forum

import (
	"errors"
	"forum/cfg"

	"go.hasen.dev/generic"
	"go.hasen.dev/vbeam"
	"go.hasen.dev/vbolt"
	"go.hasen.dev/vpack"

	"golang.org/x/crypto/bcrypt"
)

var dbInfo vbolt.Info

type User struct {
	Id       int
	Username string
	Email    string
	IsAdmin  bool
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

// var usernames []string
// var usernames = make([]string, 0)

func fetchUsers(tx *vbolt.Tx) (users []User) {
	vbolt.IterateAll(tx, UsersBkt, func(key int, value User) bool {
		generic.Append(&users, value)
		return true
	})
	return
}

type AddUserRequest struct {
	Username string
	Email    string
	Password string
}

type UserListResponse struct {
	Users []User
}

var UsernameTaken = errors.New("UsernameTaken")

func AddUser(ctx *vbeam.Context, req AddUserRequest) (resp UserListResponse, err error) {
	if vbolt.HasKey(ctx.Tx, UsernameBkt, req.Username) {
		err = UsernameTaken
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	// start write transaction
	vbeam.UseWriteTx(ctx)

	var user User
	user.Id = vbolt.NextIntId(ctx.Tx, UsersBkt)
	user.Username = req.Username
	user.Email = req.Email
	user.IsAdmin = user.Id < 2

	vbolt.Write(ctx.Tx, UsersBkt, user.Id, &user)
	vbolt.Write(ctx.Tx, PasswdBkt, user.Id, &hash)
	vbolt.Write(ctx.Tx, UsernameBkt, user.Username, &user.Id)

	resp.Users = fetchUsers(ctx.Tx)
	generic.EnsureSliceNotNil(&resp.Users)

	// commit transaction
	vbolt.TxCommit(ctx.Tx)
	return
}

type EmptyRequest struct{}

func ListUsers(ctx *vbeam.Context, req EmptyRequest) (resp UserListResponse, err error) {
	resp.Users = fetchUsers(ctx.Tx)
	generic.EnsureSliceNotNil(&resp.Users)
	return
}
