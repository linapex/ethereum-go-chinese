
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450098914398208>

package metrics

//HealthChecks保存一个描述任意上/下状态的错误值。
type Healthcheck interface {
	Check()
	Error() error
	Healthy()
	Unhealthy(error)
}

//new healthcheck构造一个新的healthcheck，它将使用给定的
//函数更新其状态。
func NewHealthcheck(f func(Healthcheck)) Healthcheck {
	if !Enabled {
		return NilHealthcheck{}
	}
	return &StandardHealthcheck{nil, f}
}

//nilhealthcheck是不允许的。
type NilHealthcheck struct{}

//支票是不允许的。
func (NilHealthcheck) Check() {}

//错误是不可操作的。
func (NilHealthcheck) Error() error { return nil }

//健康是禁忌。
func (NilHealthcheck) Healthy() {}

//不健康是禁忌。
func (NilHealthcheck) Unhealthy(error) {}

//StandardHealthCheck是HealthCheck的标准实现，
//存储状态和用于调用以更新状态的函数。
type StandardHealthcheck struct {
	err error
	f   func(Healthcheck)
}

//check运行healthcheck函数更新healthcheck的状态。
func (h *StandardHealthcheck) Check() {
	h.f(h)
}

//错误返回HealthCheck的状态，如果它是健康的，则为零。
func (h *StandardHealthcheck) Error() error {
	return h.err
}

//健康将健康检查标记为健康。
func (h *StandardHealthcheck) Healthy() {
	h.err = nil
}

//不健康将健康检查标记为不健康。存储错误并
//可以通过错误方法检索。
func (h *StandardHealthcheck) Unhealthy(err error) {
	h.err = err
}

