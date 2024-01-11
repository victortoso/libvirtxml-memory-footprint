package main

import (
	_ "embed"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
	"unsafe"

	"github.com/spf13/pflag"
	"kubevirt.io/kubevirt/pkg/virt-launcher/virtwrap/api"
	"libvirt.org/go/libvirtxml"
)

type memory struct {
	nbytes  uint64
	nallocs uint64
}

type stats struct {
	libvirt  memory
	kubevirt memory
}

type line struct {
	benchmark *stats
	runtime   *stats
}

var gbl_current_xmldata []byte

func BenchmarkLibvirtxml(b *testing.B) {
	domain := new(libvirtxml.Domain)
	if err := xml.Unmarshal(gbl_current_xmldata, domain); err != nil {
		panic(err)
	}
}

func BenchmarkKubevirtSchema(b *testing.B) {
	domain := new(api.DomainSpec)
	if err := xml.Unmarshal(gbl_current_xmldata, domain); err != nil {
		panic(err)
	}
}

func do_test_benchmark(xmldata []byte) *stats {
	gbl_current_xmldata = xmldata
	l := testing.Benchmark(BenchmarkLibvirtxml)
	k := testing.Benchmark(BenchmarkKubevirtSchema)
	return &stats{
		libvirt: memory{
			nbytes:  l.MemBytes,
			nallocs: l.MemAllocs,
		},
		kubevirt: memory{
			nbytes:  k.MemBytes,
			nallocs: k.MemAllocs,
		},
	}
}

func do_runtime_check(xmldata []byte) *stats {
	var m1, m2, m3 runtime.MemStats

	runtime.ReadMemStats(&m1)
	{
		var s1 libvirtxml.Domain
		if err := xml.Unmarshal(xmldata, &s1); err != nil {
			panic(err)
		}
	}
	runtime.ReadMemStats(&m2)
	{
		var s2 api.DomainSpec
		if err := xml.Unmarshal(xmldata, &s2); err != nil {
			panic(err)
		}
	}
	runtime.ReadMemStats(&m3)
	return &stats{
		libvirt: memory{
			nbytes:  m2.TotalAlloc - m1.TotalAlloc,
			nallocs: m2.Mallocs - m1.Mallocs,
		},
		kubevirt: memory{
			nbytes:  m3.TotalAlloc - m2.TotalAlloc,
			nallocs: m3.Mallocs - m2.Mallocs,
		},
	}
}

func do_unsize_of() {
	lv := unsafe.Sizeof(libvirtxml.Domain{})
	kv := unsafe.Sizeof(api.DomainSpec{})
	fmt.Printf("unsafe.Sizeof() of each domain struct\n")
	fmt.Printf(" * libvirtxml: %d bytes\n", lv)
	fmt.Printf(" * kubevirt  : %d bytes\n\n", kv)
}

func input() []string {
	var files []string
	var directory string
	var singleFile string
	pflag.StringVar(&singleFile, "xml", "", "a single libvirt domain")
	pflag.StringVar(&directory, "dir", "", "a directory with libvirt domain xmls")
	pflag.Parse()

	if singleFile != "" {
		files = append(files, singleFile)
	}

	if directory != "" {
		err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(info.Name(), ".xml") {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	}

	return files
}

func printHeader(footer bool) {
	if !footer {
		fmt.Printf("compared:bytes is the percentage of libvirtxml in comparison to KubeVirt\n")
		fmt.Printf("e.g: Value: 111, means livirtxml's domain struct is 11%% bigger than KubeVirt domain struct\n\n")
	}
	fmt.Printf("%40.0s | %19.0s | <------- libvirtxml ----> | <-------- kubevirt -----> |\n", "", "")
	fmt.Printf("%40.0s | compared: bytes (%%) |  benchmark  |   runtime   |  benchmark  |   runtime   |\n", "")
	fmt.Printf("%40.40s | benchmark | runtime | KiB | alloc | KiB | alloc | KiB | alloc | KiB | alloc |\n", "filename")
}

func printLine(path string, bench, runt *stats) {
	fmt.Printf("%40.40s | %9d | %7d | %3d | %5d | %3d | %5d | %3d | %5d | %3d | %5d |\n",
		filepath.Base(path),
		bench.libvirt.nbytes*100/bench.kubevirt.nbytes,
		runt.libvirt.nbytes*100/runt.kubevirt.nbytes,
		bench.libvirt.nbytes/1024, bench.libvirt.nallocs,
		runt.libvirt.nbytes/1024, runt.libvirt.nallocs,
		bench.kubevirt.nbytes/1024, bench.kubevirt.nallocs,
		runt.kubevirt.nbytes/1024, runt.kubevirt.nallocs,
	)
}

func main() {
	files := input()
	do_unsize_of()
	printHeader(false)

	lines := make(map[string]line)
	for _, file := range files {
		bytes, err := os.ReadFile(file)
		if err != nil {
			panic(err)
		}

		lines[file] = line{
			benchmark: do_test_benchmark(bytes),
			runtime:   do_runtime_check(bytes),
		}
	}

	// sort files by ascending order of percentage of the size comparison
	// of libvirtxml to kubevirt
	sort.SliceStable(files, func(i, j int) bool {
		l1, _ := lines[files[i]]
		l2, _ := lines[files[j]]
		p1 := l1.benchmark.libvirt.nbytes * 100 / l1.benchmark.kubevirt.nbytes
		p2 := l2.benchmark.libvirt.nbytes * 100 / l2.benchmark.kubevirt.nbytes
		return p1 <= p2
	})

	avg_libvirt := uint64(0)
	avg_kubevirt := uint64(0)
	for _, filename := range files {
		line, _ := lines[filename]
		printLine(filename, line.benchmark, line.runtime)
		avg_libvirt += line.benchmark.libvirt.nbytes
		avg_kubevirt += line.benchmark.kubevirt.nbytes
	}
	printHeader(true)
	avg_libvirt = avg_libvirt / uint64(len(files))
	avg_kubevirt = avg_kubevirt / uint64(len(files))

	fmt.Printf("The average size of Libvirt : %d KiB\n", avg_libvirt/1024)
	fmt.Printf("The average size of Kubevirt: %d KiB\n", avg_kubevirt/1024)
	fmt.Printf("In average, libvirtxml is %d %% compared to Kubevirt",
		avg_libvirt*100/avg_kubevirt,
	)
}
