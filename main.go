package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	minSpins                 = 30
	maxSpins                 = 80
	nToDisplayBeforeAndAfter = 2
)

var (
	spinPause     = time.Millisecond * 50
	inputFileFlag = flag.String("file", "", "Path to a file, containing a list of sentences, new-line delimited")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		fmt.Println("ERR: ", err)
	}
}

func run() error {
	opts, err := loadOpts(*inputFileFlag)
	if err != nil {
		return err
	}

	if len(opts) == 0 {
		return fmt.Errorf("no options found in file")
	}

	optionsList := makeLinkedLoop(opts)

	startingIndex := rand.Int() % len(opts)
	node := optionsList[startingIndex]
	nSpins := (rand.Int() % (maxSpins - minSpins)) + minSpins

	for i := 0; i < nSpins; i++ {
		currentValue := node.Val
		before, after := getDisplayOptions(node, nToDisplayBeforeAndAfter)

		clearStdout()
		printList(before)
		printChoice(currentValue)
		printList(after)

		time.Sleep(spinPause)
		node = node.Next
	}

	time.Sleep(time.Minute * 5)

	return nil
}

func clearStdout() {
	// linux & macOS
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func getDisplayOptions[T any](currentNode *listNode[T], nBeforeAndAfter int) (before, after []T) {
	afterNode := currentNode.Next

	beforeNode := currentNode
	// step back to start of befores
	for range make([]bool, nBeforeAndAfter) {
		beforeNode = beforeNode.Prev
	}

	for range make([]bool, nBeforeAndAfter) {
		before = append(before, beforeNode.Val)
		after = append(after, afterNode.Val)

		beforeNode = beforeNode.Next
		afterNode = afterNode.Next
	}

	return
}

func loadOpts(path string) ([]string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return []string{}, fmt.Errorf("failed to read file: %w", err)
	}

	var s []string
	for _, bb := range bytes.Split(b, []byte{'\n'}) {
		if strings.TrimSpace(string(bb)) != "" {
			s = append(s, string(bb))
		}
	}
	return s, nil
}

type listNode[T any] struct {
	Next, Prev *listNode[T]
	Val        T
}

func makeLinkedLoop[T any](values []T) []*listNode[T] {
	length := len(values)

	nodes := make([]*listNode[T], length)
	for i, v := range values {
		nodes[i] = &listNode[T]{
			Val: v,
		}
	}

	for i, n := range nodes {
		// adding the length allows negative cases to loop in reverse
		prevIndex := ((i - 1) + length) % length
		nextIndex := (i + 1) % length

		n.Next = nodes[nextIndex]
		n.Prev = nodes[prevIndex]
	}

	return nodes
}

var (
	current = color.New(color.FgGreen, color.Bold)
	others  = color.New(color.FgYellow, color.Faint)
)

func printChoice(s string) {
	current.Printf("> %s <\n", s)
}

func printList(ss []string) {
	for _, s := range ss {
		others.Println(s)
	}
}
