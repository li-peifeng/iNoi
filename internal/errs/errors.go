package errs

import (
	"errors"
	"fmt"

	pkgerr "github.com/pkg/errors"
)

var (
	NotImplement = errors.New("功能未实现")
	NotSupport   = errors.New("不支持此操作")
	RelativePath = errors.New("不允许使用相对路径")

	UploadNotSupported = errors.New("不支持上传功能")
	MetaNotFound       = errors.New("元数据不存在")
	StorageNotFound    = errors.New("存储不存在")
	StorageNotInit     = errors.New("存储未初始化")
	StreamIncomplete   = errors.New("上传/下载流不完整，可能是网络问题")
	StreamPeekFail     = errors.New("流预览失败")

	UnknownArchiveFormat      = errors.New("未知的压缩文件格式")
	WrongArchivePassword      = errors.New("压缩包密码错误")
	DriverExtractNotSupported = errors.New("驱动不支持解压操作")

	WrongShareCode  = errors.New("分享码错误")
	InvalidSharing  = errors.New("分享无效")
	SharingNotFound = errors.New("分享不存在")
)

// NewErr wrap constant error with an extra message
// use errors.Is(err1, StorageNotFound) to check if err belongs to any internal error
func NewErr(err error, format string, a ...any) error {
	return fmt.Errorf("%w; %s", err, fmt.Sprintf(format, a...))
}

func IsNotFoundError(err error) bool {
	return errors.Is(pkgerr.Cause(err), ObjectNotFound) || errors.Is(pkgerr.Cause(err), StorageNotFound)
}

func IsNotSupportError(err error) bool {
	return errors.Is(pkgerr.Cause(err), NotSupport)
}
func IsNotImplementError(err error) bool {
	return errors.Is(pkgerr.Cause(err), NotImplement)
}
