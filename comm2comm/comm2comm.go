package comm2comm

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/jmhodges/pac2pac/commparse"
)

// Transaction is a Schedule B transaction of money sent from SendingCommittee
// to ReceivingCommittee.
type Transaction struct {
	SendingCommittee   commparse.CommitteeID
	ReceivingCommittee commparse.CommitteeID
}

// ParseFile parses a file of transactions of money sent from committees (PACs)
// to other committees (PACs). This is the "Any transaction from one committee
// to another" file on the FEC website at
// https://www.fec.gov/data/advanced/?tab=bulk-data .
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

type PACSet map[commparse.CommitteeID]struct{}

func (p PACSet) Add(committee commparse.CommitteeID) {
	p[committee] = struct{}{}
}
func (p PACSet) Has(committee commparse.CommitteeID) bool {
	_, found := p[committee]
	return found
}
func (p PACSet) Del(committee commparse.CommitteeID) {
	delete(p, committee)
}

// HighTrafficSendingLimit is the default lowerLimit used in the
// HighTrafficSendingPACs to determine if a PAC should be returned for sending
// to many other PACs. It was determined by looking at the histogram of PACs and
// identifying
const HighTrafficSendingLimit = 100

// HighTrafficSendingPACs returns a PACSet of the committees that sent to more
// committees (PACs) than lowerLimit. If lowerLimit is less than zero, the
// default HighTrafficSendingLimit constant is used as the default.
func HighTrafficSendingPACs(data Maps, lowerLimit int) PACSet {
	if lowerLimit < 0 {
		lowerLimit = HighTrafficSendingLimit
	}

	set := make(PACSet)
	for committee, recvs := range data.SendingCommitteeToReceivers {
		if len(recvs) > lowerLimit {
			set.Add(committee)
		}
	}
	return set
}

// Maps is a struct for hold what committees send money to what committees
// (SendingCommitteeToReceivers) and what committees receive money from what
// committees (ReceivingCommitteeToSenders).
type Maps struct {
	// SendingCommitteeToReceivers is a map from a commitee id that sent money
	// to the committes it sent money to.
	SendingCommitteeToReceivers map[commparse.CommitteeID]PACSet

	// ReceivingCommitteeToSenders is a map from a commitee id that received
	// money to the committes that sent money to it.
	ReceivingCommitteeToSenders map[commparse.CommitteeID]PACSet
}

func MoneyMapsFromTransactions(trans []Transaction) Maps {
	data := Maps{
		SendingCommitteeToReceivers: make(map[commparse.CommitteeID]PACSet),
		ReceivingCommitteeToSenders: make(map[commparse.CommitteeID]PACSet),
	}
	for _, t := range trans {
		recvs, found := data.SendingCommitteeToReceivers[t.SendingCommittee]
		if !found {
			recvs = make(PACSet, 1)
		}
		recvs.Add(t.ReceivingCommittee)
		if !found {
			data.SendingCommitteeToReceivers[t.SendingCommittee] = recvs
		}

		sends, found := data.ReceivingCommitteeToSenders[t.SendingCommittee]
		if !found {
			sends = make(PACSet, 1)
		}
		sends.Add(t.SendingCommittee)
		if !found {
			data.ReceivingCommitteeToSenders[t.ReceivingCommittee] = sends
		}
	}
	return data
}
