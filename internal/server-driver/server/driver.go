package server

import (
	"context"
	"encoding/json"
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
				// 验证 OCR 识别结果与用户填写信息是否一致
				if parsedResult.RealName != "" && in.RealName != "" && parsedResult.RealName != in.RealName {
					return nil, errors.New("身份证姓名与上传照片信息不一致")
				}
				if parsedResult.IdCard != "" && in.IdCardNo != "" && parsedResult.IdCard != in.IdCardNo {
					return nil, errors.New("身份证号码与上传照片信息不一致")
				}
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
				// 验证 OCR 识别结果与用户填写信息是否一致
				if parsedResult.IdCard != "" && in.DriverLicenseNo != "" && parsedResult.IdCard != in.DriverLicenseNo {
					return nil, errors.New("驾驶证号码与上传照片信息不一致")
				}
				// 使用 OCR 识别结果填充字段
				if parsedResult.IdCard != "" && in.DriverLicenseNo == "" {
					in.DriverLicenseNo = parsedResult.IdCard
				}
			}
		}
	}

	// 第二次实名认证（司机认证专用）
	if in.RealName != "" && in.IdCardNo != "" {
		isValid, err := utils.VerifyIdCard(in.RealName, in.IdCardNo)
		if err != nil {
			return nil, fmt.Errorf("身份证验证失败: %v", err)
		}
		if !isValid {
			return nil, errors.New("身份证信息不匹配")
		}
	}

	// 验证驾驶证信息
	if in.DriverLicenseNo != "" {
		// 这里可以添加驾驶证验证逻辑，例如调用第三方驾驶证验证服务
		// 确保驾驶证信息与身份证信息匹配，防止盗用他人驾驶证
		// 注：实际业务中需要集成驾驶证验证API
		// 暂时添加简单的格式验证
		if len(in.DriverLicenseNo) < 10 {
			return nil, errors.New("驾驶证号码格式不正确")
		}
	}

	// 司机信息必须与第二次实名认证信息保持一致
	// 确保司机提交的所有身份相关信息都与实名认证信息一致
	if in.RealName == "" || in.IdCardNo == "" {
		return nil, errors.New("司机认证必须提供真实姓名和身份证号")
	}
	if in.DriverLicenseNo == "" || in.DriverLicenseUrl == "" {
		return nil, errors.New("司机认证必须提供驾驶证号码和照片")
	}

	// 注：实际业务中可能需要根据法律法规要求，确保司机身份信息的真实性和合法性

	newDriver := &model.ShunDriver{
		UserId:                  uint64(in.UserId),
		DriverNo:                utils.DriverNoRandom(in.UserId, in.CityCode),
		RealName:                utils.EnPwdCode([]byte(in.RealName)),
		IdCardNo:                utils.EnPwdCode([]byte(in.IdCardNo)),
		IdCardFrontUrl:          utils.EnPwdCode([]byte(in.IdCardFrontUrl)),
		IdCardBackUrl:           utils.EnPwdCode([]byte(in.IdCardBackUrl)),
		IdCardExpireTime:        utils.StringTransformationTime(in.IdCardExpireTime),
		DriverLicenseNo:         utils.EnPwdCode([]byte(in.DriverLicenseNo)),
		DriverLicenseUrl:        utils.EnPwdCode([]byte(in.DriverLicenseUrl)),
		DriverLicenseGetTime:    utils.StringTransformationTime(in.DriverLicenseGetTime),
		DriverLicenseExpireTime: utils.StringTransformationTime(in.DriverLicenseExpireTime),
		DrivingAge:              uint8(in.DriverAge),
		HealthCertUrl:           utils.EnPwdCode([]byte(in.HealthCertUrl)),
		ResidencePermitUrl:      utils.EnPwdCode([]byte(in.ResidencePermitUrl)),
		CityCode:                in.CityCode,
	}
	//司机信息认证必须要有车辆信息
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

// DetailDriver 司机详情
//
// 参数:
//   - ctx context.Context: 上下文
//   - in *proto.DetailDriverReq: 司机详情请求，包含用户ID
//
// 返回值:
//   - *proto.DetailDriverResp: 司机详情响应
//   - error: 错误信息
func (s *Server) DetailDriver(ctx context.Context, in *proto.DetailDriverReq) (*proto.DetailDriverResp, error) {
	// 生成缓存键
	cacheKey := fmt.Sprintf("driver:detail:%d", in.UserId)

	// 尝试从缓存获取
	cachedData, err := global.Rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		// 缓存命中，解析数据
		var resp proto.DetailDriverResp
		if json.Unmarshal([]byte(cachedData), &resp) == nil {
			return &resp, nil
		}
	}

	// 缓存未命中，从数据库获取
	//根据用户 ID 传输并验证 jwt
	var user model.ShunUser
	if err := user.GetUserById(global.DB, int(in.UserId)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	//根据用户 ID 验证司机注册信息
	var driver model.ShunDriver
	if err := driver.GetDriverByUserId(global.DB, in.UserId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("未认证司机")
		}
	}
	//根据查出的司机 ID 验证车辆信息
	var car model.ShunCar
	if err := car.GetCarByDriverId(global.DB, driver.Id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	//构建响应
	resp := &proto.DetailDriverResp{
		DriverLicenseNo:         string(utils.DePwdCode(driver.DriverLicenseNo)),
		DriverLicenseUrl:        string(utils.DePwdCode(driver.DriverLicenseUrl)),
		DriverLicenseGetTime:    utils.TimeTransformationString(driver.DriverLicenseGetTime),
		DriverLicenseExpireTime: utils.TimeTransformationString(driver.DriverLicenseExpireTime),
		DriverAge:               int64(driver.DrivingAge),
		HealthCertUrl:           string(utils.DePwdCode(driver.HealthCertUrl)),
		ResidencePermitUrl:      string(utils.DePwdCode(driver.ResidencePermitUrl)),
		CityCode:                driver.CityCode,
		VehicleNo:               car.VehicleNo,
		VehicleType:             car.VehicleType,
		VehicleBrand:            car.VehicleBrand,
		VehicleModel:            car.VehicleModel,
		VehicleColor:            car.VehicleColor,
		LicenseNo:               car.LicenseNo,
		LicenseExpireDate:       utils.TimeTransformationString(car.LicenseExpireDate),
		InsuranceExpireDate:     utils.TimeTransformationString(car.InsuranceExpireDate),
	}

	// 存入缓存，设置过期时间
	if data, err := json.Marshal(resp); err == nil {
		global.Rdb.Set(ctx, cacheKey, data, time.Minute*30)
	}

	return resp, nil
}

