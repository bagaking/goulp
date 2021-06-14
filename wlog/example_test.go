package wlog

import (
	"context"
	"testing"
)

func TestExample(t *testing.T) {
	LDev.Log().Debug("dev msg 1")
	wlog, err := NewWLog(createStderrLogger())
	if err != nil {
		panic(err)
	}

	wlog.Common().Info("dev msg 2")

	Common("ok").Info("print by default \\ wlog instance")
	Common("ok").Dev().Info("print by dev \\ wlog instance")

	ctx := context.Background()
	ByCtx(ctx, "l1").Info("print by byCtx entry")
}

func TestMFP(t *testing.T) {
	ctx := context.Background()
	l1, ctx1 := ByCtxAndCache(ctx, "l1")
	l1.Info("l1")

	l2, ctx2 := ByCtxAndCache(ctx1, "l2")
	l2.Info("l2")

	l3, ctx3 := ByCtxAndRemoveCache(ctx2, "l3")
	l3.Info("l3")

	l4, _ := ByCtxAndRemoveCache(ctx3, "l4")
	l4.Info("l4")
}

func BenchmarkExample(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Common("ok").Info("print by default wlog instance")
	}
}

func BenchmarkDisableExample(b *testing.B) {
	DevEnabled = false
	for i := 0; i < b.N; i++ {
		Common("ok").Dev().Info("print by default wlog instance")
	}
}

func BenchmarkDisableExample2(b *testing.B) {
	DevEnabled = false
	d := Common("ok").Dev()
	for i := 0; i < b.N; i++ {
		d.Info("print by default wlog instance")
	}
}
