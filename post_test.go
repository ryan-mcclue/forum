package forum

import (
	"os"
	"testing"

	"go.hasen.dev/vbeam"
	"go.hasen.dev/vbolt"
)

// no need to mock database

func TestPosting(t *testing.T) {
	testDBPath := "test.db"
	db := OpenDB(testDBPath)
	defer os.Remove(testDBPath)

	reqs := []CreatePostRequest{
		{UserId: 1, Content: "Hello #World #T1"},
		{UserId: 1, Content: "Hello #World #T2"},
		{UserId: 1, Content: "#Hello World #T3"},

		{UserId: 2, Content: "Hello #World #T1"},
		{UserId: 2, Content: "#Hello World #T3"},

		{UserId: 3, Content: "#Hello #World #T3"},
	}

	tagsCount := map[string]int{
		"T1":    2,
		"T2":    1,
		"T3":    3,
		"World": 4,
		"Hello": 3,
	}

	for _, req := range reqs {
		var ctx vbeam.Context
		ctx.Tx = vbolt.ReadTx(db)

		_, err := CreatePost(&ctx, req)
		if err != nil {
			t.Fatalf("Post Creation Failed: %v", err)
		}

		vbolt.TxClose(ctx.Tx)
	}

	for tag, count := range tagsCount {
		// go test -v
		t.Logf("%s =>", tag)
		var ctx vbeam.Context
		ctx.Tx = vbolt.ReadTx(db)

		res, err := PostsByHashTag(&ctx, ByHashTagReq{req})
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Posts) != count {
			t.Fatalf("Expected: %d, actual: %d", count, len(res.Posts))
		}

		vbolt.TxClose(ctx.Tx)
	}
}
