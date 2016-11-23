package sockaddr_test

import (
	"testing"

	sockaddr "github.com/hashicorp/go-sockaddr"
)

func TestVisitAllRFCs(t *testing.T) {
	const expectedNumRFCs = 28
	numRFCs := 0
	sockaddr.VisitAllRFCs(func(rfcNum uint, sas sockaddr.SockAddrs) {
		numRFCs++
	})
	if numRFCs != expectedNumRFCs {
		t.Fatalf("wrong number of RFCs: %d", numRFCs)
	}
}
