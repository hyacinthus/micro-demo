# Golang 微服务项目模板

这个模板主要使用 echo 提供 restful api

结合 repo x 使用

## 非业务部分

有一些地方实在抽象不出去，复制项目后需要核对修改

- config 配置
- store 用 config 初始化 mysql 和 redis , 顺便用 gorm 建表。
- util 需要用到 config 的自定义类型

## 业务部分

- entity 是个 demo，每个实体一个文件写下去就好
- entity 的模型在子 package ske 中，为了模型可以导出被别的服务使用
- 在 package ske 中，可以直接将模型的纯函数还有 sdk 和模型写在一起

## 其他

- endpoint 要不要认证在路由中根据分组决定
- 认证的接口必定能在 context 中取到相关字段

## 微服务支持

为了支持微服务，并保持代码简洁不做过多分层，做了如下工作：

- 服务间访问依然使用 http 协议，这样代码量少
- 将业务模型、模型的方法、sdk 放在 ske 这个 package
- 别的语言直接通过内部 rest api 访问，为 golang 提供 sdk
- 如果追求极致性能，服务间访问可切换至 grpc , 只需要修改 ske package 的内容即可

## Dockerfile

- 使用两段构建，先 build ，再把二级制文件复制到生产镜像
- 生产镜像使用了我略微修改的
  [debian](https://github.com/hyacinthus/docker-debian)
  - 增加证书，用以在镜像内访问别的 https 资源，否则会报 x509
  - 指定了中国时区
