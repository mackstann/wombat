package main


import "bufio"
import "fmt"
import "os"
import "strings"
import "time"

type TimeModeID string

const (
	WorkID TimeModeID = "w"
	BreakID = "b"
)

type TimeMode struct {
	id TimeModeID
	message string
}

var (
	Work = TimeMode{id: WorkID, message: "Work!"}
	Break = TimeMode{id: BreakID, message: "Take a break."}
)

func (mode TimeMode) oppositeMode() TimeMode {
	if mode == Work {
		return Break
	}
	return Work
}

func printInputOptions(currentMode TimeMode) {
	fmt.Println(" Options:")
	if currentMode == Break {
		fmt.Println("  [w]ork")
	} else if currentMode == Work {
		fmt.Println("  [b]reak")
	}
	fmt.Println("  [q]uit")
}

func runTimer(currentMode TimeMode) bool {
	fmt.Println(currentMode.message)

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
				if stdin == string(currentMode.oppositeMode().id) {
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
	mode := Work
	for {
		keepGoing := runTimer(mode)
		if !keepGoing {
			break
		}
		mode = mode.oppositeMode()
	}
}
