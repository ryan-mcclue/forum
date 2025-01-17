package forum

import (
	"forum/cfg"

	"go.hasen.dev/vbeam"
	"go.hasen.dev/vbolt"
)

// local server from filesystem, dev from RAM
func MakeApplication() *vbeam.Application {
	// will shutdown existing server if running on same port
	vbeam.RunBackServer(cfg.Backport)

	db := vbolt.Open(cfg.DBPath)

	// set up RPC
	var app = vbeam.NewApplication("HandcraftedForum", db)

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