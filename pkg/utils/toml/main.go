package preproc

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
)

//TomlFile interface for model with Validate function
type TomlFile interface {
	Validate() []error
}

var (
	includeRegexp, _         = regexp.Compile("^#include \"(.+)\"$")
	includeFilenameRegexp, _ = regexp.Compile("\"(.+)\"$")
)

//ReadToml function for reading toml files, run validation and preproccesing
func ReadToml(file string, out TomlFile) (err error) {
	var reader io.Reader
	if reader, err = preprocess(file); err != nil {
		return
	}
	if _, err = toml.DecodeReader(reader, out); err != nil {
		return
	}
	if errs := out.Validate(); errs != nil {
		var errWritter bytes.Buffer
		for _, e := range errs {
			errWritter.WriteString(e.Error())
		}
		return errors.New(errWritter.String())
	}
	return nil
}

func preprocess(path string) (reader io.Reader, err error) {
	var file *os.File
	var writer bytes.Buffer
	if file, err = os.Open(path); err != nil {
		return nil, err
	}
	defer file.Close()
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		msgb := scan.Bytes()
		if includeRegexp.Match(msgb) {
			var fbytes []byte
			fnameb := includeFilenameRegexp.Find(msgb)
			fname := strings.Trim(string(fnameb), "\"")
			if fbytes, err = ioutil.ReadFile(fname); err != nil {
				return nil, err
			}
			writer.Write(fbytes)
		} else {
			writer.Write(msgb)
		}
		writer.WriteRune('\n')
	}
	reader = bytes.NewReader(writer.Bytes())
	return reader, nil
}
