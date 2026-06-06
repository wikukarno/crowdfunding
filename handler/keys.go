package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var safeExt = regexp.MustCompile(`^\.[a-z0-9]{1,8}$`)

// objectKey builds a collision-free, non-guessable storage key. The original
// filename is dropped (only a sanitized extension is kept) so a user can't
// influence the path, cause collisions, or leak their local filename.
func objectKey(prefix, filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if !safeExt.MatchString(ext) {
		ext = ""
	}
	return fmt.Sprintf("%s/%s%s", prefix, randomHex(16), ext)
}

func randomHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return strings.Repeat("0", n*2)
	}
	return hex.EncodeToString(b)
}
