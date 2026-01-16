package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	gourl "net/url"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
)

// 腾讯云实名认证
func calcAuthorization(secretId string, secretKey string) (auth string, datetime string, err error) {
	timeLocation, _ := time.LoadLocation("Etc/GMT")
	datetime = time.Now().In(timeLocation).Format("Mon, 02 Jan 2006 15:04:05 GMT")
	signStr := fmt.Sprintf("x-date: %s", datetime)

	// hmac-sha1
	mac := hmac.New(sha1.New, []byte(secretKey))
	mac.Write([]byte(signStr))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	auth = fmt.Sprintf("{\"id\":\"%s\", \"x-date\":\"%s\", \"signature\":\"%s\"}",
		secretId, datetime, sign)

	return auth, datetime, nil
}

func urlencode(params map[string]string) string {
	var p = gourl.Values{}
	for k, v := range params {
		p.Add(k, v)
	}
	return p.Encode()
}

// VerifyIdCard 身份证验证
// realName: 真实姓名
// idCardNo: 身份证号
// 返回验证结果和错误信息
func VerifyIdCard(realName, idCardNo string) (bool, error) {
	// 云市场分配的密钥 Id
	secretId := os.Getenv("QySi8XnshaDPEfby")
	// 云市场分配的密钥 Key
	secretKey := os.Getenv("pCpcU2LxKaYZv6eLmpapB25RWUNoYvGy")

	// 签名
	auth, _, err := calcAuthorization(secretId, secretKey)
	if err != nil {
		return false, err
	}

	// 请求方法
	method := "POST"
	reqID, err := uuid.GenerateUUID()
	if err != nil {
		return false, err
	}
	// 请求头
	headers := map[string]string{"Authorization": auth, "request-id": reqID}

	// 查询参数
	queryParams := make(map[string]string)

	// body 参数
	bodyParams := make(map[string]string)
	bodyParams["cardNo"] = idCardNo
	bodyParams["realName"] = realName
	bodyParamStr := urlencode(bodyParams)
	// url 参数拼接
	url := "https://ap-beijing.cloudmarket-apigw.com/service-18c38npd/idcard/VerifyIdcardv2"

	if len(queryParams) > 0 {
		url = fmt.Sprintf("%s?%s", url, urlencode(queryParams))
	}

	bodyMethods := map[string]bool{"POST": true, "PUT": true, "PATCH": true}
	var body io.Reader = nil
	if bodyMethods[method] {
		body = strings.NewReader(bodyParamStr)
		headers["Content-Type"] = "application/x-www-form-urlencoded"
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return false, err
	}
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	response, err := client.Do(request)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	// 这里需要根据实际的返回结果进行解析
	// 由于没有实际的返回格式示例，这里只提供一个基本的实现
	// 实际使用时需要根据云市场的API文档进行调整
	result := string(bodyBytes)
	fmt.Println("身份证验证结果:", result)

	// 假设返回结果中包含 "success" 表示验证成功
	if strings.Contains(result, "success") {
		return true, nil
	}

	return false, fmt.Errorf("身份证验证失败: %s", result)
}

// StudentVerification 学生认证
// realName: 真实姓名
// studentId: 学号
// schoolName: 学校名称
// 返回验证结果和错误信息
func StudentVerification(realName, studentId, schoolName string) (bool, error) {
	// 这里实现学生认证的逻辑
	// 实际使用时，需要调用相应的API进行验证
	fmt.Printf("学生认证: 姓名=%s, 学号=%s, 学校=%s\n", realName, studentId, schoolName)

	// 暂时返回成功，实际使用时需要根据API返回结果进行判断
	return true, nil
}

// FaceRecognition 人脸识别
// faceImageUrl: 人脸图片URL
// idCardFrontUrl: 身份证正面图片URL
// 返回验证结果和错误信息
func FaceRecognition(faceImageUrl, idCardFrontUrl string) (bool, error) {
	// 这里实现人脸识别的逻辑
	// 实际使用时，需要调用相应的API进行验证
	fmt.Printf("人脸识别: 人脸图片=%s, 身份证图片=%s\n", faceImageUrl, idCardFrontUrl)

	// 暂时返回成功，实际使用时需要根据API返回结果进行判断
	return true, nil
}
