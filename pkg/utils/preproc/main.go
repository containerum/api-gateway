package preproc

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"regexp"
)

var (
	includeRegexp, _ = regexp.Compile("^#include \"(.+)\"$")
)

func Preprocess(path string) (io.Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var writer bytes.Buffer

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		writer.Write(scan.Bytes())
		writer.WriteRune('\n')
	}

	r := bytes.NewReader(writer.Bytes())
	return r, nil
}
