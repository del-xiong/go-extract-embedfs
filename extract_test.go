package extractfs

import (
	"embed"
	"os"
	"path"
	"testing"
)

//go:embed testFiles
var testFiles embed.FS

func TestExtractToPath(t *testing.T) {
	testDir := path.Join(os.TempDir(), "extractfstest")
	os.MkdirAll(testDir, os.ModePerm)
	info, err := ExtractToPath(&testFiles, testDir)
	defer os.RemoveAll(testDir)
	if err != nil {
		t.Fatal(err)
	}
	if info.ExtractedFiles != 3 || info.ExtractedFolders != 2 {
		t.Fatal("extract or save failed")
	}
	t.Logf("extracted %d folders and %d files", info.ExtractedFolders, info.ExtractedFiles)
}
