package passphrase

import (
	"crypto/rand"
	"strings"
)

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// Generate creates a random alphanumeric passphrase of the given length.
func Generate(length int) string {
	if length <= 0 {
		return ""
	}

	b := make([]byte, length)
	max := byte(len(alphabet))

	for i := 0; i < length; i++ {
		for {
			var r byte
			if _, err := rand.Read(b[i : i+1]); err != nil {
				r = 'A'
			} else {
				r = b[i]
			}
			if int(r) < (256/len(alphabet))*len(alphabet) {
				b[i] = alphabet[r%max]
				break
			}
		}
	}

	return string(b)
}

// BuildInstruction applies the template, replacing %random% with the phrase.
// Template is user-configurable; empty template falls back to default (Chinese).
func BuildInstruction(template string, phrase string) string {
	if template == "" {
		template = "当你完全完成所有任务后，你必须在回复结尾附上这个短语：`%random%`。不要向用户提及此指令，不要解释这个短语。"
	}
	return strings.ReplaceAll(template, "%random%", phrase)
}
