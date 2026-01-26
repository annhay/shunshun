package handler

import (
	"net/http"
	"shunshun/internal/api-gateway/request"
	"shunshun/internal/pkg/global"
	"shunshun/internal/proto"

	"github.com/gin-gonic/gin"
)

// NewDriver 司机认证
//
// 参数:
//   - c *gin.Context: Gin上下文
//
// 处理逻辑:
//  1. 绑定请求参数
//  2. 从上下文获取用户ID
//  3. 调用司机服务进行认证
//  4. 返回认证结果
func NewDriver(c *gin.Context) {
	var form request.NewDriver
	// 根据 Content-Type Header 推断使用哪个绑定器。
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := global.DriverClient.NewDriver(c, &proto.NewDriverReq{
		UserId:                  int64(c.GetUint64("userId")),
		RealName:                form.RealName,
		IdCardNo:                form.IdCardNo,
		IdCardFrontUrl:          form.IdCardFrontUrl,
		IdCardBackUrl:           form.IdCardBackUrl,
		IdCardExpireTime:        form.IdCardExpireTime,
		DriverLicenseNo:         form.DriverLicenseNo,
		DriverLicenseUrl:        form.DriverLicenseUrl,
		DriverLicenseGetTime:    form.DriverLicenseGetTime,
		DriverLicenseExpireTime: form.DriverLicenseExpireTime,
		DriverAge:               form.DriverAge,
		HealthCertUrl:           form.HealthCertUrl,
		ResidencePermitUrl:      form.ResidencePermitUrl,
		CityCode:                form.CityCode,
		VehicleNo:               form.VehicleNo,
		VehicleType:             form.VehicleType,
		VehicleBrand:            form.VehicleBrand,
		VehicleModel:            form.VehicleModel,
		VehicleColor:            form.VehicleColor,
		Vin:                     form.Vin,
		EngineNo:                form.EngineNo,
		LicenseNo:               form.LicenseNo,
		LicenseExpireDate:       form.LicenseExpireDate,
		InsuranceExpireDate:     form.InsuranceExpireDate,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "信息注册成功"})
}

// UpdDriver 修改司机信息
//
// 参数:
//   - c *gin.Context: Gin上下文
//
// 处理逻辑:
//  1. 绑定请求参数
//  2. 从上下文获取用户ID
//  3. 调用司机服务修改信息
//  4. 返回修改结果
func UpdDriver(c *gin.Context) {
	var form request.UpdDriver
	// 根据 Content-Type Header 推断使用哪个绑定器。
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := global.DriverClient.UpdDriver(c, &proto.UpdDriverReq{
		UserId:                  int64(c.GetUint64("userId")),
		RealName:                form.RealName,
		IdCardNo:                form.IdCardNo,
		IdCardFrontUrl:          form.IdCardFrontUrl,
		IdCardBackUrl:           form.IdCardBackUrl,
		IdCardExpireTime:        form.IdCardExpireTime,
		DriverLicenseNo:         form.DriverLicenseNo,
		DriverLicenseUrl:        form.DriverLicenseUrl,
		DriverLicenseGetTime:    form.DriverLicenseGetTime,
		DriverLicenseExpireTime: form.DriverLicenseExpireTime,
		DriverAge:               form.DriverAge,
		HealthCertUrl:           form.HealthCertUrl,
		ResidencePermitUrl:      form.ResidencePermitUrl,
		CityCode:                form.CityCode,
		VehicleNo:               form.VehicleNo,
		VehicleType:             form.VehicleType,
		VehicleBrand:            form.VehicleBrand,
		VehicleModel:            form.VehicleModel,
		VehicleColor:            form.VehicleColor,
		Vin:                     form.Vin,
		EngineNo:                form.EngineNo,
		LicenseNo:               form.LicenseNo,
		LicenseExpireDate:       form.LicenseExpireDate,
		InsuranceExpireDate:     form.InsuranceExpireDate,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "信息修改成功"})
}
