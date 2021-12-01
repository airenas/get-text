package extractor

import (
	"bytes"
	"os/exec"

	"github.com/pkg/errors"
)

//EBookConverter wrapper for text extraction
type EBookConverter struct {
	extractFunc func([]string) error
}

//NewEBookConverter return new extractor instance
func NewEBookConverter() (*EBookConverter, error) {
	res := EBookConverter{}
	res.extractFunc = runCmd
	return &res, nil
}

//Extract returns name of new written txt
func (e *EBookConverter) Extract(nameIn, nameOut string) error {
	params := []string{"ebook-convert", nameIn, nameOut}
	err := e.extractFunc(params)
	if err != nil {
		return err
	}
	return nil
}

func runCmd(cmdArr []string) error {
	cmd := exec.Command(cmdArr[0], cmdArr[1:]...)
	var outputBuffer bytes.Buffer
	cmd.Stdout = &outputBuffer
	cmd.Stderr = &outputBuffer
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "output: "+outputBuffer.String())
	}
	return nil
}
