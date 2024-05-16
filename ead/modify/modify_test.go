package modify

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func failOnError(t *testing.T, err error, label string) {
	if err != nil {
		t.Errorf("%s: %s", label, err)
		t.FailNow()
	}
}

func TestFABifyEAD(t *testing.T) {
	t.Run("Modify EAD For Discovery System (FAB): creator, location", func(t *testing.T) {

		testFixturePath := filepath.Join(".", "testdata")
		testTmpDirPath := filepath.Join(testFixturePath, "tmp")

		sourceFile := filepath.Join(testFixturePath, "modify-input.xml")
		referenceFile := filepath.Join(testFixturePath, "modify-expected.xml")

		EADXML, err := os.ReadFile(sourceFile)
		failOnError(t, err, "Unexpected error")

		doc, errors := FABifyEAD(EADXML)
		if len(errors) != 0 {
			failOnError(t, fmt.Errorf("%s", strings.Join(errors, "\n")), "problem modifying EAD")
		}

		got := []byte(doc)
		want, err := os.ReadFile(referenceFile)
		failOnError(t, err, "Unexpected error reading reference file")

		if !bytes.Equal(want, got) {
			errTmpFile := filepath.Join(testTmpDirPath, "ERR-modify-ead.xml")
			err = os.WriteFile(errTmpFile, got, 0644)
			failOnError(t, err, fmt.Sprintf("Unexpected error writing %s", errTmpFile))

			errMsg := fmt.Sprintf("The modified EAD does not match the reference file.\ndiff %s %s", errTmpFile, referenceFile)
			t.Errorf(errMsg)
		}
	})

	t.Run("Modify EAD For Discovery System (FAB): unitid @aspace_uri", func(t *testing.T) {

		testFixturePath := filepath.Join(".", "testdata")
		testTmpDirPath := filepath.Join(testFixturePath, "tmp")

		sourceFile := filepath.Join(testFixturePath, "unitid-aspace_uri-input.xml")
		referenceFile := filepath.Join(testFixturePath, "unitid-aspace_uri-expected.xml")

		EADXML, err := os.ReadFile(sourceFile)
		failOnError(t, err, "Unexpected error")

		doc, errors := FABifyEAD(EADXML)
		if len(errors) != 0 {
			failOnError(t, fmt.Errorf("%s", strings.Join(errors, "\n")), "problem modifying EAD")
		}

		got := []byte(doc)
		want, err := os.ReadFile(referenceFile)
		failOnError(t, err, "Unexpected error reading reference file")

		if !bytes.Equal(want, got) {
			errTmpFile := filepath.Join(testTmpDirPath, "ERR-unitid-aspace_uri-ead.xml")
			err = os.WriteFile(errTmpFile, got, 0644)
			failOnError(t, err, fmt.Sprintf("Unexpected error writing %s", errTmpFile))

			errMsg := fmt.Sprintf("The modified EAD does not match the reference file.\ndiff %s %s", errTmpFile, referenceFile)
			t.Errorf(errMsg)
		}
	})
}
