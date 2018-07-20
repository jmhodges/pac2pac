package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/jmhodges/pac2pac/comm2comm"
	"github.com/jmhodges/pac2pac/commparse"
)

var (
	comm2commFile = flag.String("comm2comm", "", "the file containing the committee to committee transactions (PAC to PAC transactions), often something like 'itoth.txt' (required)")
	commFile      = flag.String("comm", "", "the file containing the list of all committees (PACs) and their information (a.k.a. the 'Committee master'), often something like 'cm.txt' (required)")
)

func main() {
	flag.Parse()
	if *comm2commFile == "" {
		fmt.Fprintf(os.Stderr, "pac2pac: missing -comm2comm argument\n")
		flag.Usage()
		os.Exit(2)
	}
	if *commFile == "" {
		fmt.Fprintf(os.Stderr, "pac2pac: missing -comm argument\n")
		flag.Usage()
		os.Exit(2)
	}
	allComms, err := commparse.ParseFile(*commFile)
	if err != nil {
		log.Fatalf("unable to parse comm file (from -comm): %s", err)
	}
	idToName := make(map[commparse.CommitteeID]string, len(allComms))
	for _, cm := range allComms {
		idToName[cm.ID] = cm.Name
	}
	trans, err := comm2comm.ParseFile(*comm2commFile)
	if err != nil {
		log.Fatalf("unable to parse comm2comm files (from -comm2comm): %s", err)
	}
	transMaps := comm2comm.MoneyMapsFromTransactions(trans)
	type commCount struct {
		id        commparse.CommitteeID
		name      string
		recvcount int
		sendcount int
	}
	counts := []commCount{}
	for comm, recvs := range transMaps.SendingCommitteeToReceivers {
		sendcount := len(transMaps.ReceivingCommitteeToSenders[comm])
		counts = append(counts, commCount{id: comm, name: idToName[comm], recvcount: len(recvs), sendcount: sendcount})
	}
	sort.Slice(counts, func(i, j int) bool {
		switch {
		case counts[i].recvcount < counts[j].recvcount:
			return true
		case counts[i].recvcount == counts[j].recvcount:
			return counts[i].id < counts[j].id
		default:
			return false
		}
	})
	csvw := csv.NewWriter(os.Stdout)
	for _, cc := range counts {
		csvw.Write([]string{strconv.Itoa(cc.recvcount), strconv.Itoa(cc.sendcount), string(cc.id), cc.name})
		// if cc.sendcount == 1 {
		// 	csvw.Flush()
		// 	fmt.Println("and senders are", transMaps.ReceivingCommitteeToSenders[cc.id])
		// }
	}
	csvw.Flush()
}
