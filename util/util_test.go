package util

/**
 * 函数测试： go test -v
 * 性能测试：go test -test.bench=".*"
 * 生成调用函数性能图：go test -bench=".*" -cpuprofile=cpu.profile ./util
 * 查看性能数据：go tool pprof util.test cpu.profile
 * 查看测试覆盖率：
*/
import "testing"

func TestGenShortId(t *testing.T) {
	shortId, err := GenShortId()

	if shortId == "" || err != nil {
		t.Error("GenShortId failed!")
	}

	t.Log("GenShortId test pass")
}

func BenchmarkGenShortId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenShortId()
	}
}

func BenchmarkGenShortIdTimeConsuming(b *testing.B) {
	b.StopTimer()

	shortId, err := GenShortId()
	if shortId == "" || err != nil {
		b.Error(err)
	}

	b.StartTimer() // 重新开始时间

	for i := 0; i < b.N; i++ {
		GenShortId()
	}
}