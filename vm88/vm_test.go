package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/vmware/govmomi/vim25/types"
)

var vmAllListTests = []struct {
	Name   string
	System string
	Self   Self
	VM     types.ManagedObjectReference
}{
	{"测试机器two", "Red Hat Enterprise Linux 7 (64-bit)", Self{Type: "VirtualMachine", Value: "vm-904"},
		types.ManagedObjectReference{Type: "VirtualMachine", Value: "vm-904"}},
	{"EVE-PRO(100.222)", "Ubuntu Linux (64-bit)", Self{Type: "VirtualMachine", Value: "vm-902"},
		types.ManagedObjectReference{Type: "VirtualMachine", Value: "vm-902"}},
}

func TestVmWare_GetAllVmClient(t *testing.T) {
	vc := os.Getenv("GOVMOMI_URL")
	user := os.Getenv("GOVMOMI_USERNAME")
	pwd := os.Getenv("GOVMOMI_PASSWORD")
	vm := NewVmWare(vc, user, pwd)
	fmt.Println(vm)
	// vm := NewVmWare("192.168.100.200", "Administrator@vsphere.local", "!@AsiaLink@2020")
	vmList, _, _ := vm.GetAllVmClient()
	for _, vm := range vmList {
		// for _, vmtest := range vmAllListTests {
		// 	if vm == vmtest {
		// 		t.Log("获取虚拟机测试通过")
		// 	}
		// }
		fmt.Println(vm)
	}
}

func TestVmWare_GetAllHost(t *testing.T) {

}
