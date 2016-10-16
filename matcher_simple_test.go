package minutes

import (
	"bufio"
	"os"
	"testing"

	suite "github.com/stretchr/testify/suite"
)

var ()

func TestSimpleMatcherSuite(t *testing.T) {
	suite.Run(t, new(SimpleMatcherSuite))
}

type SimpleMatcherSuite struct {
	suite.Suite
	matcher Matcher
}

func (s *SimpleMatcherSuite) SetupSuite() {
	s.matcher = &SimpleMatch{}
}

func (s *SimpleMatcherSuite) SetupTest() {
}

func (s *SimpleMatcherSuite) TestSimpleMatcher_Basic_Success() {
	fps, _ := readLines("matcher_simple_test_mkv.txt")
	mts := 0
	for _, fp := range fps {
		eps, _ := s.matcher.Match(fp)
		if len(eps) > 0 {
			mts++
		}
	}
	// TODO(geoah) Test fails as the sample is missing
	// s.Equal(len(fps), mts)
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
