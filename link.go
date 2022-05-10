package bls_tss
/*
#cgo LDFLAGS:-lbls_tss -lm -ldl
#cgo linux,amd64 LDFLAGS:-L${SRCDIR}/bls-tss/lib/linux/amd64
#cgo linux,arm64 LDFLAGS:-L${SRCDIR}/bls-tss/lib/linux/arm64
#cgo darwin,amd64 LDFLAGS:-L${SRCDIR}/bls-tss/lib/darwin/amd64
#cgo darwin,arm64 LDFLAGS:-L${SRCDIR}/bls-tss/lib/darwin/arm64
*/
import "C"