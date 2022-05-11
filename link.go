package bls_tss
/*
#cgo LDFLAGS:-lbls_tss -lm -ldl
#cgo windows,amd64 LDFLAGS: -lws2_32 -luserenv -lbcrypt
#cgo linux,amd64 LDFLAGS:-L${SRCDIR}/bls-tss/lib/linux/amd64
#cgo linux,arm64 LDFLAGS:-L${SRCDIR}/bls-tss/lib/linux/arm64
#cgo darwin,amd64 LDFLAGS:-L${SRCDIR}/bls-tss/lib/darwin/amd64
#cgo darwin,arm64 LDFLAGS:-L${SRCDIR}/bls-tss/lib/darwin/arm64
#cgo windows,amd64 LDFLAGS:-L${SRCDIR}/bls-tss/lib/windows/amd64
*/
import "C"