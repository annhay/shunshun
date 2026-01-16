package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"shunshun/internal/pkg/global"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/green"
)

// OCRResult 结构化的OCR识别结果
type OCRResult struct {
	RealName   string `json:"real_name"`   // 真实姓名
	IdCard     string `json:"id_card"`     // 身份证号
	Birthday   string `json:"birthday"`    // 出生日期
	Gender     string `json:"gender"`      // 性别
	Address    string `json:"address"`     // 地址
	SchoolName string `json:"school_name"` // 学校名称
	StudentId  string `json:"student_id"`  // 学号
}

// AliOCR 阿里图片信息自动识别补充字段
// imageURL: 待检测图片链接地址
// cardType: 身份证类型，"id-card-front" 表示正面，"id-card-back" 表示反面
// 返回识别结果和错误信息
func AliOCR(imageURL, cardType string) (string, error) {
	/**
	 * 注意：此处实例化的client尽可能重复使用，提升检测性能。避免重复建立连接。
	 * 常见获取环境变量方式：
	 *     获取RAM用户AccessKey ID：os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")
	 *     获取RAM用户AccessKey Secret：os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")
	 */
	// accessorisesKey
	accessKeyID := os.Getenv(global.AppConf.AliYun.AccessKeyID)
	accessKeySecret := os.Getenv(global.AppConf.AliYun.AccessKeySecret)

	client, err := green.NewClientWithAccessKey(
		"cn-shanghai",
		accessKeyID,
		accessKeySecret)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	task1 := map[string]interface{}{"dataId": fmt.Sprintf("%d", GetTimestamp()), "url": imageURL}
	cardExtras := map[string]interface{}{"card": cardType}
	// scenes：检测场景。
	content, _ := json.Marshal(
		map[string]interface{}{
			"tasks": []map[string]interface{}{task1}, "scenes": []string{"ocr"}, "bizType": "shunshun-ocr", "extras": cardExtras,
		},
	)

	request := green.CreateImageSyncScanRequest()
	request.SetContent(content)
	response, _err := client.ImageSyncScan(request)
	if _err != nil {
		fmt.Println(_err.Error())
		return "", _err
	}
	if response.GetHttpStatus() != 200 {
		statusMsg := "response not success. status:" + strconv.Itoa(response.GetHttpStatus())
		fmt.Println(statusMsg)
		return "", fmt.Errorf(statusMsg)
	}

	result := response.GetHttpContentString()
	fmt.Println(result)
	return result, nil
}

// ParseOCRResult 解析阿里OCR返回结果
// ocrResult: 阿里OCR返回的JSON字符串
// cardType: 身份证类型，"id-card-front" 表示正面，"id-card-back" 表示反面
// 返回结构化的OCR识别结果
func ParseOCRResult(ocrResult, cardType string) (*OCRResult, error) {
	// 解析 JSON 字符串
	var result map[string]interface{}
	err := json.Unmarshal([]byte(ocrResult), &result)
	if err != nil {
		return nil, err
	}

	// 提取 OCR 识别结果
	oCRResultStruct := &OCRResult{}

	// 根据卡片类型处理不同的解析逻辑
	switch cardType {
	case "id-card-front":
		// 处理身份证正面
		if data, ok := result["data"].([]interface{}); ok && len(data) > 0 {
			if item, ok := data[0].(map[string]interface{}); ok {
				if ocrData, ok := item["ocr"].(map[string]interface{}); ok {
					if idCardInfo, ok := ocrData["idCard"].(map[string]interface{}); ok {
						if name, ok := idCardInfo["name"].(string); ok {
							oCRResultStruct.RealName = name
						}
						if id, ok := idCardInfo["id"].(string); ok {
							oCRResultStruct.IdCard = id
						}
						if birthday, ok := idCardInfo["birthday"].(string); ok {
							oCRResultStruct.Birthday = birthday
						}
						if gender, ok := idCardInfo["gender"].(string); ok {
							oCRResultStruct.Gender = gender
						}
						if address, ok := idCardInfo["address"].(string); ok {
							oCRResultStruct.Address = address
						}
					}
				}
			}
		}
	case "id-card-back":
		// 处理身份证反面
		if data, ok := result["data"].([]interface{}); ok && len(data) > 0 {
			if item, ok := data[0].(map[string]interface{}); ok {
				if ocrData, ok := item["ocr"].(map[string]interface{}); ok {
					if _, ok := ocrData["idCard"].(map[string]interface{}); ok {
						// 身份证反面可能包含有效期等信息
						// 这里可以根据实际需要进行解析 ，例如：
						// if validDate, ok := idCardInfo["validDate"].(string); ok {
						// 	oCRResultStruct.ValidDate = validDate
						// }
					}
				}
			}
		}
	case "student-card":
		// 处理学生证
		if data, ok := result["data"].([]interface{}); ok && len(data) > 0 {
			if item, ok := data[0].(map[string]interface{}); ok {
				if ocrData, ok := item["ocr"].(map[string]interface{}); ok {
					if studentCardInfo, ok := ocrData["studentCard"].(map[string]interface{}); ok {
						if schoolName, ok := studentCardInfo["schoolName"].(string); ok {
							oCRResultStruct.SchoolName = schoolName
						}
						if studentId, ok := studentCardInfo["studentId"].(string); ok {
							oCRResultStruct.StudentId = studentId
						}
					}
				}
			}
		}
	default:
		// 处理其他类型的卡片
	}

	return oCRResultStruct, nil
}

// GetTimestamp 获取当前时间戳
func GetTimestamp() int64 {
	return time.Now().UnixNano() / 1000000
}
