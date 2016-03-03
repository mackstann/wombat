package main


import "bufio"
import "fmt"
import "os"
import "strings"
import "time"

func otherMode(currentMode string) string {
	if currentMode == "w" {
		return "b"
	}
	return "w"
}

func modeName(currentMode string) string {
	if currentMode == "w" {
		return "Work!"
	}
	return "Take a break."
}

func isOtherMode(currentMode, input string) bool {
	if currentMode == "w" {
		return input == "b"
	}
	return input == "w"
}

func printInputOptions(currentMode string) {
	fmt.Println(" Options:")
	if currentMode == "b" {
		fmt.Println("  [w]ork")
	} else if currentMode == "w" {
		fmt.Println("  [b]reak")
	}
	fmt.Println("  [q]uit")
}

func runTimer(currentMode string) bool {
	fmt.Println(modeName(currentMode))

	stdinCh := make(chan string)
	go func(ch chan string) {
		reader := bufio.NewReader(os.Stdin)
		for {
			s, err := reader.ReadString('\n')
			if err != nil {
				// io.EOF and other errors
				close(ch)
				return
			}
			ch <- s
		}
	}(stdinCh)

	ticker := time.NewTicker(time.Second)
	start := time.Now()
	for {
		select {
		case <- ticker.C:
			now := time.Now()
			elapsed := now.Sub(start)
			minutes := int(elapsed.Minutes())
			seconds := int(elapsed.Seconds()) % 60
			fmt.Printf("\r%02d:%02d", minutes, seconds)
		case stdin, ok := <-stdinCh:
			if !ok {
				fmt.Println("stdin was lost")
				return false
			} else {
				stdin = strings.TrimSpace(stdin)
				if isOtherMode(currentMode, stdin) {
					return true
				} else if stdin == "q" {
					return false
				} else {
					printInputOptions(currentMode)
				}
			}
		}
	}
}

func main() {
	mode := "w"
	for {
		keepGoing := runTimer(mode)
		if !keepGoing {
			break
		}
		mode = otherMode(mode)
	}
}
