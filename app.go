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

	return app
}

// local server from filesystem, dev from RAM

//func main() {
//	fmt.Println("hi there")
//}
