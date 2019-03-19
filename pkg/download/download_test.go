package download

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractIssuePublication(t *testing.T) {
	url1 := "https://pros.lekiosk.com/fr/pageproduct/1920587/1920588"
	url2 := "https://pros.lekiosk.com/fr/reader/3125124/341212312"

	res1, err1 := ExtractIssuePublication(url1)
	assert.NoError(t, err1)
	assert.Equal(t, "1920587", res1.Publication)
	assert.Equal(t, "1920588", res1.Issue)

	res2, err2 := ExtractIssuePublication(url2)
	assert.NoError(t, err2)
	assert.Equal(t, "3125124", res2.Publication)
	assert.Equal(t, "341212312", res2.Issue)

}
