package forum

import (
	"forum/cfg"
	"strings"
	"time"

	"go.hasen.dev/generic"
	"go.hasen.dev/vbeam"
	"go.hasen.dev/vbolt"
	"go.hasen.dev/vpack"
)

var dbInfo vbolt.Info

type Post struct {
	Id        int
	UserId    int
	CreatedAt time.Time
	Content   string

	Tags []string
}

func ExtractHashTags(content string) (tags []string) {
	const max_tag_length = 20
	at := 0
	for at < len(content) {
		hash := strings.IndexByte(content[at:], '#')
		if hash == -1 {
			break
		}
		at = hash + 1
		end_tag := strings.IndexAny(content[at:], " \n\t")
		if end_tag == -1 {
			end_tag = at + len(content[at:])
		}
		tag := content[at:end_tag]
		if len(tag) > max_tag_length {
			tag = tag[:max_tag_length]
		}
		tags = append(tags, tag)

		at = end_tag
	}
	return
}

func PackPost(self *Post, buf *vpack.Buffer) {
	vpack.Version(1, buf)
	vpack.Int(&self.Id, buf)
	vpack.Int(&self.UserId, buf)
	vpack.UnixTime(&self.CreatedAt, buf)
	vpack.String(&self.Content, buf)
	vpack.Slice(&self.Tags, vpack.String, buf)
}

var PostsBkt = vbolt.Bucket(&dbInfo, "posts", vpack.FInt, PackPost)

// cursor pagination based on bytes, so cannot say jump to page 20 like with an offset
// however, apart from that more efficient
// keys are 3-tuple (term, priority, target)
// (user_id1, 12:01, post_id1)
// (user_id1, 12:02, post_id2)
var UserPostsIdx = vbolt.IndexExt(&dbInfo, "user-posts", vpack.FInt, vpack.UnixTimeKey, vpack.FInt)

// hashtag -> post_id
var HashTagsIdx = vbolt.IndexExt(&dbInfo, "hashtags", vpack.StringZ, vpack.UnixTimeKey, vpack.FInt)

type CreatePostRequest struct {
	UserId  int
	Content string
}

// TODO: unify with codepaths over intertwining

func CreatePost(ctx *vbeam.Context, req CreatePostRequest) (resp Post, err error) {
	const MaxPostSize = 1024
	if len(req.Content) > MaxPostSize {
		req.Content = req.Content[:MaxPostSize]
	}
	tags := ExtractHashTags(req.Content)

	vbeam.UseWriteTx(ctx)

	resp.Id = vbolt.NextIntId(ctx.Tx, PostsBkt)
	resp.UserId = req.UserId
	resp.Content = req.Content
	resp.CreatedAt = time.Now()
	resp.Tags = tags

	vbolt.Write(ctx.Tx, PostsBkt, resp.Id, &resp)

	vbolt.SetTargetSingleTermExt(
		ctx.Tx,
		UserPostsIdx,
		resp.Id, // targets
		resp.CreatedAt,
		resp.UserId, // key
	)

	vbolt.SetTargetTermsUniform(
		ctx.Tx,
		HashTagsIdx,
		resp.Id,
		tags, // key
		resp.CreatedAt,
	)

	vbolt.TxCommit(ctx.Tx)

	return
}

type Posts struct {
	Posts  []Post
	Cursor []byte
}
type ByUserReq struct {
	UserId int
	Cursor []byte
}

const Limit = 2

func PostsByUser(ctx *vbeam.Context, req ByUserReq) (resp Posts, err error) {
	var window = vbolt.Window{
		Limit:     Limit,
		Direction: vbolt.IterateReverse,
		Cursor:    req.Cursor,
	}
	var postIds []int
	resp.Cursor = vbolt.ReadTermTargets(
		ctx.Tx,
		UserPostsIdx,
		req.UserId,
		&postIds,
		window,
	)
	vbolt.ReadSlice(ctx.Tx, PostsBkt, postIds, &resp.Posts)

	generic.EnsureSliceNotNil(&resp.Posts)
	generic.EnsureSliceNotNil(&resp.Cursor)

	return
}

func OpenDB(dbpath string) *vbolt.DB {
	db := vbolt.Open(dbpath)
	vbolt.WithWriteTx(db, func(tx *vbolt.Tx) {
		vbolt.TxRawBucket(tx, "proc") // special
		vbolt.EnsureBuckets(tx, &dbInfo)
		tx.Commit()
	})
	return db
}

// local server from filesystem, dev from RAM
func MakeApplication() *vbeam.Application {
	// will shutdown existing server if running on same port
	vbeam.RunBackServer(cfg.Backport)

	db := OpenDB(cfg.DBPath)

	var app = vbeam.NewApplication("HandcraftedForum", db)

	// set up RPC
	vbeam.RegisterProc(app, AddUser)
	vbeam.RegisterProc(app, ListUsers)

	return app
}

// REST is a network architecture that enforces a hypermedia constraint, i.e. returns HTML and so is self-describing.
// This is in contrast to RPC (e.g. JSON), where you would have to know how to interpret fields, i.e. have out-of-band knowledge (server and client are coupled)
