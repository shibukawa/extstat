package extstat

import (
	"io/ioutil"
	"os"
	"runtime"
	"testing"
	"time"
)

var filePath string = "test.txt"

func TestNew(t *testing.T) {
	err := ioutil.WriteFile(filePath, []byte("hello"), 0666)
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(filePath)
	now := time.Now()

	btimeBefore := now.Add(+time.Second)
	btimeAfter := now.Add(-time.Second)
	mtimeBefore := now.Add(time.Second * 4)
	mtimeAfter := now.Add(time.Second * 2)
	atimeBefore := now.Add(time.Second * 7)
	atimeAfter := now.Add(time.Second * 5)

	t.Log("wait for changing mtime...")
	time.Sleep(time.Second * 3)
	fileForModify, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fileForModify.Close()
		t.Error(err)
		return
	}
	fileForModify.Write([]byte(" world"))
	fileForModify.Close()

	t.Log("wait for changing atime...")
	time.Sleep(time.Second * 3)
	content, err := ioutil.ReadFile(filePath)
	t.Log(string(content))
	if err != nil {
		t.Error(err)
		return
	}

	stat, err := os.Stat(filePath)
	if err != nil {
		t.Error(err)
		return
	}

	extStat := New(stat)

	// plan9 doesn't have correct ctime attribute
	if runtime.GOOS != "plan9" {
		if !extStat.BirthTime.Before(btimeBefore) || !extStat.BirthTime.After(btimeAfter) {
			t.Error("btime is wrong:", btimeAfter, "<", extStat.BirthTime, "<", btimeBefore, "  now: ", now)
		}
	}
	if !extStat.ModTime.Before(mtimeBefore) || !extStat.ModTime.After(mtimeAfter) {
		t.Error("mtime is wrong:", mtimeAfter, "<", extStat.ModTime, "<", mtimeBefore, "  now: ", now)
	}
	if !extStat.ModTime.Before(mtimeBefore) || !extStat.ChangeTime.After(mtimeAfter) {
		t.Error("mtime is wrong:", mtimeAfter, "<", extStat.ModTime, "<", mtimeBefore, "  now: ", now)
	}
	if !extStat.AccessTime.Before(atimeBefore) || !extStat.AccessTime.After(atimeAfter) {
		t.Error("atime is wrong:", atimeAfter, "<", extStat.AccessTime, "<", atimeBefore, "  now: ", now)
	}
}
