package main

import (
	_ "embed"
	"encoding/xml"
	"fmt"
	"runtime"

	"kubevirt.io/kubevirt/pkg/virt-launcher/virtwrap/api"
	"libvirt.org/go/libvirtxml"
)

//go:embed testdata/vmi-pvc-clipboard.xml
var libvirtDomainXML []byte

func main() {
	var m1, m2, m3 runtime.MemStats
	var s1 api.DomainSpec
	var s2 libvirtxml.Domain

	runtime.ReadMemStats(&m1)

	if err := xml.Unmarshal(libvirtDomainXML, &s1); err != nil {
		panic(err)
	}
	runtime.ReadMemStats(&m2)

	if err := xml.Unmarshal(libvirtDomainXML, &s2); err != nil {
		panic(err)
	}
	runtime.ReadMemStats(&m3)

	fmt.Printf("[Kubevirt  ] total=%d, mallocs=%d\n", m2.TotalAlloc-m1.TotalAlloc, m2.Mallocs-m1.Mallocs)
	fmt.Printf("[Libvirtxml] total=%d, mallocs=%d\n", m3.TotalAlloc-m2.TotalAlloc, m3.Mallocs-m2.Mallocs)
}
