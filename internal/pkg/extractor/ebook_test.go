package extractor

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var converter *EBookConverter

func initTest(t *testing.T) {
	var err error
	converter, err = NewEBookConverter()
	assert.Nil(t, err)
}

func TestExract_Invokes(t *testing.T) {
	initTest(t)
	var s []string
	converter.extractFunc = func(cmd []string) error {
		s = cmd
		return nil
	}
	err := converter.Extract("/dir/file.epub", "/dir/file.txt")
	assert.Equal(t, []string{"ebook-convert", "/dir/file.epub", "/dir/file.txt"}, s)
	assert.Nil(t, err)
}

func TestExract_Fails(t *testing.T) {
	initTest(t)
	converter.extractFunc = func(cmd []string) error {
		return errors.New("olia")
	}
	err := converter.Extract("/dir/file.epub", "/dir/file.txt")
	assert.NotNil(t, err)
}

func Test_runCmd(t *testing.T) {
	type args struct {
		cmdArr []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "success", args: args{cmdArr: []string{"ls", "-la"}}, wantErr: false},
		{name: "fails", args: args{cmdArr: []string{"missing_cmd"}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := runCmd(tt.args.cmdArr); (err != nil) != tt.wantErr {
				t.Errorf("runCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
