package es

import (
	"go03/db"
	"log"
	"testing"
)

func TestFindByPage(t *testing.T) {
	page := HistoryEs.FindByPage(10, 1)
	log.Println(page)
}

func TestHistoryEs_FindByScroll(t *testing.T) {
	HistoryEs.FindByScroll(10, func(histories []*db.History) {
		log.Println("size:", len(histories), histories)
	})
}
