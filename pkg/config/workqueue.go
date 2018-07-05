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

type ContextResult struct {
	ID      string
	Context *clientcmdapi.Context
	Message string
	Valid   bool
}

func (m *Manager) ValidateAndAddToWorkqueue(id string, context *clientcmdapi.Context) {
	valid, message := m.Validate(context)
	m.workqueue <- ContextResult{id, context, message, valid}
	m.contextedValidated++

	if m.contextedValidated == m.totalContexts {
		close(m.workqueue)
	}
}

func (m *Manager) runWorkqueue() {
	for m.processNextWorkItem() {
	}

	return
}

func (m *Manager) processNextWorkItem() bool {
	s.Start()
	for n := range m.workqueue {
		s.Stop()
		fmt.Printf("\nTesting context's '%s' cluster...\n", n.ID)
		fmt.Println(n.Message)

		if !n.Valid {
			m.RemoveContext(n.ID, n.Context)
		}

		s.Start()
	}
	s.Stop()

	return false
}
