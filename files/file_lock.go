package files

import (
	"fmt"
	"os"
	// "sync"
	"syscall"

)

//文件锁
type FileLock struct {
    dir string
    f   *os.File
}

func NewFileLock(dir string) *FileLock {
    return &FileLock{
        dir: dir,
    }
}

//加锁
func (l *FileLock) Lock() error {
    f, err := os.OpenFile(l.dir,os.O_CREATE|os.O_WRONLY,0644)
    if err != nil {
        return err
    }
    l.f = f
    err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
    if err != nil {
        return fmt.Errorf("cannot flock directory %s - %s", l.dir, err)
    }
    return nil
}

//释放锁
func (l *FileLock) Unlock() error {
    defer l.f.Close()
    return syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
}