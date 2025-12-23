package main

import (
	"bufio"
	"fmt"
	"os"
	"texas-holdem-go/poker"
)

func main() {
	if len(os.Args) > 1 {
		// Evaluate arguments as hands
		for _, arg := range os.Args[1:] {
			evaluateStr(arg)
		}
		return
	}

	// Interactive mode
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter poker hands (e.g., 'H2 SQ C2 D2 CQ') or 'quit' to exit:")
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if line == "quit" || line == "exit" {
			break
		}
		evaluateStr(line)
	}
}

func evaluateStr(s string) {
	hand, err := poker.NewHand(s)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	val, kickers := hand.Evaluate()
	fmt.Printf("Hand: %s\nValue: %s\nKickers: %v\n", s, val, kickers)
}

