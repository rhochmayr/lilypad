package bacalhau

import (
	"strings"
	"testing"

	"github.com/lilypad-tech/lilypad/pkg/data"
	"github.com/lilypad-tech/lilypad/pkg/system"
	"github.com/stretchr/testify/assert"
)

func TestParseBacalhauJobID(t *testing.T) {
	tests := []struct {
		name          string
		output        string
		expectedID    string
		expectedError bool
	}{
		{
			name:          "Valid output",
			output:        "jobID\n",
			expectedID:    "jobID",
			expectedError: false,
		},
		{
			name:          "Output with error",
			output:        "jobID\nerror message",
			expectedID:    "",
			expectedError: true,
		},
		{
			name:          "Empty output",
			output:        "",
			expectedID:    "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := parseBacalhauJobID([]byte(tt.output))
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedID, id)
		})
	}
}

func TestGetJobID(t *testing.T) {
	ipfsClient := &ipfs.Client{}
	executor, err := NewBacalhauExecutor(BacalhauExecutorOptions{}, ipfsClient)
	assert.NoError(t, err)

	deal := data.DealContainer{
		ID: "test-deal",
	}
	module := data.Module{
		Job: data.Job{
			Spec: data.Spec{
				Engine: data.EngineDocker,
				Docker: data.JobSpecDocker{
					Image: "ubuntu",
				},
			},
		},
	}

	system.EnsureDataDir("bacalhau-job-specs/test-deal")
	system.WriteFile("bacalhau-job-specs/test-deal/job.json", []byte(`{"Spec":{"Engine":"Docker","Docker":{"Image":"ubuntu"}}}`))

	id, err := executor.getJobID(deal, module)
	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(id, "jobID"))
}
