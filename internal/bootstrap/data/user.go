package data

import (
	"fmt"
	"os"

	"github.com/OpenListTeam/OpenList/v4/cmd/flags"
	"github.com/OpenListTeam/OpenList/v4/internal/db"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/internal/op"
	"github.com/OpenListTeam/OpenList/v4/pkg/utils"
	"github.com/OpenListTeam/OpenList/v4/pkg/utils/random"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func initUser() {
	admin, err := op.GetAdmin()
	adminPassword := random.String(8)
	envpass := os.Getenv("INOI_ADMIN_PASSWORD")
	if flags.Dev {
		adminPassword = "admin"
	} else if len(envpass) > 0 {
		adminPassword = envpass
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			salt := random.String(16)
			admin = &model.User{
				Username: "admin",
				Salt:     salt,
				PwdHash:  model.TwoHashPwd(adminPassword, salt),
				Role:     model.ADMIN,
				BasePath: "/",
				Authn:    "[]",
				// 0(can see hidden) - 8(webdav read) & 12(can read archives) - 14(can share)
				Permission: 0x71FF,
			}
			if err := op.CreateUser(admin); err != nil {
				panic(err)
			} else {
				// DO NOT output the password to log file. Only output to console.
				// utils.Log.Infof("Successfully created the admin user and the initial password is: %s", adminPassword)
				fmt.Printf("管理员用户创建成功，初始密码为: %s\n", adminPassword)
			}
		} else {
			utils.Log.Fatalf("[初始化用户] 获取管理员用户失败: %v", err)
		}
	}
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
			if err := db.CreateUser(guest); err != nil {
				utils.Log.Fatalf("[初始化用户] 创建访客用户失败: %v", err)
			}
		} else {
			utils.Log.Fatalf("[初始化用户] 获取访客用户失败: %v", err)
		}
	}
}