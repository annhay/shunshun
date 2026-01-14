package server

import (
	"context"
	"errors"
	"shunshun/internal/pkg/global"
	"shunshun/internal/pkg/model"
	"shunshun/internal/pkg/utils"
	"shunshun/internal/proto"
	"time"

	"gorm.io/gorm"
)

type Server struct {
	proto.UnimplementedDriverServer
}

// NewDriver 司机认证（司机信息添加）
func (s *Server) NewDriver(_ context.Context, in *proto.NewDriverReq) (*proto.NewDriverResp, error) {
	var user model.ShunUser
	if err := user.GetUserById(global.DB, int(in.UserId)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	var driver model.ShunDriver
	if err := driver.GetDriverByUserId(global.DB, in.UserId); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("您已认证")
		}
	}
	newDriver := &model.ShunDriver{
		UserId:                  uint64(in.UserId),
		DriverNo:                utils.DriverNoRandom(in.UserId, in.CityCode),
		RealName:                in.RealName,
		IdCardNo:                in.IdCardNo,
		IdCardFrontUrl:          in.IdCardFrontUrl,
		IdCardBackUrl:           in.IdCardBackUrl,
		IdCardExpireTime:        utils.StringTransformationTime(in.IdCardExpireTime),
		DriverLicenseNo:         in.DriverLicenseNo,
		DriverLicenseUrl:        in.DriverLicenseUrl,
		DriverLicenseGetTime:    utils.StringTransformationTime(in.DriverLicenseGetTime),
		DriverLicenseExpireTime: utils.StringTransformationTime(in.DriverLicenseExpireTime),
		DrivingAge:              uint8(in.DriverAge),
		HealthCertUrl:           in.HealthCertUrl,
		ResidencePermitUrl:      in.ResidencePermitUrl,
		CityCode:                in.CityCode,
	}
	if err := newDriver.CreateDriver(global.DB); err != nil {
		return nil, err
	}
	car := &model.ShunCar{
		DriverId:            newDriver.Id,
		VehicleNo:           in.VehicleNo,
		VehicleType:         in.VehicleType,
		VehicleBrand:        in.VehicleBrand,
		VehicleModel:        in.VehicleModel,
		VehicleColor:        in.VehicleColor,
		Vin:                 in.Vin,
		EngineNo:            in.EngineNo,
		RegisterDate:        time.Now(),
		LicenseNo:           in.LicenseNo,
		LicenseExpireDate:   utils.StringTransformationTime(in.LicenseExpireDate),
		InsuranceExpireDate: utils.StringTransformationTime(in.InsuranceExpireDate),
	}
	if err := car.CreateCar(global.DB); err != nil {
		return nil, err
	}
	return &proto.NewDriverResp{}, nil
}

// UpdDriver 司机信息修改
func (s *Server) UpdDriver(_ context.Context, in *proto.UpdDriverReq) (*proto.UpdDriverResp, error) {
	var user model.ShunUser
	if err := user.GetUserById(global.DB, int(in.UserId)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	var driver model.ShunDriver
	if err := driver.GetDriverByUserId(global.DB, in.UserId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("未注册司机")
		}
	}

	// 只更新提供的字段
	driver.RealName = in.RealName
	driver.IdCardNo = in.IdCardNo
	driver.IdCardFrontUrl = in.IdCardFrontUrl
	driver.IdCardBackUrl = in.IdCardBackUrl
	driver.IdCardExpireTime = utils.StringTransformationTime(in.IdCardExpireTime)
	driver.DriverLicenseNo = in.DriverLicenseNo
	driver.DriverLicenseUrl = in.DriverLicenseUrl
	driver.DriverLicenseGetTime = utils.StringTransformationTime(in.DriverLicenseGetTime)
	driver.DrivingAge = uint8(in.DriverAge)
	driver.HealthCertUrl = in.HealthCertUrl
	driver.ResidencePermitUrl = in.ResidencePermitUrl
	driver.CityCode = in.CityCode
	if err := driver.Editor(global.DB); err != nil {
		return nil, err
	}
	// 处理车辆信息：如果提供了车辆信息，就停用之前的车辆并创建新车辆
	if in.VehicleNo != "" {
		// 获取司机之前的车辆
		var car model.ShunCar
		cars, err := car.GetCarsByDriverId(global.DB, driver.Id)
		if err != nil {
			return nil, err
		}

		// 停用之前的车辆
		for _, oldCar := range cars {
			oldCar.Status = "3" // 3-停用
			oldCar.UpdatedAt = time.Now()
			if err := oldCar.UpdateCar(global.DB); err != nil {
				return nil, err
			}
		}

		// 创建新的车辆记录
		newCar := &model.ShunCar{
			DriverId:            driver.Id,
			VehicleNo:           in.VehicleNo,
			VehicleType:         in.VehicleType,
			VehicleBrand:        in.VehicleBrand,
			VehicleModel:        in.VehicleModel,
			VehicleColor:        in.VehicleColor,
			Vin:                 in.Vin,
			EngineNo:            in.EngineNo,
			RegisterDate:        time.Now(),
			LicenseNo:           in.LicenseNo,
			LicenseExpireDate:   utils.StringTransformationTime(in.LicenseExpireDate),
			InsuranceExpireDate: utils.StringTransformationTime(in.InsuranceExpireDate),
		}
		if err := newCar.CreateCar(global.DB); err != nil {
			return nil, err
		}
	}

	return &proto.UpdDriverResp{}, nil
}
