package teamsnap

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const basePath = "/v3/"

func (r TeamSnap) loadRoot() teamSnapResult {


	tr, ok := r.makeRequest(r.buildTeamSnapURL(basePath))
	if !ok {
		panic("Banana - Failed to get TeamSnap root")
	}

	return tr
}

func (r TeamSnap) buildTeamSnapURL(path string) string {
	return fmt.Sprintf("%s%s", r.configuration.TeamSnapServer, path)
}

func generateHash(salt string, data string) string {
	// Make string
	tmp := fmt.Sprintf("%s|%s", salt, data)
	bytes := []byte(tmp)

	// Converts string to sha2
	h := sha256.New()                   // new sha256 object
	h.Write(bytes)                      // data is now converted to hex
	code := h.Sum(nil)                  // code is now the hex sum
	codestr := hex.EncodeToString(code) // converts hex to string

	return codestr
}

