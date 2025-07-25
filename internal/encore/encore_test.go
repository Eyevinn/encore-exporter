package encore

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestDeserializeEncoreResponse(t *testing.T) {
	t.Skip("Test data not pushed yet, needs cleaning")
	is := is.New(t)
	testData, err := os.ReadFile("../test_data/json/encoreJobResult.json")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}
	var response encoreResponse
	err = json.Unmarshal(testData, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	jobs := response.Jobs()
	is.Equal(len(jobs), 4)
}
