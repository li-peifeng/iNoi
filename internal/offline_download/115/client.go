package _115

import (
	"context"
	"fmt"

	"github.com/OpenListTeam/OpenList/v4/internal/conf"
	"github.com/OpenListTeam/OpenList/v4/internal/setting"

	_115 "github.com/OpenListTeam/OpenList/v4/drivers/115"
	"github.com/OpenListTeam/OpenList/v4/internal/errs"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/internal/offline_download/tool"
	"github.com/OpenListTeam/OpenList/v4/internal/op"
)

type Cloud115 struct {
	refreshTaskCache bool
}

func (p *Cloud115) Name() string {
	return "115 Cloud"
}

func (p *Cloud115) Items() []model.SettingItem {
	return nil
}

func (p *Cloud115) Run(task *tool.DownloadTask) error {
	return errs.NotSupport
}

func (p *Cloud115) Init() (string, error) {
	p.refreshTaskCache = false
	return "完成", nil
}

func (p *Cloud115) IsReady() bool {
	tempDir := setting.GetStr(conf.Pan115TempDir)
	if tempDir == "" {
		return false
	}
	storage, _, err := op.GetStorageAndActualPath(tempDir)
	if err != nil {
		return false
	}
	if _, ok := storage.(*_115.Pan115); !ok {
		return false
	}
	return true
}

func (p *Cloud115) AddURL(args *tool.AddUrlArgs) (string, error) {
	// 添加新任务刷新缓存
	p.refreshTaskCache = true
	storage, actualPath, err := op.GetStorageAndActualPath(args.TempDir)
	if err != nil {
		return "", err
	}
	driver115, ok := storage.(*_115.Pan115)
	if !ok {
		return "", fmt.Errorf("不支持此存储，仅支持 115 Cloud")
	}

	ctx := context.Background()

	if err := op.MakeDir(ctx, storage, actualPath); err != nil {
		return "", err
	}

	parentDir, err := op.GetUnwrap(ctx, storage, actualPath)
	if err != nil {
		return "", err
	}

	hashs, err := driver115.OfflineDownload(ctx, []string{args.Url}, parentDir)
	if err != nil || len(hashs) < 1 {
		return "", fmt.Errorf("添加离线下载任务失败: %w", err)
	}

	return hashs[0], nil
}

func (p *Cloud115) Remove(task *tool.DownloadTask) error {
	storage, _, err := op.GetStorageAndActualPath(task.TempDir)
	if err != nil {
		return err
	}
	driver115, ok := storage.(*_115.Pan115)
	if !ok {
		return fmt.Errorf("不支持此存储，仅支持 115 Cloud")
	}

	ctx := context.Background()
	if err := driver115.DeleteOfflineTasks(ctx, []string{task.GID}, false); err != nil {
		return err
	}
	return nil
}

func (p *Cloud115) Status(task *tool.DownloadTask) (*tool.Status, error) {
	storage, _, err := op.GetStorageAndActualPath(task.TempDir)
	if err != nil {
		return nil, err
	}
	driver115, ok := storage.(*_115.Pan115)
	if !ok {
		return nil, fmt.Errorf("不支持此存储，仅支持 115 Cloud")
	}

	tasks, err := driver115.OfflineList(context.Background())
	if err != nil {
		return nil, err
	}

	s := &tool.Status{
		Progress:  0,
		NewGID:    "",
		Completed: false,
		Status:    "任务已被删除",
		Err:       nil,
	}
	for _, t := range tasks {
		if t.InfoHash == task.GID {
			s.Progress = t.Percent
			s.Status = t.GetStatus()
			s.Completed = t.IsDone()
			s.TotalBytes = t.Size
			if t.IsFailed() {
				s.Err = fmt.Errorf(t.GetStatus())
			}
			return s, nil
		}
	}
	s.Err = fmt.Errorf("任务已被删除")
	return nil, nil
}

var _ tool.Tool = (*Cloud115)(nil)

func init() {
	tool.Tools.Add(&Cloud115{})
}
