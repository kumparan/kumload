package kumload

import (
	"testing"
)

func TestParseConfig(t *testing.T) {
	tests := []struct {
		Name         string
		RelativePath string
		ExpectError  bool
	}{
		{
			Name:         "FileExist",
			RelativePath: "./config.yml.example",
			ExpectError:  false,
		},
		{
			Name:         "FileNotExist",
			RelativePath: "./configg.yml",
			ExpectError:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			_, err := ParseConfig(test.RelativePath)
			if err != nil && !test.ExpectError {
				t.Errorf("expected error: %t, got: %s", test.ExpectError, err)
			}
		})
	}
}
