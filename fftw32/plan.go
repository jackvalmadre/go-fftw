package fftw32

// #include <fftw3.h>
import "C"

import "unsafe"

type Plan struct {
	fftw_p C.fftwf_plan
}

func (p *Plan) Execute() *Plan {
	C.fftwf_execute(p.fftw_p)
	return p
}

func (p *Plan) Destroy() {
	C.fftwf_destroy_plan(p.fftw_p)
}

func NewPlan(in, out *Array, dir Direction, flag Flag) *Plan {
	// TODO: check that len(in) == len(out)
	n := in.Len()
	var (
		n_    = C.int(n)
		in_   = (*C.fftwf_complex)(unsafe.Pointer(in.ptr()))
		out_  = (*C.fftwf_complex)(unsafe.Pointer(out.ptr()))
		dir_  = C.int(dir)
		flag_ = C.uint(flag)
	)
	p := C.fftwf_plan_dft_1d(n_, in_, out_, dir_, flag_)
	return &Plan{p}
}

func NewPlan2(in, out *Array2, dir Direction, flag Flag) *Plan {
	// TODO: check that in and out have the same dimensions
	n0, n1 := in.Dims()
	var (
		n0_   = C.int(n0)
		n1_   = C.int(n1)
		in_   = (*C.fftwf_complex)(unsafe.Pointer(in.ptr()))
		out_  = (*C.fftwf_complex)(unsafe.Pointer(out.ptr()))
		dir_  = C.int(dir)
		flag_ = C.uint(flag)
	)
	p := C.fftwf_plan_dft_2d(n0_, n1_, in_, out_, dir_, flag_)
	return &Plan{p}
}

func NewPlan3(in, out *Array3, dir Direction, flag Flag) *Plan {
	// TODO: check that in and out have the same dimensions
	n0, n1, n2 := in.Dims()
	var (
		n0_   = C.int(n0)
		n1_   = C.int(n1)
		n2_   = C.int(n2)
		in_   = (*C.fftwf_complex)(unsafe.Pointer(in.ptr()))
		out_  = (*C.fftwf_complex)(unsafe.Pointer(out.ptr()))
		dir_  = C.int(dir)
		flag_ = C.uint(flag)
	)
	p := C.fftwf_plan_dft_3d(n0_, n1_, n2_, in_, out_, dir_, flag_)
	return &Plan{p}
}
