package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type prompter struct {
	reader *bufio.Reader
}

func NewPrompter() *prompter {
	return &prompter{
		bufio.NewReader(os.Stdin),
	}
}

func (p *prompter) getValidYNInput(s string) bool {
	text := strings.ToLower(p.getInput(s))

	if text != "n" && text != "y" {
		return p.getValidYNInput(s)
	}

	return text == "y"
}

func (p *prompter) getInput(s string) string {
	fmt.Printf(s)
	text, _ := p.reader.ReadString('\n')
	return strings.Trim(text, " \n")
}

func (p *prompter) RemoveContext(id string) bool {
	return p.getValidYNInput(fmt.Sprintf("Remove '%s' from context (Y/n)? ", id))
}

func (p *prompter) RemoveUser(id string) bool {
	return p.getValidYNInput(fmt.Sprintf("Remove unused user '%s' from config (Y/n)? ", id))
}

func (p *prompter) WriteConfig() bool {
	return p.getValidYNInput(fmt.Sprintf("Overwrite kubeconfig (Y/n)? "))
}

func (p *prompter) WriteConfigToPath() bool {
	return p.getValidYNInput(fmt.Sprintf("Write kubeconfig to a different path (Y/n)? "))
}

func (p *prompter) GetPath() string {
	return p.getInput(fmt.Sprintf("Path? "))
}
