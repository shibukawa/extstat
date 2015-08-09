package extstat

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
	"runtime"
)

func TestNew(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Error(err)
		return
	}
	filePath := f.Name()
	f.Write([]byte("test"))
	defer os.Remove(filePath)
	now := time.Now()

	ctimeBefore := now.Add(+time.Second)
	ctimeAfter :=  now.Add(-time.Second)
	mtimeBefore := now.Add(time.Second*4)
	mtimeAfter := now.Add(time.Second*2)
	atimeBefore := now.Add(time.Second*7)
	atimeAfter := now.Add(time.Second*5)
	f.Close()

	t.Log("wait for changing mtime...")
	time.Sleep(time.Second * 3)
	fileForModify, err := os.OpenFile(filePath, os.O_APPEND, 0666)
	if err != nil {
		fileForModify.Close()
		t.Error(err)
		return
	}
	fileForModify.Write([]byte("hello world"))
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
		if !extStat.CreatedTime.Before(ctimeBefore) || !extStat.CreatedTime.After(ctimeAfter) {
			t.Error("ctime is wrong:", ctimeAfter, "<", extStat.CreatedTime, "<", ctimeBefore)
		}
	}
	if !extStat.ModTime.Before(mtimeBefore) || !extStat.ModTime.After(mtimeAfter) {
		t.Error("mtime is wrong:", mtimeAfter, "<", extStat.ModTime, "<", mtimeBefore)
	}
	if !extStat.AccessTime.Before(atimeBefore) || !extStat.AccessTime.After(atimeAfter) {
		t.Error("atime is wrong:", atimeAfter, "<", extStat.AccessTime, "<", atimeBefore)
	}
}
