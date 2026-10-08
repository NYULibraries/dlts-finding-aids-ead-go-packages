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
	t.Run("duplicate container parent IDs return an error without panicking", func(t *testing.T) {
		input := []byte(`<ead xmlns="urn:isbn:1-931666-22-9">
			<archdesc>
				<did><container id="root">Box 1</container></did>
				<dsc><c><did>
					<container id="child1" parent="root">Folder 1</container>
					<container id="child2" parent="root">Folder 2</container>
				</did></c></dsc>
			</archdesc>
		</ead>`)

		defer func() {
			if recovered := recover(); recovered != nil {
				t.Fatalf("FABifyEAD panicked for duplicate container parent IDs: %v", recovered)
			}
		}()

		_, errors := FABifyEAD(input)
		if !strings.Contains(strings.Join(errors, "\n"), "multiple subcontainers with the same parentID: root") {
			t.Fatalf("expected duplicate-parent error, got %v", errors)
		}
	})

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
