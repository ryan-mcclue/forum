//go:build !release

package cfg

const IsRelease = false
const DBPath = ".serve/data/db.bolt"
const StaticDir = ".serve/static/"
