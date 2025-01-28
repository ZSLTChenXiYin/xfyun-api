# xfyun-api #
本项目为Golang开发者提供了一个讯飞云服务的接口库。
## 内容导引 ##
* [介绍](#介绍)
* [里程碑](#里程碑)
* [版本实现](#版本实现)
* [部署](#部署)
* [快速上手](#快速上手)
* [问题反馈](#问题反馈)
## 介绍 ##
#### 如何使用讯飞云服务
* [讯飞云开放平台](https://www.xfyun.com/)
* [讯飞云服务WebAPI](https://www.xfyun.cn/doc/)
#### 目前为止，xfyun-api 可以提供以下支持：
* 人脸识别服务调用（已齐全）
## 里程碑 ##
* 实现各服务的基本调用接口 < latest
## 版本实现 ##
#### 测试版
* master 仅实现人脸识别服务接口
#### 正式版
* 无
## 部署 ##
xfyun-api 的部署依赖 Go modules，如果你还没有 go mod，你需要首先初始化:
```sh
go mod init myproject
```
安装 xfyun-api
```sh
go get -u github.com/ZSLTChenXiYin/xfyun-api
```
## 快速上手 ##
* 请参考 [api_test](https://github.com/ZSLTChenXiYin/xfyun-api/tree/master/api_test) 中的示例
## 问题反馈 ##
* 陈汐胤会在每周五至周日查看 [Issues](https://github.com/ZSLTChenXiYin/xfyun-api/issues)，还会不定期地在 bilibili 直播
>> 陈汐胤的 e-mail: imjfoy@163.com
>> 
>> 陈汐胤的 bilibili UID: 352456302
