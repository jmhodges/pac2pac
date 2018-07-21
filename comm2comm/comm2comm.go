package comm2comm

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/jmhodges/pac2pac/commparse"
)

// Transaction is a transaction of money sent from SendingCommittee
// to ReceivingCommittee.
type Transaction struct {
	SendingCommittee   commparse.CommitteeID
	ReceivingCommittee commparse.CommitteeID
}

// ParseFile parses a file of transactions of money sent from committees to
// other committees. This is the "Any transaction from one committee to another"
// file on the FEC website at https://www.fec.gov/data/advanced/?tab=bulk-data .
func ParseFile(fp string) ([]Transaction, error) {
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, err
	}
	lines := bytes.Split(b, []byte{'\n'})
	trans := make([]Transaction, 0, len(lines))
	for i, line := range lines {
		if len(line) == 0 {
			continue
		}
		cols := bytes.Split(line, []byte{'|'})
		if len(cols) != 21 {
			return nil, fmt.Errorf("on line %d, expected 21 columns, got %d", i+1, cols)
		}
		trans = append(trans, Transaction{
			SendingCommittee:   commparse.CommitteeID(string(cols[0])),
			ReceivingCommittee: commparse.CommitteeID(string(cols[15])),
		})
	}
	return trans, nil
}

type CommIDSet map[commparse.CommitteeID]struct{}

func (p CommIDSet) Add(committee commparse.CommitteeID) {
	p[committee] = struct{}{}
}
func (p CommIDSet) Has(committee commparse.CommitteeID) bool {
	_, found := p[committee]
	return found
}
func (p CommIDSet) Del(committee commparse.CommitteeID) {
	delete(p, committee)
}

// Maps is a struct for hold what committees send money to what committees
// (SendingCommitteeToReceivers) and what committees receive money from what
// committees (ReceivingCommitteeToSenders).
type Maps struct {
	// SendingCommitteeToReceivers is a map from a commitee id that sent money
	// to the committes it sent money to.
	SendingCommitteeToReceivers map[commparse.CommitteeID]CommIDSet

	// ReceivingCommitteeToSenders is a map from a commitee id that received
	// money to the committes that sent money to it.
	ReceivingCommitteeToSenders map[commparse.CommitteeID]CommIDSet
}

func MoneyMapsFromTransactions(trans []Transaction) Maps {
	data := Maps{
		SendingCommitteeToReceivers: make(map[commparse.CommitteeID]CommIDSet),
		ReceivingCommitteeToSenders: make(map[commparse.CommitteeID]CommIDSet),
	}
	for _, t := range trans {
		recvs, found := data.SendingCommitteeToReceivers[t.SendingCommittee]
		if !found {
			recvs = make(CommIDSet, 1)
		}
		recvs.Add(t.ReceivingCommittee)
		if !found {
			data.SendingCommitteeToReceivers[t.SendingCommittee] = recvs
		}

		sends, found := data.ReceivingCommitteeToSenders[t.SendingCommittee]
		if !found {
			sends = make(CommIDSet, 1)
		}
		sends.Add(t.SendingCommittee)
		if !found {
			data.ReceivingCommitteeToSenders[t.ReceivingCommittee] = sends
		}
	}
	return data
}
