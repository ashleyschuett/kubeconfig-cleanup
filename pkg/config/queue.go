package config

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var s *spinner.Spinner

func init() {
	s = spinner.New(spinner.CharSets[9], 100*time.Millisecond)
}

type Node struct {
	ID      string
	Context *clientcmdapi.Context
	Valid   bool
	Message string
}

// Stack is a basic LIFO stack that resizes as needed.
type Stack struct {
	nodes     []*Node
	count     int
	toProcess int
	processed int
}

// NewStack returns a new stack.
func NewStack(tc int) *Stack {
	return &Stack{
		toProcess: tc,
		processed: 0,
	}
}

func (s *Stack) NewNode(id, message string, valid bool, c *clientcmdapi.Context) *Node {
	return &Node{
		id,
		c,
		valid,
		message,
	}
}

// Push adds a node to the stack.
func (s *Stack) Push(n *Node) {
	s.nodes = append(s.nodes[:s.count], n)
	s.count++
}

// Pop removes and returns a node from the stack in last to first order.
func (s *Stack) Pop() *Node {
	if s.count == 0 {
		return nil
	}
	s.count--
	return s.nodes[s.count]
}

func (m *Manager) runWorkqueue() {
	for m.processNextWorkItem() {
	}

	return
}

func (m *Manager) processNextWorkItem() bool {
	if m.workqueue.processed == m.workqueue.toProcess {
		return false
	}
	s.Start()

	n := m.workqueue.Pop()
	if n == nil {
		return true
	}
	s.Stop()

	fmt.Printf("Testing context's '%s' cluster...\n", n.ID)
	fmt.Println(n.Message)

	if !n.Valid {
		m.RemoveContext(n.ID, n.Context)
	}

	m.workqueue.processed++
	return true
}
