package handler

import (
	"net/http"
	"shunshun/internal/api-gateway/request"
	"shunshun/internal/pkg/global"
	"shunshun/internal/proto"

	"github.com/gin-gonic/gin"
)

// NewDriver 司机认证
func NewDriver(c *gin.Context) {
	var form request.NewDriver
	// 根据 Content-Type Header 推断使用哪个绑定器。
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := global.DriverClient.NewDriver(c, &proto.NewDriverReq{
		UserId:                  int64(c.GetUint("userId")),
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
func UpdDriver(c *gin.Context) {
	var form request.UpdDriver
	// 根据 Content-Type Header 推断使用哪个绑定器。
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := global.DriverClient.UpdDriver(c, &proto.UpdDriverReq{
		UserId:                  int64(c.GetUint("userId")),
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
