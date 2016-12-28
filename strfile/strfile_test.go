package strfile

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func run(input string) bytes.Buffer {
	inFile, _ := ioutil.TempFile("", "in")
	defer os.Remove(inFile.Name())
	fmt.Fprint(inFile, input)

	outFile, _ := ioutil.TempFile("", "out")
	defer os.Remove(outFile.Name())

	cmd := exec.Command("strfile", "-s", inFile.Name(), outFile.Name())
	cmd.Run()

	var buf bytes.Buffer
	io.Copy(&buf, outFile)

	return buf
}

func TestStrfileWithTwoLines(t *testing.T) {
	var buf bytes.Buffer
	input := `a
%
b`

	err := Strfile(strings.NewReader(input), &buf)
	if err != nil {
		t.Fatal(err)
	}

	exp := run(input)

	if buf.String() != exp.String() {
		t.Fatalf("Expected %q\n"+
			"            to equal %q", buf.String(), exp.String())
	}
}

func TestStrfileWithThreeLines(t *testing.T) {
	var buf bytes.Buffer
	input := `abc
%
bc
%
c`

	err := Strfile(strings.NewReader(input), &buf)
	if err != nil {
		t.Fatal(err)
	}

	exp := run(input)

	if buf.String() != exp.String() {
		t.Fatalf("Expected %q\n"+
			"            to equal %q", buf.String(), exp.String())
	}
}

func TestStrfileWithDelimInLines(t *testing.T) {
	var buf bytes.Buffer
	input := `%%
%
%%
%
%%`

	err := Strfile(strings.NewReader(input), &buf)
	if err != nil {
		t.Fatal(err)
	}

	exp := run(input)

	if buf.String() != exp.String() {
		t.Fatalf("Expected %q\n"+
			"            to equal %q", buf.String(), exp.String())
	}
}

func TestStrfileWithManyLines(t *testing.T) {
	var buf bytes.Buffer

	input, _ := ioutil.ReadFile("/usr/share/games/fortunes/fortunes")
	exp, _ := ioutil.ReadFile("/usr/share/games/fortunes/fortunes.dat")

	err := Strfile(bytes.NewReader(input), &buf)
	if err != nil {
		t.Fatal(err)
	}

	if buf.String() != string(exp) {
		t.Fatalf("Expected %q\n"+
			"            to equal %q", buf.String(), string(exp))
	}
}
