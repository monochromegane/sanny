package sanny

/*
#cgo CFLAGS: -mavx -std=c99
#include <stdio.h>
#include <stdlib.h>
#include <immintrin.h> //AVX: -mavx
void avx_add(const size_t n, float *x, float *y, float *z)
{
    static const size_t single_size = 8;
    const size_t end = n / single_size;
    __m256 *vz = (__m256 *)z;
    __m256 *vx = (__m256 *)x;
    __m256 *vy = (__m256 *)y;
    for(size_t i=0; i<end; ++i) {
      vz[i] = _mm256_add_ps(vx[i], vy[i]);
    }
}

void avx_sub(const size_t n, float *x, float *y, float *z)
{
    static const size_t single_size = 8;
    const size_t end = n / single_size;
    __m256 *vz = (__m256 *)z;
    __m256 *vx = (__m256 *)x;
    __m256 *vy = (__m256 *)y;
    for(size_t i=0; i<end; ++i) {
      vz[i] = _mm256_sub_ps(vx[i], vy[i]);
    }
}

void avx_mul(const size_t n, float *x, float *y, float *z)
{
    static const size_t single_size = 8;
    const size_t end = n / single_size;
    __m256 *vz = (__m256 *)z;
    __m256 *vx = (__m256 *)x;
    __m256 *vy = (__m256 *)y;
    for(size_t i=0; i<end; ++i) {
      vz[i] = _mm256_mul_ps(vx[i], vy[i]);
    }
}

float avx_dot(const size_t n, float *x, float *y)
{
    static const size_t single_size = 8;
    const size_t end = n / single_size;
    __m256 *vx = (__m256 *)x;
    __m256 *vy = (__m256 *)y;
    __m256 vsum = {0};
    for(size_t i=0; i<end; ++i) {
      vsum = _mm256_add_ps(vsum, _mm256_mul_ps(vx[i], vy[i]));
    }
    __attribute__((aligned(32))) float t[8] = {0};
    _mm256_store_ps(t, vsum);
    return t[0] + t[1] + t[2] + t[3] + t[4] + t[5] + t[6] + t[7];
}
*/
import "C"
import (
	"reflect"
	"unsafe"
)

func mmMalloc(size int) []float32 {
	ptr := C._mm_malloc((C.size_t)(C.sizeof_float*size), 32)
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(ptr)),
		Len:  size,
		Cap:  size,
	}
	goSlice := *(*[]float32)(unsafe.Pointer(&hdr))
	return goSlice
}

func mmFree(v []float32) {
	C._mm_free(unsafe.Pointer(&v[0]))
}

func avxAdd(size int, x, y, z []float32) {
	C.avx_add((C.size_t)(size), (*C.float)(&x[0]), (*C.float)(&y[0]), (*C.float)(&z[0]))
}

func avxMul(size int, x, y, z []float32) {
	C.avx_mul((C.size_t)(size), (*C.float)(&x[0]), (*C.float)(&y[0]), (*C.float)(&z[0]))
}

func avxSub(size int, x, y, z []float32) {
	C.avx_sub((C.size_t)(size), (*C.float)(&x[0]), (*C.float)(&y[0]), (*C.float)(&z[0]))
}

func avxDot(size int, x, y []float32) float32 {
	dot := C.avx_dot((C.size_t)(size), (*C.float)(&x[0]), (*C.float)(&y[0]))
	return float32(dot)
}
