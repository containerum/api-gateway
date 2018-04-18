package preproc

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var (
	includeRegexp, _         = regexp.Compile("^#include \"(.+)\"$")
	includeFilenameRegexp, _ = regexp.Compile("\"(.+)\"$")
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
		msgb := scan.Bytes()
		if includeRegexp.Match(msgb) {
			fnameb := includeFilenameRegexp.Find(msgb)
			fname := strings.Trim(string(fnameb), "\"")
			file, err := ioutil.ReadFile(fname)
			if err != nil {
				return nil, err
			}
			writer.Write(file)
		} else {
			writer.Write(msgb)
		}
		writer.WriteRune('\n')
	}

	r := bytes.NewReader(writer.Bytes())
	return r, nil
}
