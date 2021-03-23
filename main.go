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
	var state string
	for scanner.Scan() {
		func() {
			line := scanner.Text()
			defer fmt.Println(line)

			switch state {
			case "":
				if strings.HasPrefix(line, "--- FAIL: Example") {
					state = "IN_EXAMPLE"
					return
				}
			case "IN_EXAMPLE":
				if line == "got:" {
					state = "GOT"
					return
				}
			case "GOT":
				if line == "want:" {
					state = "WANT"
					return
				}
				got += (line + "\n")
			case "WANT":
				if strings.HasPrefix(line, "--- FAIL:") || line == "FAIL" {
					// print diff
					b, err := diff("example-diff", []byte(got), []byte(want))
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println("diff:")
					fmt.Println(formatDiff(string(b)))

					// reset
					got = ""
					want = ""
					if strings.HasPrefix(line, "--- FAIL: Example") {
						// skip to the IN_EXAMPLE state if we detect the start of another failing test
						state = "IN_EXAMPLE"
					} else {
						state = ""
					}
					return
				}
				want += (line + "\n")
			}
		}()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func formatDiff(diff string) string {
	if diff == "" {
		return diff
	}
	lines := strings.Split(diff, "\n")
	if len(lines) > 3 && strings.HasPrefix(lines[0], "---") && strings.HasPrefix(lines[1], "+++") {
		diff = strings.Join(lines[2:len(lines)-1], "\n")
	}
	return diff
}

// Diff returns diff of two arrays of bytes in diff tool format.
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
