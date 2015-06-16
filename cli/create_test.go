package cli_test

import (
	"fmt"
	"os"
	"testing"
)

func TestCreate(t *testing.T) {
	ctx, tempDir := NewTestContext(t)
	defer cleanTestDir(t, tempDir)

	// create the pipeline
	if err := ctx.Create("foobar"); err != nil {
		t.Errorf("TestCreate: '%v'", err)
	}

	project := fmt.Sprintf("%s/foobar", ctx.PipeDir)
	// check 1. directory exists
	if _, err := os.Stat(project); err != nil {
		t.Errorf("TestCreate: directory didn't get created, '%v'", err)
	}

	// check 2. it has git initialized (check for .git directory)
	if _, err := os.Stat(fmt.Sprintf("%s/.git", project)); err != nil {
		t.Errorf("TestCreate: directory did not have git initialized, '%v'", err)
	}
}

// when passing an empty string, should give an error
func TestEmptyCreate(t *testing.T) {
	ctx, tempDir := NewTestContext(t)
	defer cleanTestDir(t, tempDir)

	// create the pipeline
	if err := ctx.Create(""); err == nil || err.Error() != "Cannot create a pipeline with no name." {
		t.Errorf("TestEmptyCreate did not fail.")
	}
}

// when creating an existing pipeline, should give an error
func TestExistingCreate(t *testing.T) {
	ctx, tempDir := NewTestContext(t)
	defer cleanTestDir(t, tempDir)

	// create the pipeline
	if err := ctx.Create("foobar"); err != nil {
		t.Errorf("TestExistingCreate: '%v'", err)
	}

	// create it again! this time, it should fail
	if err := ctx.Create("foobar"); err == nil || err.Error() != "Pipeline already exists." {
		t.Errorf("TestingExistingCreate did not detect already existing pipeline.")
	}

}

// path, err := PipelinePath(name)
// if err != nil {
// 	return "", err
// }
// // make sure file exists and we have permissions, etc.
// if _, err := os.Stat(path); err != nil {
// 	return "", err
// }
// return path, nil
