package v1

import (
	"errors"
	"time"

	"github.com/1Panel-dev/1Panel/backend/app/api/v1/helper"
	"github.com/1Panel-dev/1Panel/backend/app/dto"
	"github.com/1Panel-dev/1Panel/backend/constant"
	"github.com/1Panel-dev/1Panel/backend/global"
	"github.com/1Panel-dev/1Panel/backend/utils/mfa"
	"github.com/1Panel-dev/1Panel/backend/utils/ntp"
	"github.com/gin-gonic/gin"
)

// Load system setting info
// @Tags System Setting
// @Summary Load system setting info
// @Description 加载系统配置信息
// @Success 200 {object} dto.SettingInfo
// @Security ApiKeyAuth
// @Router /settings/search [post]
func (b *BaseApi) GetSettingInfo(c *gin.Context) {
	setting, err := settingService.GetSettingInfo()
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, setting)
}

// Load daemon.json path
// @Tags System Setting
// @Summary Load daemon.json path
// @Description 加载 docker 配置路径
// @Success 200 {string} path
// @Security ApiKeyAuth
// @Router /settings/daemonjson [get]
func (b *BaseApi) GetDaemonjson(c *gin.Context) {
	value, err := settingService.GetDaemonjson()
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, value)
}

// Update system setting
// @Tags System Setting
// @Summary Update system setting
// @Description 更新系统配置
// @Accept json
// @Param request body dto.SettingUpdate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /settings/update [post]
// @x-panel-log {"bodyKeys":["key","value"],"paramKeys":[],"BeforeFuntions":[],"formatZH":"修改系统配置 [key] => [value]","formatEN":"update system setting [key] => [value]"}
func (b *BaseApi) UpdateSetting(c *gin.Context) {
	var req dto.SettingUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
		return
	}
	if err := global.VALID.Struct(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
		return
	}

	if err := settingService.Update(c, req.Key, req.Value); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// Update system password
// @Tags System Setting
// @Summary Update system password
// @Description 更新系统登录密码
// @Accept json
// @Param request body dto.PasswordUpdate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /settings/password/update [post]
// @x-panel-log {"bodyKeys":[],"paramKeys":[],"BeforeFuntions":[],"formatZH":"修改系统密码","formatEN":"update system password"}
func (b *BaseApi) UpdatePassword(c *gin.Context) {
	var req dto.PasswordUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
		return
	}
	if err := global.VALID.Struct(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
		return
	}

	if err := settingService.UpdatePassword(c, req.OldPassword, req.NewPassword); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// Reset system password expired
// @Tags System Setting
// @Summary Reset system password expired
// @Description 重置过期系统登录密码
// @Accept json
// @Param request body dto.PasswordUpdate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /settings/expired/handle [post]
// @x-panel-log {"bodyKeys":[],"paramKeys":[],"BeforeFuntions":[],"formatZH":"重置过期密码","formatEN":"reset an expired Password"}
func (b *BaseApi) HandlePasswordExpired(c *gin.Context) {
	var req dto.PasswordUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
		return
	}
	if err := global.VALID.Struct(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
		return
	}

	if err := settingService.HandlePasswordExpired(c, req.OldPassword, req.NewPassword); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// Sync system time
// @Tags System Setting
// @Summary Sync system time
// @Description 系统时间同步
// @Success 200 {string} ntime
// @Security ApiKeyAuth
// @Router /settings/time/sync [post]
// @x-panel-log {"bodyKeys":[],"paramKeys":[],"BeforeFuntions":[],"formatZH":"系统时间同步","formatEN":"sync system time"}
func (b *BaseApi) SyncTime(c *gin.Context) {
	ntime, err := ntp.Getremotetime()
	if err != nil {
		helper.SuccessWithData(c, time.Now().Format("2006-01-02 15:04:05 MST -0700"))
		return
	}

	ts := ntime.Format("2006-01-02 15:04:05")
	if err := ntp.UpdateSystemDate(ts); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, ntime.Format("2006-01-02 15:04:05 MST -0700"))
}

// Clean monitor datas
// @Tags System Setting
// @Summary Clean monitor datas
// @Description 清空监控数据
// @Success 200
// @Security ApiKeyAuth
// @Router /settings/monitor/clean [post]
// @x-panel-log {"bodyKeys":[],"paramKeys":[],"BeforeFuntions":[],"formatZH":"清空监控数据","formatEN":"clean monitor datas"}
func (b *BaseApi) CleanMonitor(c *gin.Context) {
	if err := global.DB.Exec("DELETE FROM monitor_bases").Error; err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	if err := global.DB.Exec("DELETE FROM monitor_ios").Error; err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	if err := global.DB.Exec("DELETE FROM monitor_networks").Error; err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, nil)
}

// Load mfa info
// @Tags System Setting
// @Summary Load mfa info
// @Description 获取 mfa 信息
// @Success 200 {object} mfa.Otp
// @Security ApiKeyAuth
// @Router /settings/mfa [get]
func (b *BaseApi) GetMFA(c *gin.Context) {
	otp, err := mfa.GetOtp("admin")
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, otp)
}

// Bind mfa
// @Tags System Setting
// @Summary Bind mfa
// @Description Mfa 绑定
// @Accept json
// @Param request body dto.MfaCredential true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /settings/mfa/bind [post]
// @x-panel-log {"bodyKeys":[],"paramKeys":[],"BeforeFuntions":[],"formatZH":"mfa 绑定","formatEN":"bind mfa"}
func (b *BaseApi) MFABind(c *gin.Context) {
	var req dto.MfaCredential
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
		return
	}
	success := mfa.ValidCode(req.Code, req.Secret)
	if !success {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, errors.New("code is not valid"))
		return
	}

	if err := settingService.Update(c, "MFAStatus", "enable"); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	if err := settingService.Update(c, "MFASecret", req.Secret); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, nil)
}
