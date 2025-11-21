package snowflake

import "testing"

func TestNewWorker(t *testing.T) {
	id, _ := NewWorker(WorkerID, WataCenterID).NextID()
	t.Log(id)
}
