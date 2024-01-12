# libvirtxml-memory-footprint

Compares memory usage from [libvirtxml][] with KubeVirt's [schema.go][]
Motivation: https://github.com/kubevirt/kubevirt/issues/10844

# Run

```
go run main.go --xml testdata/real-use-cases/vmi-pvc-clipboard.xml

unsafe.Sizeof() of each domain struct
 * libvirtxml: 520 bytes
 * kubevirt  : 920 bytes

compared:bytes is the percentage of libvirtxml in comparison to KubeVirt
e.g: Value: 111, means livirtxml's domain struct is 11% bigger than KubeVirt domain struct

                                         |                     | <------- libvirtxml ----> | <-------- kubevirt -----> |
                                         | compared: bytes (%) |  benchmark  |   runtime   |  benchmark  |   runtime   |
                                filename | benchmark | runtime | KiB | alloc | KiB | alloc | KiB | alloc | KiB | alloc |
                   vmi-pvc-clipboard.xml |       111 |     111 | 113 |  2970 | 113 |  2970 | 101 |  2867 | 101 |  2867 |
                                         |                     | <------- libvirtxml ----> | <-------- kubevirt -----> |
                                         | compared: bytes (%) |  benchmark  |   runtime   |  benchmark  |   runtime   |
                                filename | benchmark | runtime | KiB | alloc | KiB | alloc | KiB | alloc | KiB | alloc |
The average size of Libvirt : 113 KiB
The average size of Kubevirt: 101 KiB
In average, libvirtxml is 111 % compared to Kubevirt
```

# Results

* [testdata/libvirt-git-xml][]
* [testdata/real-use-cases][]

[libvirtxml]: https://gitlab.com/libvirt/libvirt-go-xml-module
[schema.go]: https://github.com/kubevirt/kubevirt/blob/main/pkg/virt-launcher/virtwrap/api/schema.go
[testdata/libvirt-git-xml]: https://gist.github.com/victortoso/fa420dd6bc385e10a991c76c30ed521b#file-from-libvirt-git-asc
[testdata/real-use-cases]: https://gist.github.com/victortoso/fa420dd6bc385e10a991c76c30ed521b#file-real-use-cases-asc
