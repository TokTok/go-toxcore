package tox

import (
	"fmt"
	"log"
	"testing"
)

func TestGid0(t *testing.T) {
	mt := NewMiniTox()
	t0 := mt.t

	go mt.Iterate()

	for i := 0; i < 5; i++ {
		gn, err := t0.ConferenceNew()
		log.Println("gn:", gn, err)
		t0.ConferenceSetTitle(gn, fmt.Sprintf("group###%d", gn))
		id, err := t0.ConferenceGetIdentifier(gn)
		log.Println("id:", id, err)
	}
}
