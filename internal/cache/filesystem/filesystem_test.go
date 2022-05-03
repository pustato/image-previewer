package filesystem

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const basePath = "/tmp/test"

type FilesystemTestSuite struct {
	suite.Suite
	basePath string
}

func (ts *FilesystemTestSuite) SetupTest() {
	ts.basePath = basePath
}

func (ts *FilesystemTestSuite) TearDownTest() {
	os.RemoveAll(ts.basePath)
}

func (ts *FilesystemTestSuite) TestConstructor() {
	ts.T().Run("error", func(t *testing.T) {
		_, err := NewDiscFilesystem("/dev/not_exists_and_not_writable")
		require.Error(ts.T(), err)
	})

	ts.T().Run("success", func(t *testing.T) {
		_, err := NewDiscFilesystem(ts.basePath)
		require.NoError(ts.T(), err)

		_, err = os.Stat(ts.basePath)
		require.True(t, !os.IsNotExist(err))
	})
}

func (ts *FilesystemTestSuite) TestComplex() {
	fs, _ := NewDiscFilesystem(ts.basePath)
	name := "file.txt"
	content := []byte("content")

	require.NoError(ts.T(), fs.WriteFile(name, []byte("content")))

	actual, err := fs.ReadFile(name)
	require.NoError(ts.T(), err)
	require.EqualValues(ts.T(), content, actual)

	require.NoError(ts.T(), fs.RemoveFile(name))

	_, err = fs.ReadFile(name)
	require.Error(ts.T(), err)
	require.ErrorIs(ts.T(), err, ErrFileNotExists)
}

func TestFilesystemTestSuite(t *testing.T) {
	suite.Run(t, new(FilesystemTestSuite))
}
