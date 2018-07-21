package comm2comm

import (
	"testing"

	"github.com/jmhodges/pac2pac/commparse"
)

func TestGolden(t *testing.T) {
	trans, err := ParseFile("./itoth-test.txt")
	if err != nil {
		t.Fatalf("ParseFile: %s", err)
	}
	expectedCount := 100
	if len(trans) != expectedCount {
		t.Errorf("number of transactions: expected %d, got %d", expectedCount, len(trans))
	}
	tr := trans[len(trans)-1]
	expectedSend := commparse.CommitteeID("C00012880")
	if tr.SendingCommittee != expectedSend {
		t.Errorf("SendingCommittee: expected %#v, got %#v", expectedSend, tr.SendingCommittee)
	}
	expectedRecv := commparse.CommitteeID("C00428052")
	if tr.ReceivingCommittee != expectedRecv {
		t.Errorf("ReceivingCommittee: expected %#v, got %#v", expectedRecv, tr.ReceivingCommittee)
	}
}
