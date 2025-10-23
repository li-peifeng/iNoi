package errs

import (
	"errors"
	pkgerr "github.com/pkg/errors"
	"testing"
)

func TestErrs(t *testing.T) {

	err1 := NewErr(StorageNotFound, "请先创建存储空间")
	t.Logf("err1: %s", err1)
	if !errors.Is(err1, StorageNotFound) {
		t.Errorf("失败, 预期 %s 实际 %s", err1, StorageNotFound)
	}
	if !errors.Is(pkgerr.Cause(err1), StorageNotFound) {
		t.Errorf("失败, 预期 %s 实际 %s", err1, StorageNotFound)
	}
	err2 := pkgerr.WithMessage(err1, "获取存储空间失败")
	t.Logf("err2: %s", err2)
	if !errors.Is(err2, StorageNotFound) {
		t.Errorf("失败, 预期 %s 实际 %s", err2, StorageNotFound)
	}
	if !errors.Is(pkgerr.Cause(err2), StorageNotFound) {
		t.Errorf("失败, 预期 %s 实际 %s", err2, StorageNotFound)
	}
}
