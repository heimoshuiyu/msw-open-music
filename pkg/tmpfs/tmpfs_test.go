package tmpfs

import "testing"

func TestTmpfs(t *testing.T) {
	t.Log("Starting ...")
	tmpfs := NewTmpfs()
	tmpfs.FileLifeTime = 1
	tmpfs.Record("/tmp/testfile")
	t.Log(tmpfs.record)
	tmpfs.wg.Wait()
}
