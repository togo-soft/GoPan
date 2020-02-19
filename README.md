# GoPan

## 加密解密的探索

[最安全的网盘MEGA的加密原理](https://love.junzimu.com/archives/3230)

[有人介绍一下MEGA下载的原理么？](https://www.v2ex.com/t/62361)

[B端加解密-XMLHttpRequest Level 2 使用指南](http://www.ruanyifeng.com/blog/2012/09/xmlhttprequest_level_2.html)

[B端对象生成URL-URL.createObjectURL()](https://developer.mozilla.org/zh-CN/docs/Web/API/URL/createObjectURL)

## 性能测试结果

```bash
# 测试命令
go test -bench=xxx -run=None -benchmem
```

| 文件大小 | 加解密 | 总时间  |
| -------- | ------ | ------- |
| 70K      | 加密   | 0.013s  |
| 70K      | 解密   | 0.013s  |
| 500K     | 加密   | 0.018s  |
| 500K     | 解密   | 0.015s  |
| 1.5M     | 加密   | 0.027s  |
| 1.5M     | 解密   | 0.027s  |
| 10M      | 加密   | 0.081s  |
| 10M      | 解密   | 0.069s  |
| 50M      | 加密   | 0.273s  |
| 50M      | 解密   | 0.232s  |
| 100M     | 加密   | 0.587s  |
| 100M     | 解密   | 0.465s  |
| 500M     | 加密   | 2.859s  |
| 500M     | 解密   | 2.393s  |
| 1G       | 加密   | 7.102s  |
| 1G       | 解密   | 4.217s  |
| 2G       | 加密   | 14.577s |
| 2G       | 解密   | 9.619s  |
| 5G       | 加密   | 38.356s |
| 5G       | 解密   | 37.827s |

## 响应码

### JWT

- 异常
	TokenError = 400
- 头字段认证信息缺失
	TokenMissHeader = 4900
- token超时
	TokenExpired = 4901
- token格式错误
	TokenMalformed = 4902
- token匹配uid失败
	TokenTampered = 4903