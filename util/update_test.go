package util_test

import (
	"flight-api/util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// --- helpers ---

func mustPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic, got none")
		}
	}()
	f()
}

// --- UpdateString ---

func TestUpdateString(t *testing.T) {
	t.Run("valid - update with non-nil src", func(t *testing.T) {
		dst := new(string)
		*dst = "old"

		src := "new"
		util.UpdateString(&dst, &src)

		assert.Equal(t, "new", *dst)
	})

	t.Run("valid - src is nil (dst unchanged)", func(t *testing.T) {
		dst := new(string)
		*dst = "keep"

		util.UpdateString(&dst, nil)

		assert.Equal(t, "keep", *dst) // tidak berubah
	})

	t.Run("invalid - wrong dereference level", func(t *testing.T) {
		str := "oops"
		// dst seharusnya **string, tapi kita coba pakai *string
		// ini tidak akan compile kalau langsung dipanggil,
		// jadi kita uji dengan interface{} + reflect, atau tandai sebagai negative test.
		// Di sini cukup tunjukkan bahwa tidak bisa compile.
		_ = str
		// 👉 Compile-time error: UpdateString(str, &str) // uncomment to see the error
	})
}

// --- UpdateBool ---
func TestUpdateBool(t *testing.T) {
	t.Run("valid - update with non-nil src", func(t *testing.T) {
		dst := new(bool)
		*dst = false

		src := true
		util.UpdateBool(&dst, &src)

		if *dst != true {
			t.Fatalf("expected dst=true, got %v", *dst)
		}
	})

	t.Run("valid - src is nil (dst unchanged)", func(t *testing.T) {
		dst := new(bool)
		*dst = false

		util.UpdateBool(&dst, nil)

		if *dst != false {
			t.Fatalf("expected dst remain false, got %v", *dst)
		}
	})

	t.Run("invalid - dst is nil pointer to pointer", func(t *testing.T) {
		var dstPtr **bool = nil
		src := true

		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("expected panic, got none")
			}
		}()

		util.UpdateBool(dstPtr, &src) // panic karena dstPtr == nil
	})
}

// --- TestInt ---

func TestUpdateInt(t *testing.T) {
	t.Run("valid - update with non-nil src - int", func(t *testing.T) {
		dst := []*int{new(int), nil}

		src := 42

		// dst has value
		*dst[0] = 2
		old := *dst[0]
		util.UpdateInt(&dst[0], &src)
		assert.IsType(t, 42, *dst[0])
		assert.Equal(t, 42, *dst[0])
		assert.NotEqual(t, old, *dst[0])

		// dst is nil
		util.UpdateInt(&dst[1], &src)
		assert.IsType(t, 42, *dst[1])
		assert.Equal(t, 42, *dst[1])
		assert.NotEqual(t, nil, *dst[1])
	})

	t.Run("valid - update with non-nil src - int32", func(t *testing.T) {
		dst := []*int32{new(int32), nil}

		src := int32(42)

		// dst has value
		*dst[0] = 2
		old := *dst[0]
		util.UpdateInt(&dst[0], &src)
		assert.IsType(t, int32(42), *dst[0])
		assert.Equal(t, int32(42), *dst[0])
		assert.NotEqual(t, old, *dst[0])

		// dst is nil
		util.UpdateInt(&dst[1], &src)
		assert.IsType(t, int32(42), *dst[1])
		assert.Equal(t, int32(42), *dst[1])
		assert.NotEqual(t, nil, *dst[1])
	})

	t.Run("valid - update with non-nil src - int64", func(t *testing.T) {
		dst := []*int64{new(int64), nil}

		src := int64(410000123456789)

		// dst has value
		*dst[0] = 2
		old := *dst[0]
		util.UpdateInt(&dst[0], &src)
		assert.IsType(t, int64(410000123456789), *dst[0])
		assert.Equal(t, int64(410000123456789), *dst[0])
		assert.NotEqual(t, old, *dst[0])

		// dst is nil
		util.UpdateInt(&dst[1], &src)
		assert.IsType(t, int64(410000123456789), *dst[1])
		assert.Equal(t, int64(410000123456789), *dst[1])
		assert.NotEqual(t, nil, *dst[1])
	})

	t.Run("valid - src is nil (dst unchanged)", func(t *testing.T) {
		dst := new(int)
		*dst = 42

		util.UpdateInt(&dst, nil)

		assert.IsType(t, 42, *dst)
		assert.Equal(t, 42, *dst)
		assert.NotEqual(t, nil, *dst)
	})

	t.Run("invalid - dst is nil pointer to pointer", func(t *testing.T) {
		var dstPtr **int = nil
		src := 42

		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("expected panic, got none")
			}
		}()

		util.UpdateInt(dstPtr, &src)

		assert.IsType(t, nil, dstPtr)
		assert.Equal(t, nil, dstPtr)
	})
}

func TestUpdateTime(t *testing.T) {
	t.Run("valid - update with non-nil src", func(t *testing.T) {
		now := time.Now()
		dst := &now

		src := time.Now().Add(2 * time.Hour)
		util.UpdateTime(&dst, &src)

		assert.IsType(t, src, *dst)
		assert.Equal(t, src, *dst)
		assert.NotEqual(t, now, *dst)
	})

	t.Run("valid - src is nil (dst unchanged)", func(t *testing.T) {
		now := time.Now()
		dst := &now

		util.UpdateTime(&dst, nil)

		assert.IsType(t, now, *dst)
		assert.Equal(t, now, *dst)
		assert.NotEqual(t, nil, *dst)
	})

	t.Run("invalid - dst is nil pointer to pointer", func(t *testing.T) {
		var dstPtr **time.Time = nil
		src := time.Now()

		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("expected panic, got none")
			}
		}()

		util.UpdateTime(dstPtr, &src)

		assert.IsType(t, nil, dstPtr)
		assert.Equal(t, nil, dstPtr)
	})
}
