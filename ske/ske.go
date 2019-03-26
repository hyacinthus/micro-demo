// Package ske xx服务模型及sdk
// 这个文件夹还用service的名字是因为别的服务引用多个service的模型时，可以区分开
package ske

var addr = "http://ske/"

// SetAddr 供使用者修改
func SetAddr(url string) {
	addr = url
}