// UpdDriver 司机信息修改
//
// 参数:
//   - ctx context.Context: 上下文
//   - in *proto.UpdDriverReq: 司机信息修改请求，包含用户ID和司机信息
//
// 返回值:
//   - *proto.UpdDriverResp: 司机信息修改响应
//   - error: 错误信息
func (s *Server) UpdDriver(ctx context.Context, in *proto.UpdDriverReq) (*proto.UpdDriverResp, error) {
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
				// 验证 OCR 识别结果与用户填写信息是否一致
				if parsedResult.RealName != "" && in.RealName != "" && parsedResult.RealName != in.RealName {
					return nil, errors.New("身份证姓名与上传照片信息不一致")
				}
				if parsedResult.IdCard != "" && in.IdCardNo != "" && parsedResult.IdCard != in.IdCardNo {
					return nil, errors.New("身份证号码与上传照片信息不一致")
				}
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
				// 验证 OCR 识别结果与用户填写信息是否一致
				if parsedResult.IdCard != "" && in.DriverLicenseNo != "" && parsedResult.IdCard != in.DriverLicenseNo {
					return nil, errors.New("驾驶证号码与上传照片信息不一致")
				}
				// 使用 OCR 识别结果填充字段
				if parsedResult.IdCard != "" && in.DriverLicenseNo == "" {
					in.DriverLicenseNo = parsedResult.IdCard
				}
			}
		}
	}

	// 第二次实名认证（司机信息修改时）
	if in.RealName != "" || in.IdCardNo != "" {
		// 如果修改身份信息，必须同时提供姓名和身份证号
		if in.RealName == "" || in.IdCardNo == "" {
			return nil, errors.New("修改身份信息必须同时提供真实姓名和身份证号")
		}
		// 验证身份信息
		isValid, err := utils.VerifyIdCard(in.RealName, in.IdCardNo)
		if err != nil {
			return nil, fmt.Errorf("身份证验证失败: %v", err)
		}
		if !isValid {
			return nil, errors.New("身份证信息不匹配")
		}
	}

	// 验证驾驶证信息修改
	if in.DriverLicenseNo != "" {
		// 这里可以添加驾驶证验证逻辑，例如调用第三方驾驶证验证服务
		// 确保驾驶证信息与身份证信息匹配，防止盗用他人驾驶证
		// 注：实际业务中需要集成驾驶证验证API
		// 暂时添加简单的格式验证
		if len(in.DriverLicenseNo) < 10 {
			return nil, errors.New("驾驶证号码格式不正确")
		}
		// 如果修改驾驶证信息，必须同时提供驾驶证照片
		if in.DriverLicenseUrl == "" {
			return nil, errors.New("修改驾驶证信息必须同时提供驾驶证照片")
		}
	}

	// 只更新提供的字段
	if in.RealName != "" {
		driver.RealName = utils.EnPwdCode([]byte(in.RealName))
	}
	if in.IdCardNo != "" {
		driver.IdCardNo = utils.EnPwdCode([]byte(in.IdCardNo))
	}
	if in.IdCardFrontUrl != "" {
		driver.IdCardFrontUrl = utils.EnPwdCode([]byte(in.IdCardFrontUrl))
	}
	if in.IdCardBackUrl != "" {
		driver.IdCardBackUrl = utils.EnPwdCode([]byte(in.IdCardBackUrl))
	}
	if in.IdCardExpireTime != "" {
		driver.IdCardExpireTime = utils.StringTransformationTime(in.IdCardExpireTime)
	}
	if in.DriverLicenseNo != "" {
		driver.DriverLicenseNo = utils.EnPwdCode([]byte(in.DriverLicenseNo))
	}
	if in.DriverLicenseUrl != "" {
		driver.DriverLicenseUrl = utils.EnPwdCode([]byte(in.DriverLicenseUrl))
	}
	if in.DriverLicenseGetTime != "" {
		driver.DriverLicenseGetTime = utils.StringTransformationTime(in.DriverLicenseGetTime)
	}
	if in.DriverAge > 0 {
		driver.DrivingAge = uint8(in.DriverAge)
	}
	if in.HealthCertUrl != "" {
		driver.HealthCertUrl = utils.EnPwdCode([]byte(in.HealthCertUrl))
	}
	if in.ResidencePermitUrl != "" {
		driver.ResidencePermitUrl = utils.EnPwdCode([]byte(in.ResidencePermitUrl))
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
	// 当司机修改信息时不注册新车辆，确保不会创建新车辆记录

	// 删除缓存，保证数据一致性
	cacheKey := fmt.Sprintf("driver:detail:%d", in.UserId)
	global.Rdb.Del(ctx, cacheKey)

	return &proto.UpdDriverResp{}, nil
}
