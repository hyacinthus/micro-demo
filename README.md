# Golang 微服务项目模板

这个模板主要使用 echo 提供 restful api

结合 repo x 使用

## 非业务部分

有一些地方实在抽象不出去，复制项目后需要核对修改

- config 配置
- store 初始化数据库和 redis
- util 自定义类型

## 业务部分

entity 是个 demo，每个实体一个文件写下去就好

## 注意事项

- endpoint 要不要认证在路由中根据分组决定
- 认证的接口必定能在 context 中取到相关字段
- docs 在构建时生成，本地调试需要在本地用 swag init 生成一份，不然会报错
