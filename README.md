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

## 注意事项

- endpoint 要不要认证在路由中根据分组决定
- 认证的接口必定能在 context 中取到相关字段
