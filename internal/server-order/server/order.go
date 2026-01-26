package server

import (
	"context"
	"errors"
	"fmt"
	"shunshun/internal/pkg/global"
	"shunshun/internal/pkg/model"
	"shunshun/internal/pkg/utils"
	"shunshun/internal/proto"
	"shunshun/internal/server-order/task"
	"strconv"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Server 订单服务结构体
// 实现了proto.UnimplementedOrderServer接口
type Server struct {
	proto.UnimplementedOrderServer
}

// NewOrder 创建订单
func (s *Server) NewOrder(ctx context.Context, in *proto.NewOrderReq) (*proto.NewOrderResp, error) {
	global.Logger.Info("开始创建订单", zap.Int64("userId", in.UserId), zap.String("tripType", in.TripType), zap.String("rideMode", in.RideMode))

	// 使用Redis setnx防止用户重复下单
	lockKey := fmt.Sprintf("order:lock:%d:%d", in.UserId, time.Now().UnixNano()/int64(time.Millisecond))
	lockValue := "1"
	lockExpiration := 30 * time.Second

	// 尝试获取锁
	success, err := global.Rdb.SetNX(ctx, lockKey, lockValue, lockExpiration).Result()
	if err != nil {
		global.Logger.Error("获取订单锁失败", zap.Int64("userId", in.UserId), zap.Error(err))
		return nil, errors.New("系统繁忙，请稍后重试")
	}

	if !success {
		global.Logger.Warn("用户重复下单", zap.Int64("userId", in.UserId))
		return nil, errors.New("您正在创建订单，请稍后重试")
	}

	// 确保在函数结束时释放锁
	defer func() {
		err := global.Rdb.Del(ctx, lockKey).Err()
		if err != nil {
			global.Logger.Error("释放订单锁失败", zap.Int64("userId", in.UserId), zap.Error(err))
		}
	}()

	// 行程模式处理
	// 1-顺风车, 2-打车
	// 根据用户要求，先不管送货和送宠物
	if in.TripType != "1" && in.TripType != "2" {
		global.Logger.Error("行程模式不支持", zap.Int64("userId", in.UserId), zap.String("tripType", in.TripType))
		return nil, errors.New("行程模式不支持")
	}

	// 打车模式处理
	if in.TripType == "2" {
		// 打车默认1人，不可设置
		in.PassengerNum = 1
		global.Logger.Info("打车模式，乘客人数固定为1", zap.Int64("userId", in.UserId))
	}

	// 顺风车模式处理
	if in.TripType == "1" {
		// 检查乘车模式
		// 1-拼单, 2-只拼一单（拼单）, 3-独享【顺风车】
		if in.RideMode != "1" && in.RideMode != "2" && in.RideMode != "3" {
			global.Logger.Error("顺风车乘车模式不支持", zap.Int64("userId", in.UserId), zap.String("rideMode", in.RideMode))
			return nil, errors.New("顺风车乘车模式不支持")
		}

		// 检查等待时间（顺风车）
		// 等待时间(分钟) (最短10分钟,最长3小时)
		if in.WaitingTime < 10 || in.WaitingTime > 180 {
			global.Logger.Error("顺风车等待时间超出范围", zap.Int64("userId", in.UserId), zap.Int64("waitingTime", in.WaitingTime))
			return nil, errors.New("顺风车等待时间必须在10-180分钟之间")
		}

		// 拼单模式处理
		if in.RideMode == "1" || in.RideMode == "2" {
			// 拼单模式需要处理拼单相关字段
			global.Logger.Info("拼单模式", zap.Int64("userId", in.UserId), zap.String("rideMode", in.RideMode))

			// 查找是否有合适的拼单主订单
			// 查找条件：顺风车、拼单模式、待接单状态、出发时间相近、路线相似
			// 这里简化处理，实际应该根据路线相似度和时间匹配度来查找
			// 暂时先创建新的主订单

			// 对于拼单模式，设置拼单相关字段
			// 在订单创建时，我们会在下面的订单创建部分设置这些字段
		}
	}

	// 车辆类型验证
	// 1-经济车 2-商务车 3-六座车
	if in.CarType != "1" && in.CarType != "2" && in.CarType != "3" {
		global.Logger.Error("车辆类型不支持", zap.Int64("userId", in.UserId), zap.String("carType", in.CarType))
		return nil, errors.New("车辆类型不支持")
	}

	// 支付方式验证
	// 1-支付宝 2-微信 3-银联 4-余额
	if in.PaymentMethod != "1" && in.PaymentMethod != "2" && in.PaymentMethod != "3" && in.PaymentMethod != "4" {
		global.Logger.Error("支付方式不支持", zap.Int64("userId", in.UserId), zap.String("paymentMethod", in.PaymentMethod))
		return nil, errors.New("支付方式不支持")
	}

	//账号信息查询，验证登录
	var user model.ShunUser
	if err := user.GetUserById(global.DB, int(in.UserId)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.Error("用户不存在", zap.Int64("userId", in.UserId))
			return nil, err
		}
	}

	//实名认证验证
	if user.RealName == "" || user.IdCard == "" {
		global.Logger.Error("用户未实名认证", zap.Int64("userId", in.UserId))
		return nil, errors.New("请前往实名认证")
	}

	global.Logger.Info("用户验证通过", zap.Int64("userId", in.UserId), zap.String("realName", user.RealName))

	// 获取起点经纬度
	startGeo, err := utils.Geocode(in.StratDetailAddress)
	if err != nil {
		global.Logger.Error("获取起点位置信息失败", zap.Int64("userId", in.UserId), zap.String("address", in.StratDetailAddress), zap.Error(err))
		return nil, errors.New("获取起点位置信息失败")
	}
	startLng, _ := strconv.ParseFloat(startGeo.Lng, 64)
	startLat, _ := strconv.ParseFloat(startGeo.Lat, 64)

	// 获取终点经纬度
	endGeo, err := utils.Geocode(in.EndDetailAddress)
	if err != nil {
		global.Logger.Error("获取终点位置信息失败", zap.Int64("userId", in.UserId), zap.String("address", in.EndDetailAddress), zap.Error(err))
		return nil, errors.New("获取终点位置信息失败")
	}
	endLng, _ := strconv.ParseFloat(endGeo.Lng, 64)
	endLat, _ := strconv.ParseFloat(endGeo.Lat, 64)

	global.Logger.Info("地理位置获取成功", zap.Int64("userId", in.UserId), zap.String("startAddress", in.StratDetailAddress), zap.String("endAddress", in.EndDetailAddress))

	// 计算行程距离

	travelDistance, err := task.CalculateTravelDistance(startLng, startLat, endLng, endLat)
	if err != nil {
		global.Logger.Error("计算行程距离失败", zap.Int64("userId", in.UserId), zap.Error(err))
		return nil, errors.New("计算行程距离失败")
	}

	global.Logger.Info("行程距离计算成功", zap.Int64("userId", in.UserId), zap.Float64("distance", travelDistance))

	// 计算预估金额
	var estimatedAmount float64
	if in.TripType == "2" {
		// 打车订单计算预估金额
		estimatedAmount = task.CalculateEstimatedAmount(travelDistance, in.CarType)
		global.Logger.Info("打车订单预估金额计算成功", zap.Int64("userId", in.UserId), zap.Float64("estimatedAmount", estimatedAmount))
	} else {
		// 顺风车订单计算预估金额
		// 基础预估金额
		baseAmount := task.CalculateEstimatedAmount(travelDistance, in.CarType)

		// 拼单模式折扣
		if in.RideMode == "1" || in.RideMode == "2" {
			// 拼单模式，价格打8折
			estimatedAmount = baseAmount * 0.8
			global.Logger.Info("顺风车拼单模式预估金额计算成功", zap.Int64("userId", in.UserId), zap.Float64("baseAmount", baseAmount), zap.Float64("estimatedAmount", estimatedAmount))
		} else {
			// 独享模式，原价
			estimatedAmount = baseAmount
			global.Logger.Info("顺风车独享模式预估金额计算成功", zap.Int64("userId", in.UserId), zap.Float64("estimatedAmount", estimatedAmount))
		}
	}

	// 处理优惠券
	couponAmount, err := task.CalculateCouponAmount(in.UserId, in.CouponId, estimatedAmount)
	if err != nil {
		global.Logger.Error("计算优惠券金额失败", zap.Int64("userId", in.UserId), zap.Int64("couponId", in.CouponId), zap.Error(err))
		return nil, err
	}

	if in.CouponId > 0 {
		global.Logger.Info("优惠券验证通过", zap.Int64("userId", in.UserId), zap.Int64("couponId", in.CouponId), zap.Float64("couponAmount", couponAmount))
	}

	// 开始事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			global.Logger.Error("创建订单发生 panic", zap.Int64("userId", in.UserId), zap.Any("recover", r))
			tx.Rollback()
		}
	}()

	// 创建订单信息
	order := &model.ShunOrder{
		OrderCode:          utils.OrderCodeRandom(in.UserId),
		UserId:             in.UserId,
		TripType:           in.TripType,
		RideMode:           in.RideMode,
		PassengerNum:       int32(in.PassengerNum),
		CarType:            in.CarType,
		DepartureTime:      utils.StringTransformationTime(in.DepartureTime),
		WaitingTime:        int32(in.WaitingTime),
		Remark:             in.Remark,
		StartDetailAddress: in.StratDetailAddress, //起点详细地址
		StartLongitude:     startLng,
		StartLatitude:      startLat,
		EndDetailAddress:   in.EndDetailAddress, //终点详细地址
		EndLongitude:       endLng,
		EndLatitude:        endLat,
		TravelDistance:     travelDistance,
		EstimatedAmount:    estimatedAmount, //预估金额
		CouponId:           in.CouponId,
		CouponAmount:       couponAmount,
		PaymentMethod:      in.PaymentMethod,
		OrderStatus:        "1", // 待接单
	}

	// 处理拼单相关字段
	if in.TripType == "1" && (in.RideMode == "1" || in.RideMode == "2") {
		// 查找是否有合适的拼单主订单
		// 这里简化处理，实际应该根据路线相似度和时间匹配度来查找
		// 查找条件：顺风车、拼单模式、待接单状态、出发时间相近、路线相似
		// 暂时先创建新的主订单，将当前订单设为主订单
		order.MainGroupOrderId = 0 // 0表示主订单
		order.GroupUserCount = 1   // 初始拼单人数为1
		global.Logger.Info("创建拼单主订单", zap.Int64("userId", in.UserId), zap.String("orderCode", order.OrderCode))
	}

	// 设置支付状态
	if in.TripType == "1" {
		// 顺风车订单：待支付（需要10分钟内支付）
		order.PaymentStatus = "1"
		global.Logger.Info("顺风车订单设置为待支付状态", zap.Int64("userId", in.UserId))
	} else {
		// 打车订单：无需立即支付
		order.PaymentStatus = "3" // 无需支付（后续根据实际金额支付）
		global.Logger.Info("打车订单设置为无需立即支付状态", zap.Int64("userId", in.UserId))
	}

	// 保存订单
	if err := tx.Create(order).Error; err != nil {
		global.Logger.Error("创建订单失败", zap.Int64("userId", in.UserId), zap.Error(err))
		tx.Rollback()
		return nil, errors.New("创建订单失败")
	}

	// 更新优惠券状态
	if in.CouponId > 0 {
		if err := task.UpdateCouponStatus(in.UserId, in.CouponId); err != nil {
			global.Logger.Error("更新优惠券状态失败", zap.Int64("userId", in.UserId), zap.Int64("couponId", in.CouponId), zap.Error(err))
			tx.Rollback()
			return nil, errors.New("更新优惠券状态失败")
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		global.Logger.Error("提交事务失败", zap.Int64("userId", in.UserId), zap.Error(err))
		return nil, errors.New("提交事务失败")
	}

	global.Logger.Info("订单创建成功", zap.Int64("userId", in.UserId), zap.String("orderCode", order.OrderCode), zap.Float64("estimatedAmount", estimatedAmount), zap.Float64("couponAmount", couponAmount))

	// 发送订单创建成功通知
	notification := utils.OrderNotification{
		OrderID:     int64(order.Id),
		UserID:      in.UserId,
		OrderStatus: order.OrderStatus,
		Message:     "订单创建成功，等待司机接单",
	}

	nErr := utils.SendOrderNotification(notification)
	if nErr != nil {
		global.Logger.Error("Failed to send order creation notification", zap.Error(nErr))
	}

	var payUrl string
	if in.TripType == "1" {
		// 顺风车订单生成支付链接（需要10分钟内支付）
		// 使用计算的预估金额
		payAmount := estimatedAmount - couponAmount
		if payAmount < 0 {
			payAmount = 0
		}
		payUrl = utils.AliPay(order.OrderCode, strconv.FormatFloat(payAmount, 'f', 2, 64))
		global.Logger.Info("顺风车订单支付链接生成成功", zap.Int64("userId", in.UserId), zap.String("orderCode", order.OrderCode), zap.Float64("payAmount", payAmount))
	} else {
		// 打车订单不需要生成支付链接
		payUrl = ""
		global.Logger.Info("打车订单不需要支付链接", zap.Int64("userId", in.UserId), zap.String("orderCode", order.OrderCode))
	}

	return &proto.NewOrderResp{
		OrderCode: order.OrderCode,
		PayUrl:    payUrl,
	}, nil
}

