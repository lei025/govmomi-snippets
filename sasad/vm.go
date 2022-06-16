package main

import (
	"context"
	"encoding/csv"
	"net/url"
	"os"
	"strconv"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/units"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"gorm.io/gorm"
	// "gorm.io/gorm"
)

var client *vim25.Client
var ctx = context.Background()

const (
	VSPHERE_IP       = "172.118.69.31"
	VSPHERE_USERNAME = "administrator@vsphere.local"
	VSPHERE_PASSWORD = "Huayun@123"
	Insecure         = true
)

// NewClient 链接vmware
func NewClient() *vim25.Client {

	u := &url.URL{
		Scheme: "https",
		Host:   VSPHERE_IP,
		Path:   "/sdk",
	}

	u.User = url.UserPassword(VSPHERE_USERNAME, VSPHERE_PASSWORD)
	client, err := govmomi.NewClient(ctx, u, Insecure)
	if err != nil {
		panic(err)
	}
	return client.Client
}

// VmsHost 主机结构体
type VmsHost struct {
	Name string
	Ip   string
}

// VmsHosts 主机列表结构体
type VmsHosts struct {
	VmsHosts []VmsHost
}

// NewVmsHosts 初始化结构体
func NewVmsHosts() *VmsHosts {
	return &VmsHosts{
		VmsHosts: make([]VmsHost, 10),
	}
}

// 虚拟机表
type Vm struct {
	gorm.Model
	Uuid       string `gorm:"type:varchar(40);not null;unique;comment:'虚拟机id'"`
	Vc         string `gorm:"type:varchar(30);comment:'Vcenter Ip'"`
	Esxi       string `gorm:"type:varchar(30);comment:'Esxi Id'"`
	Name       string `gorm:"type:varchar(90);comment:'Vm名字'"`
	Ip         string `gorm:"type:varchar(20);comment:'Vm ip'"`
	PowerState string `gorm:"type:varchar(20);comment:'Vm state'"`
}

// AddHost 新增主机
func (vmshosts *VmsHosts) AddHost(name string, ip string) {
	host := &VmsHost{name, ip}
	vmshosts.VmsHosts = append(vmshosts.VmsHosts, *host)
}

// SelectHost 查询主机ip
func (vmshosts *VmsHosts) SelectHost(name string) string {
	ip := "None"
	for _, hosts := range vmshosts.VmsHosts {
		if hosts.Name == name {
			ip = hosts.Ip
		}
	}
	return ip
}

// GetHosts 读取主机信息
func GetHosts(client *vim25.Client, vmshosts *VmsHosts) {
	m := view.NewManager(client)
	v, err := m.CreateContainerView(ctx, client.ServiceContent.RootFolder, []string{"HostSystem"}, true)
	if err != nil {
		panic(err)
	}
	defer v.Destroy(ctx)
	var hss []mo.HostSystem
	err = v.Retrieve(ctx, []string{"HostSystem"}, []string{"summary"}, &hss)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("主机名:\t%s\n", hss[0].Summary.Host.Value)
	// fmt.Printf("IP:\t%s\n", hss[0].Summary.Config.Name)
	for _, hs := range hss {
		vmshosts.AddHost(hs.Summary.Host.Value, hs.Summary.Config.Name)
	}
}

// GetVms获取所有vm信息
func GetVms(client *vim25.Client, vmshosts *VmsHosts) {
	m := view.NewManager(client)
	v, err := m.CreateContainerView(ctx, client.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		panic(err)
	}
	defer v.Destroy(ctx)
	var vms []mo.VirtualMachine
	err = v.Retrieve(ctx, []string{"VirtualMachine"}, []string{"summary", "runtime", "datastore"}, &vms)
	if err != nil {
		panic(err)
	}
	// 输出虚拟机信息到csv
	file, _ := os.OpenFile("./vms.csv", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	//防止中文乱码
	file.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(file)
	w.Write([]string{"宿主机", "虚拟机", "系统", "状态", "IP地址", "资源"})
	w.Flush()
	for _, vm := range vms {
		//虚拟机资源信息
		res := strconv.Itoa(int(vm.Summary.Config.MemorySizeMB)) + " MB " + strconv.Itoa(int(vm.Summary.Config.NumCpu)) + " vCPU(s) " + units.ByteSize(vm.Summary.Storage.Committed+vm.Summary.Storage.Uncommitted).String()
		w.Write([]string{vmshosts.SelectHost(vm.Summary.Runtime.Host.Value), vm.Summary.Config.Name, vm.Summary.Config.GuestFullName, string(vm.Summary.Runtime.PowerState), vm.Summary.Guest.IpAddress, res})
		w.Flush()
	}
	file.Close()

	// 批量插入到数据库
	var modelVms []*Vm
	for _, vm := range vms {
		modelVms = append(modelVms, &Vm{
			Uuid:       vm.Summary.Config.Uuid,
			Vc:         VSPHERE_IP,
			Esxi:       vm.Summary.Runtime.Host.Value,
			Name:       vm.Summary.Config.Name,
			Ip:         vm.Summary.Guest.IpAddress,
			PowerState: string(vm.Summary.Runtime.PowerState),
		})
	}
}

func main() {
	var vmhosts *VmsHosts
	GetHosts(client, vmhosts)
}
