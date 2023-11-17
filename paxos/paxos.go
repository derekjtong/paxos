package paxos

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Node struct {
	ID          int
	Proposal    string
	Accepted    bool
}

func (n *Node) Prepare(proposal string) (string, bool) {
	if !n.Accepted {
		n.Proposal = proposal
		return n.Proposal, true
	}
	return n.Proposal, false
}

func (n *Node) Accept(proposal string) bool {
	if !n.Accepted {
		n.Proposal = proposal
		n.Accepted = true
		return true
	}
	return false
}

func RunPaxos(nodes []*Node, proposal string) {
	// Simulating Prepare phase
	var wg sync.WaitGroup
	for _, node := range nodes {
		wg.Add(1)
		go func(n *Node) {
			defer wg.Done()
			// Simulate network delay
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			prop, ok := n.Prepare(proposal)
			if ok {
				fmt.Printf("Node %d prepared proposal: %s\n", n.ID, prop)
			} else {
				fmt.Printf("Node %d already accepted proposal: %s\n", n.ID, n.Proposal)
			}
		}(node)
	}
	wg.Wait()

	// Simulating Accept phase
	var acceptedProposal string
	for _, node := range nodes {
		if node.Proposal == proposal {
			acceptedProposal = proposal
			break
		}
	}
	if acceptedProposal == proposal {
		// Proposal was accepted by majority, proceed with the Accept phase
		var acceptWg sync.WaitGroup
		for _, node := range nodes {
			acceptWg.Add(1)
			go func(n *Node) {
				defer acceptWg.Done()
				// Simulate network delay
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
				ok := n.Accept(acceptedProposal)
				if ok {
					fmt.Printf("Node %d accepted proposal: %s\n", n.ID, acceptedProposal)
				} else {
					fmt.Printf("Node %d already accepted proposal: %s\n", n.ID, n.Proposal)
				}
			}(node)
		}
		acceptWg.Wait()
	}

	// After consensus, check the accepted proposal
	for _, node := range nodes {
		if node.Accepted {
			fmt.Printf("Node %d has accepted proposal: %s\n", node.ID, node.Proposal)
		}
	}
}
