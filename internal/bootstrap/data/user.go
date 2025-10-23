package data

import (
	"fmt"
	"os"

	"github.com/OpenListTeam/OpenList/v4/cmd/flags"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/internal/op"
	"github.com/OpenListTeam/OpenList/v4/pkg/utils"
	"github.com/OpenListTeam/OpenList/v4/pkg/utils/random"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func initUser() {
	admin, err := op.GetAdmin()
	
	// 先确定密码来源（不立即赋值）
	var adminPassword string
	if flags.Dev {
		adminPassword = "admin" // 开发环境
	} else if envpass := os.Getenv("OPENLIST_ADMIN_PASSWORD"); envpass != "" {
		adminPassword = envpass // 生产环境优先用环境变量
	} else {
		adminPassword = "iNoi-PSWD" // 默认值
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			salt := random.String(16)
			admin = &model.User{
				Username:   "admin",
				Salt:       salt,
				PwdHash:    model.TwoHashPwd(adminPassword, salt),
				Role:       model.ADMIN,
				BasePath:   "/",
				Authn:      "[]",
				Permission: 0x71FF, // 权限位
			}
			if err := op.CreateUser(admin); err != nil {
				panic(err)
			}
			// 安全提示：仅输出到控制台
			fmt.Printf("管理员用户创建成功，初始密码为: %s\n", adminPassword)
		} else {
			utils.Log.Fatalf("[初始化用户] 获取管理员用户失败: %v", err)
		}
	}
	
	// Guest用户创建修正（使用op.CreateUser）
	_, err = op.GetGuest()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			salt := random.String(16)
			guest := &model.User{
				Username:   "guest",
				PwdHash:    model.TwoHashPwd("guest", salt),
				Salt:       salt,
				Role:       model.GUEST,
				BasePath:   "/",
				Permission: 0,
				Disabled:   true,
				Authn:      "[]",
			}
			if err := op.CreateUser(guest); err != nil {
				utils.Log.Fatalf("[初始化用户] 创建访客用户失败: %v", err)
			}
		} else {
			utils.Log.Fatalf("[初始化用户] 获取访客用户失败: %v", err)
		}
	}
}