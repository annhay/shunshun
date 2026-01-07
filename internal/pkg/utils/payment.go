package utils

import (
	"fmt"

	"github.com/smartwalle/alipay/v3"
)

// 支付

// AliPay 支付宝支付
func AliPay(orderCode, totalAmount string) string {
	var appId = "2021000149636760"
	var privateKey = "MIIEpAIBAAKCAQEAjfA0d6gruqjDnyUjIBFbT3jkx2VYs4/rZU0JR1u2N8h7oJ5PNKHZReFfHyQ9tEAaU41/NmWqJTNtbtO+RZ6wJenOCDvBmCkUWJSu3hw1/Fn4AajHn/SHXZ3Ml6WIy1sC/3KE8eHS6NLhF4FgFgc9DvH526aICIEaPbu6VMt6Y/nsRf1m6gPaN7xDGWBsqVsu0MPQ0N8nygQUUBVcnz5rEzI/7T6K86NlRIpRJOkqcIpqcXN/STGYFoN3ikc9hnMut8QrDHxTulLVHNKa/5yEeQaFF8Yw+yhGVu2Nc6ArW7DBhosr0q5dcrX38c2pmbFPIwph1T70uxX1+/mMONqsAQIDAQABAoIBAAqwJtNH15sjsC9gtYdppy2R1fBp4kcLNFeZeHRmJI+Iyj5rDV3SPjEz6lzG9tqG5TSbeBPZjfllKP1qdm55p5wDQh9+mHJjzYNqFszk5O/Ouo0tb3LNEBBtnIVi0q01ekFQF1C7h40+q/KALIMcIm3orL7siFvTlO1HIJ3YAKxcfOQGdmtROqKkHuLJGAySYng3p/OqEez7T05jYBVcK3UQNgXZZt7EtHk/V/GLwANQWUyqgBXHWBvDdkGZpW+JYNYAOjQZQDUIOMf6CgvHvvbCpsEhy9Ol1iK4euQyxdoCr4z4yiAdbsdoMWeS2yENzBn945HvXjItJcLcFzVzW00CgYEA9F501tlOc7lrRhz537dAX5/ghWgh1f3cT3m5QkKMMhmzDHrOdH/2/sYZ92iPadOaYhifDDIKOv4Gd2VLk8fIVoTQJeQzAyn3a8rfBtVCyK8KUZ3A3jq9adbNkYpRxxi6fikzFavZ0R94BrwBqB4RudRCyCQQ2ABRgjrrUKe4l4sCgYEAlLGrZjbNePepxRH7Oyg9Myxrj0vK4UeAGJQHa2dROhqO3NBO9FyPPALiNHGdkzg7ECCbw4h6G++U6eMbkF5l76BnozOfkN+ZTqV6gLjO7lcaeBDXx5n++VBRE8pNgddYvlMbcLqx9BULyhvl6BGzIb4aQkr9p57x5bjDbSF4XCMCgYEAyIEcIxEIUuGniE7MI2iLtCpNIYkQgjGaa8d3X0uVFqKJi8rTzTkV43ON6LdtPKq3uJd1IJ+KT18Q1TRS771zvrGYzA5SYN01OsepeUTQWDNvJwpmLrFJqybpYup4MQE0O8H4PWbVAMZuSDBIt7V8W9oytV8KRwDz4AQSAgqr5gMCgYA8BvGzxOH0OL8vkI/ElP0H4KHXaniPs4ax5WiNYls3Qqtz1yBYo9krF9rr4wYC/ctSOmfHaxwolPKf7RAemw05zJ6qEtgS60F/r2wh9PmM3FsSJ3KE4NU/Hr5sZ9ocVaw8wV4thyD58VkeEV8h7atMLut44b8+4Pq0i39RWha10wKBgQCyKCamCoB7Lhjn7aYh2RBdmLh1tk/zDNBPtZryvQfXWD8o5ENqqpfjQxIJ69EEGt0/Y6WSF9FhSXYlJxecgGTygXGArFigacCUhO3dtgcdTDr+iBrzYgLnf6w3LogiCDpWGROsca+vnVV8WzF9dEiqN9uuvmu6f9aZm5YmhpyjGA==" // 必须，上一步中使用 RSA签名验签工具 生成的私钥
	var client, err = alipay.New(appId, privateKey, false)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	var p = alipay.TradeWapPay{}
	p.NotifyURL = "http://xxx"
	p.ReturnURL = "http://xxx"
	p.Subject = "支付宝支付"
	p.OutTradeNo = orderCode
	p.TotalAmount = totalAmount
	p.ProductCode = "QUICK_WAP_WAY"
	url, err := client.TradeWapPay(p)
	if err != nil {
		fmt.Println(err)
	}
	// 这个 payURL 既是用于打开支付宝支付页面的 URL，可将输出的内容复制，到浏览器中访问该 URL 即可打开支付页面。
	var payURL = url.String()
	return payURL
}
