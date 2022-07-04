package errorpanic

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/chai2010/errors"
)

func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.Create(dstName)
	if err != nil {
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

func loadConfig() error {
	_, err := ioutil.ReadFile("path/to/file")
	if err != nil {
		return errors.Wrap(err, "read file")
	}
	// ...
	return nil
}

func setup() error {
	err := loadConfig()
	if err != nil {
		return errors.Wrap(err, "invalid config")
	}
	// ...
	return nil
}

func WrapErrors() {
	if err := setup(); err != nil {
		log.Fatal(err)
	}
	// ...
}

func RecoverPanic() {
	defer func() {
		if r := recover(); r != nil {
			// ...
		}
		// 虽然总是返回 nil，但是可以恢复异常状态
	}()
	// 警告: 以 nil 为参数抛出异常
	panic(nil)
}
