package main

import (
	"fmt"

	"github.com/ashleyschuett/kubeconfig-cleanup/pkg/config"
)

func main() {
	m := config.NewManager()
	valid := make(map[string]bool, 0)

	for id, context := range m.Original.Contexts {
		fmt.Printf("Testing cluster for context %s...\n", id)
		ok, tested := valid[context.Cluster]

		if !tested {
			ok = m.Validate(context)
			valid[context.Cluster] = ok
		}

		if !ok {
			m.RemoveContext(id, context)
		}
		fmt.Println()
	}

	m.RemoveUnusedUsers()
	m.Finish()
}
