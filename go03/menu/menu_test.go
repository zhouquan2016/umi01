package menu

import (
	"log"
	"regexp"
	"testing"
)

func TestPath(t *testing.T) {
	path := "/123/a_a-/dqd1"
	ok, err := regexp.Match("^(/[\\w_-]+)+$", []byte(path))
	log.Println(ok, err)
}
func TestLen(t *testing.T) {

	log.Println(len([]rune("菜单")))
}