//支付异步回调
//支付同步回调
//取消订单
//用户订单列表
//用户订单详情
//司机匹配的订单列表

// AcceptOrders 司机接单
func (s *Server) AcceptOrders(ctx context.Context, in *proto.AcceptOrdersReq) (*proto.AcceptOrdersResp, error) {
	global.Logger.Info("司机开始接单", zap.Int64("userId", in.UserId), zap.Int64("orderId", in.OrderId))
	// 查询用户信息
	var user model.ShunUser
	if err := user.GetUserById(global.DB, int(in.UserId)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.Error("用户不存在", zap.Int64("userId", in.UserId))
			return nil, err
		}
	}
	// 实名认证验证
	if user.RealName == "" || user.IdCard == "" {
		global.Logger.Error("司机未实名认证", zap.Int64("userId", in.UserId))
		return nil, errors.New("账号未实名认证")
	}

	// 查询司机信息
	var driver model.ShunDriver
	if err := driver.GetDriverByUserId(global.DB, in.UserId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.Error("司机未注册", zap.Int64("userId", in.UserId))
			return nil, err
		}
	}

	// 查询车辆信息
	var car model.ShunCar
	if err := car.GetCarByDriverIdOnNormal(global.DB, driver.Id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.Error("司机无可用车辆", zap.Int64("userId", in.UserId), zap.Int("driverId", int(driver.Id)))
			return nil, err
		}
	}

	// 查询订单信息
	var order model.ShunOrder
	if err := global.DB.First(&order, in.OrderId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.Error("订单不存在", zap.Int64("userId", in.UserId), zap.Int64("orderId", in.OrderId))
			return nil, err
		}
	}

	// 检查订单状态
	if order.OrderStatus != "1" {
		global.Logger.Error("订单状态不正确", zap.Int64("userId", in.UserId), zap.Int64("orderId", in.OrderId), zap.String("orderStatus", order.OrderStatus))
		return nil, errors.New("订单状态不正确，无法接单")
	}

	// 计算司机当前已接但未完成的订单总载客人数
	// 查询司机已接单但未完成的订单
	var activeOrders []model.ShunOrder
	if err := global.DB.Where("driver_id = ? AND order_status IN ('2', '3', '4', '5', '6')", driver.Id).Find(&activeOrders).Error; err != nil {
		global.Logger.Error("查询司机活跃订单失败", zap.Int64("userId", in.UserId), zap.Int("driverId", int(driver.Id)), zap.Error(err))
		return nil, errors.New("查询司机活跃订单失败")
	}

	// 计算当前订单的载客人数
	currentOrderPassengerNum := int(order.PassengerNum)
	if order.TripType == "2" && in.ActualPassengerNum > 0 {
		// 打车订单使用实际乘车人数
		currentOrderPassengerNum = int(in.ActualPassengerNum)
	}

	// 计算总载客人数
	totalPassengerNum := currentOrderPassengerNum
	for _, activeOrder := range activeOrders {
		// 只计算非拼单订单的载客人数，因为拼单订单的乘客是一起出行的
		if activeOrder.TripType == "1" && (activeOrder.RideMode == "1" || activeOrder.RideMode == "2") {
			// 拼单订单，不计算入总载客人数
			continue
		}
		// 打车订单使用实际乘车人数
		if activeOrder.TripType == "2" && activeOrder.ActualPassengerNum > 0 {
			totalPassengerNum += int(activeOrder.ActualPassengerNum)
		} else {
			totalPassengerNum += int(activeOrder.PassengerNum)
		}
	}

	// 根据车型确定载客量限制
	var maxPassengerNum int
	switch car.VehicleType {
	case "1": // 经济车
		maxPassengerNum = 4
	case "2": // 商务车
		maxPassengerNum = 5
	case "3": // 六座车
		maxPassengerNum = 6
	default:
		maxPassengerNum = 4 // 默认经济车载客量
	}

	// 检查是否超载
	if totalPassengerNum > maxPassengerNum {
		global.Logger.Error("车辆超载，无法接单",
			zap.Int64("userId", in.UserId),
			zap.Int("driverId", int(driver.Id)),
			zap.String("vehicleType", car.VehicleType),
			zap.Int("maxPassengerNum", maxPassengerNum),
			zap.Int("totalPassengerNum", totalPassengerNum),
			zap.Int64("orderId", in.OrderId),
		)
		return nil, errors.New("车辆超载，无法接单")
	}

	// 检查是否一次性多接单
	// 只有拼单订单（顺风车+拼单模式）才能一次性多接单
	// 非拼单订单，司机当前不能有其他未完成的订单
	if !(order.TripType == "1" && (order.RideMode == "1" || order.RideMode == "2")) {
		// 非拼单订单，检查是否已有其他未完成的订单
		if len(activeOrders) > 0 {
			global.Logger.Error("非拼单订单，无法同时接多个订单",
				zap.Int64("userId", in.UserId),
				zap.Int("driverId", int(driver.Id)),
				zap.Int64("orderId", in.OrderId),
				zap.Int("activeOrderCount", len(activeOrders)),
			)
			return nil, errors.New("非拼单订单，无法同时接多个订单")
		}
	}

	global.Logger.Info("车辆载客量检查通过",
		zap.Int64("userId", in.UserId),
		zap.Int("driverId", int(driver.Id)),
		zap.String("vehicleType", car.VehicleType),
		zap.Int("maxPassengerNum", maxPassengerNum),
		zap.Int("totalPassengerNum", totalPassengerNum),
		zap.Int64("orderId", in.OrderId),
		zap.Int("activeOrderCount", len(activeOrders)),
	)

	// 行程模式处理
	// 1-顺风车, 2-打车
	if order.TripType != "1" && order.TripType != "2" {
		global.Logger.Error("行程模式不支持", zap.Int64("userId", in.UserId), zap.String("tripType", order.TripType), zap.Int64("orderId", in.OrderId))
		return nil, errors.New("行程模式不支持")
	}

	// 根据行程模式进行不同的处理
	if order.TripType == "1" {
		// 顺风车模式
		global.Logger.Info("顺风车接单", zap.Int64("userId", in.UserId), zap.Int64("orderId", in.OrderId), zap.String("rideMode", order.RideMode))
		// 顺风车需要处理拼单、等待时间等特殊逻辑
	} else if order.TripType == "2" {
		// 打车模式
		global.Logger.Info("打车接单", zap.Int64("userId", in.UserId), zap.Int64("orderId", in.OrderId), zap.Int64("actualPassengerNum", in.ActualPassengerNum))
		// 打车需要处理实际乘客人数等特殊逻辑

		// 验证实际乘车人数
		if in.ActualPassengerNum <= 0 {
			global.Logger.Error("实际乘车人数无效", zap.Int64("userId", in.UserId), zap.Int64("orderId", in.OrderId), zap.Int64("actualPassengerNum", in.ActualPassengerNum))
			return nil, errors.New("实际乘车人数必须大于0")
		}
	}

	// 计算行程距离
	var travelDistance float64 = 0
	if order.StartLongitude != 0 && order.StartLatitude != 0 && order.EndLongitude != 0 && order.EndLatitude != 0 {
		distance, err := task.CalculateTravelDistance(order.StartLongitude, order.StartLatitude, order.EndLongitude, order.EndLatitude)
		if err != nil {
			global.Logger.Warn("计算行程距离失败，使用默认值", zap.Int64("userId", in.UserId), zap.Int64("orderId", in.OrderId), zap.Error(err))
		} else {
			travelDistance = distance
			global.Logger.Info("行程距离计算成功", zap.Int64("userId", in.UserId), zap.Int64("orderId", in.OrderId), zap.Float64("distance", travelDistance))
		}
	}

	// 开始事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			global.Logger.Error("接单过程发生 panic", zap.Int64("userId", in.UserId), zap.Int64("orderId", in.OrderId), zap.Any("recover", r))
			tx.Rollback()
		}
	}()

	// 更新订单信息
	orderUpdates := map[string]interface{}{
		"DriverId":       int64(driver.Id),
		"CarId":          int64(car.Id),
		"TravelDistance": travelDistance,
		"OrderStatus":    "2", // 已接单
	}

	// 对于打车订单，添加实际乘车人数和重新计算金额
	if order.TripType == "2" {
		// 添加实际乘车人数
		orderUpdates["ActualPassengerNum"] = in.ActualPassengerNum

		// 根据实际乘车人数重新计算金额
		// 这里需要根据实际的计价规则来计算，暂时使用简单的计算方式
		// 基础金额 + 人数加价
		baseAmount := order.EstimatedAmount
		personAddAmount := float64(in.ActualPassengerNum-1) * 0.5 // 每多一人加价0.5元
		actualAmount := baseAmount + personAddAmount

		orderUpdates["ActualAmount"] = actualAmount

		global.Logger.Info("打车订单金额重新计算", zap.Int64("userId", in.UserId), zap.Int64("orderId", in.OrderId), zap.Float64("baseAmount", baseAmount), zap.Float64("personAddAmount", personAddAmount), zap.Float64("actualAmount", actualAmount))
	}

	if err := tx.Model(&model.ShunOrder{}).Where("id = ?", in.OrderId).Updates(orderUpdates).Error; err != nil {
		global.Logger.Error("更新订单信息失败", zap.Int64("userId", in.UserId), zap.Int64("orderId", in.OrderId), zap.Error(err))
		tx.Rollback()
		return nil, errors.New("更新订单信息失败")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		global.Logger.Error("提交事务失败", zap.Int64("userId", in.UserId), zap.Int64("orderId", in.OrderId), zap.Error(err))
		return nil, errors.New("提交事务失败")
	}

	global.Logger.Info("司机接单成功", zap.Int64("userId", in.UserId), zap.Int64("orderId", in.OrderId), zap.Int("driverId", int(driver.Id)), zap.Int("carId", int(car.Id)))

	// 发送订单接单成功通知给用户
	userNotification := utils.OrderNotification{
		OrderID:     in.OrderId,
		UserID:      order.UserId,
		DriverID:    int64(driver.Id),
		OrderStatus: "2", // 已接单
		Message:     "订单已被司机接单，司机正在赶来",
	}

	nErr := utils.SendOrderNotification(userNotification)
	if nErr != nil {
		global.Logger.Error("Failed to send order acceptance notification to user", zap.Error(nErr))
	}

	// 发送订单接单成功通知给司机
	driverNotification := utils.OrderNotification{
		OrderID:     in.OrderId,
		UserID:      order.UserId,
		DriverID:    int64(driver.Id),
		OrderStatus: "2", // 已接单
		Message:     "接单成功，请及时联系乘客",
	}

	nErr = utils.SendOrderNotification(driverNotification)
	if nErr != nil {
		global.Logger.Error("Failed to send order acceptance notification to driver", zap.Error(nErr))
	}

	return &proto.AcceptOrdersResp{}, nil
}

//司机接单列表
//司机接单详情
