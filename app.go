package forum

import (
	"forum/cfg"

	"go.hasen.dev/vbeam"
	"go.hasen.dev/vbolt"
)

func MakeApplication() *vbeam.Application {
	// will shutdown existing server if running on same port
	vbeam.RunBackServer(cfg.Backport)

	db := vbolt.Open(cfg.DBPath)

	// set up RPC
	var app = vbeam.NewApplication("HandcraftedForum", db)

	vbeam.RegisterProc(app, AddUser)

	return app
}

var usernames []string

type AddUserRequest struct {
	Username string
}

type UserListResponse struct {
	AllUsernames []string
}

// https://hasen.substack.com/p/automagic-go-typescript-interface

func AddUser(ctx *vbeam.Context, req AddUserRequest) (resp UserListResponse, err error) {
	usernames = append(usernames, req.Username)
	resp.AllUsernames = usernames
	return
}

// local server from filesystem, dev from RAM

//func main() {
//	fmt.Println("hi there")
//}
