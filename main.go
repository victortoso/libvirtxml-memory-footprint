package main

import (
	_ "embed"
	"encoding/xml"
	"fmt"
	"runtime"
	"testing"

	"kubevirt.io/kubevirt/pkg/virt-launcher/virtwrap/api"
	"libvirt.org/go/libvirtxml"
)

//go:embed testdata/vmi-pvc-clipboard.xml
var libvirtDomainXML []byte

func BenchmarkLibvirtxml(b *testing.B) {
	domain := new(libvirtxml.Domain)
	if err := xml.Unmarshal(libvirtDomainXML, domain); err != nil {
		panic(err)
	}
}

func BenchmarkKubevirtSchema(b *testing.B) {
	domain := new(api.DomainSpec)
	if err := xml.Unmarshal(libvirtDomainXML, domain); err != nil {
		panic(err)
	}
}

func do_test_benchmark() {
	res1 := testing.Benchmark(BenchmarkLibvirtxml)
	fmt.Printf("Libvirtxml:\n")
	fmt.Printf("Memory allocations  : %d \n", res1.MemAllocs)
	fmt.Printf("Number of bytes ... : %d \n", res1.MemBytes)
	fmt.Printf("Number of run ..... : %d \n", res1.N)
	fmt.Printf("Time taken ........ : %s \n\n", res1.T)

	res2 := testing.Benchmark(BenchmarkKubevirtSchema)
	fmt.Printf("Kubevirt Schema:\n")
	fmt.Printf("Memory allocations : %d \n", res2.MemAllocs)
	fmt.Printf("Number of bytes .. : %d \n", res2.MemBytes)
	fmt.Printf("Number of run .... : %d \n", res2.N)
	fmt.Printf("Time taken ....... : %s \n\n", res2.T)

	fmt.Printf("+ Libvirtxml > Schema by %d bytes\n\n", res1.MemBytes-res2.MemBytes)
}

func do_runtime_check() {
	var m1, m2, m3 runtime.MemStats
	var s1 libvirtxml.Domain
	var s2 api.DomainSpec

	runtime.ReadMemStats(&m1)
	if err := xml.Unmarshal(libvirtDomainXML, &s1); err != nil {
		panic(err)
	}
	runtime.ReadMemStats(&m2)
	if err := xml.Unmarshal(libvirtDomainXML, &s2); err != nil {
		panic(err)
	}
	runtime.ReadMemStats(&m3)

	total1 := m2.TotalAlloc - m1.TotalAlloc
	total2 := m3.TotalAlloc - m2.TotalAlloc
	fmt.Printf("[Libvirtxml] total=%d, mallocs=%d\n", total1, m2.Mallocs-m1.Mallocs)
	fmt.Printf("[Kubevirt  ] total=%d, mallocs=%d\n\n", total2, m3.Mallocs-m2.Mallocs)
	fmt.Printf("+ Libvirtxml > Schema by %d bytes\n\n", total1-total2)
}

func main() {
	fmt.Println("* Using testing.Benchmark\n")
	do_test_benchmark()

	fmt.Println("* Using runtime.ReadMemStats\n")
	do_runtime_check()
}
