package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// Match the following:
	//	--- FAIL: ExampleParagraphComponent_main (0.00s)
	//	got:
	//	<got>
	//	want:
	//	<want>
	//	FAIL

	var got, want string
	currentState := "UNKNOWN"
	for scanner.Scan() {
		line := scanner.Text()
		nextState := transitionState(currentState, line)

		// keep appending lines if there's no state change
		if currentState == nextState {
			switch currentState {
			case "GOT":
				got += (line + "\n")
			case "WANT":
				want += (line + "\n")
			}
		}

		// print diff if we're transitioning out of the WANT state
		if currentState == "WANT" && currentState != nextState {
			printDiff(got, want)
			got = ""
			want = ""
		}

		currentState = nextState
		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func transitionState(state, line string) string {
	switch state {
	case "UNKNOWN":
		if strings.HasPrefix(line, "--- FAIL: Example") {
			state = "IN_EXAMPLE"
		}
	case "IN_EXAMPLE":
		if line == "got:" {
			state = "GOT"
		}
	case "GOT":
		if line == "want:" {
			state = "WANT"
		}
	case "WANT":
		if strings.HasPrefix(line, "--- FAIL: Example") {
			state = "IN_EXAMPLE"
		} else if strings.HasPrefix(line, "--- FAIL:") || line == "FAIL" {
			state = "UNKNOWN"
		}
	}
	return state
}

func printDiff(got, want string) {
	b, err := diff("example-diff", []byte(got), []byte(want))
	if err != nil {
		log.Fatal(err)
	}

	diff := string(b)
	lines := strings.Split(diff, "\n")
	if len(lines) > 3 && strings.HasPrefix(lines[0], "---") && strings.HasPrefix(lines[1], "+++") {
		diff = strings.Join(lines[2:len(lines)-1], "\n")
	}
	fmt.Println("diff:")
	fmt.Println(diff)
}

// Diff returns diff of two arrays of bytes in diff tool format. Lifted from
// src/cmd/internal/diff in the Go stdlib.
func diff(prefix string, b1, b2 []byte) ([]byte, error) {
	f1, err := writeTempFile(prefix, b1)
	if err != nil {
		return nil, err
	}
	defer os.Remove(f1)

	f2, err := writeTempFile(prefix, b2)
	if err != nil {
		return nil, err
	}
	defer os.Remove(f2)

	data, err := exec.Command("diff", "-u", f1, f2).CombinedOutput()
	if len(data) > 0 {
		// diff exits with a non-zero status when the files don't match.
		// Ignore that failure as long as we get output.
		err = nil
	}
	return data, err
}

func writeTempFile(prefix string, data []byte) (string, error) {
	file, err := ioutil.TempFile("", prefix)
	if err != nil {
		return "", err
	}
	_, err = file.Write(data)
	if err1 := file.Close(); err == nil {
		err = err1
	}
	if err != nil {
		os.Remove(file.Name())
		return "", err
	}
	return file.Name(), nil
}
