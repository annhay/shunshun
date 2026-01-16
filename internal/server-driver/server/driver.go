package server

import (
	"context"
	"errors"
	"fmt"
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

	// 如果司机上传了身份证正面照片，使用OCR自动识别信息
	if in.IdCardFrontUrl != "" {
		// 实际使用时，需要根据前端传递的参数进行调整
		ocrResult, err := utils.AliOCR(in.IdCardFrontUrl, "id-card-front")
		if err == nil {
			// 解析 OCR 识别结果
			parsedResult, parseErr := utils.ParseOCRResult(ocrResult, "id-card-front")
			if parseErr == nil {
				// 使用 OCR 识别结果填充字段
				if parsedResult.RealName != "" && in.RealName == "" {
					in.RealName = parsedResult.RealName
				}
				if parsedResult.IdCard != "" && in.IdCardNo == "" {
					in.IdCardNo = parsedResult.IdCard
				}
			}
		}
	}

	// 如果司机上传了驾驶证照片，使用OCR自动识别信息
	if in.DriverLicenseUrl != "" {
		// 这里假设 in.DriverLicenseUrl 是驾驶证照片的URL
		// 实际使用时，需要根据前端传递的参数进行调整
		ocrResult, err := utils.AliOCR(in.DriverLicenseUrl, "id-card-front")
		if err == nil {
			// 解析 OCR 识别结果
			parsedResult, parseErr := utils.ParseOCRResult(ocrResult, "id-card-front")
			if parseErr == nil {
				// 使用 OCR 识别结果填充字段
				if parsedResult.IdCard != "" && in.DriverLicenseNo == "" {
					in.DriverLicenseNo = parsedResult.IdCard
				}
			}
		}
	}

	// 身份证验证
	if in.RealName != "" && in.IdCardNo != "" {
		isValid, err := utils.VerifyIdCard(in.RealName, in.IdCardNo)
		if err != nil {
			return nil, fmt.Errorf("身份证验证失败: %v", err)
		}
		if !isValid {
			return nil, errors.New("身份证信息不匹配")
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

	// 如果司机上传了身份证正面照片，使用OCR自动识别信息
	if in.IdCardFrontUrl != "" {
		// 这里假设 in.IdCardFrontUrl 是身份证正面照片的URL
		// 实际使用时，需要根据前端传递的参数进行调整
		ocrResult, err := utils.AliOCR(in.IdCardFrontUrl, "id-card-front")
		if err == nil {
			// 解析 OCR 识别结果
			parsedResult, parseErr := utils.ParseOCRResult(ocrResult, "id-card-front")
			if parseErr == nil {
				// 使用 OCR 识别结果填充字段
				if parsedResult.RealName != "" && in.RealName == "" {
					in.RealName = parsedResult.RealName
				}
				if parsedResult.IdCard != "" && in.IdCardNo == "" {
					in.IdCardNo = parsedResult.IdCard
				}
			}
		}
	}

	// 如果司机上传了驾驶证照片，使用OCR自动识别信息
	if in.DriverLicenseUrl != "" {
		// 这里假设 in.DriverLicenseUrl 是驾驶证照片的URL
		// 实际使用时，需要根据前端传递的参数进行调整
		ocrResult, err := utils.AliOCR(in.DriverLicenseUrl, "id-card-front")
		if err == nil {
			// 解析 OCR 识别结果
			parsedResult, parseErr := utils.ParseOCRResult(ocrResult, "id-card-front")
			if parseErr == nil {
				// 使用 OCR 识别结果填充字段
				if parsedResult.IdCard != "" && in.DriverLicenseNo == "" {
					in.DriverLicenseNo = parsedResult.IdCard
				}
			}
		}
	}

	// 只更新提供的字段
	if in.RealName != "" {
		driver.RealName = in.RealName
	}
	if in.IdCardNo != "" {
		driver.IdCardNo = in.IdCardNo
	}
	if in.IdCardFrontUrl != "" {
		driver.IdCardFrontUrl = in.IdCardFrontUrl
	}
	if in.IdCardBackUrl != "" {
		driver.IdCardBackUrl = in.IdCardBackUrl
	}
	if in.IdCardExpireTime != "" {
		driver.IdCardExpireTime = utils.StringTransformationTime(in.IdCardExpireTime)
	}
	if in.DriverLicenseNo != "" {
		driver.DriverLicenseNo = in.DriverLicenseNo
	}
	if in.DriverLicenseUrl != "" {
		driver.DriverLicenseUrl = in.DriverLicenseUrl
	}
	if in.DriverLicenseGetTime != "" {
		driver.DriverLicenseGetTime = utils.StringTransformationTime(in.DriverLicenseGetTime)
	}
	if in.DriverAge > 0 {
		driver.DrivingAge = uint8(in.DriverAge)
	}
	if in.HealthCertUrl != "" {
		driver.HealthCertUrl = in.HealthCertUrl
	}
	if in.ResidencePermitUrl != "" {
		driver.ResidencePermitUrl = in.ResidencePermitUrl
	}
	if in.CityCode != "" {
		driver.CityCode = in.CityCode
	}
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
