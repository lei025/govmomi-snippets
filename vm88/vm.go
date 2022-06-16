package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/performance"
	"github.com/vmware/govmomi/units"
	"github.com/vmware/govmomi/vapi/library"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vapi/vcenter"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

type VirtualMachines struct {
	Name   string
	System string
	Self   Self
	VM     types.ManagedObjectReference
}

type TemplateInfo struct {
	Name   string
	System string
	Self   Self
	VM     types.ManagedObjectReference
}

type DatastoreSummary struct {
	Datastore          Datastore `json:"Datastore"`
	Name               string    `json:"Name"`
	URL                string    `json:"Url"`
	Capacity           int64     `json:"Capacity"`
	FreeSpace          int64     `json:"FreeSpace"`
	Uncommitted        int64     `json:"Uncommitted"`
	Accessible         bool      `json:"Accessible"`
	MultipleHostAccess bool      `json:"MultipleHostAccess"`
	Type               string    `json:"Type"`
	MaintenanceMode    string    `json:"MaintenanceMode"`
	DatastoreSelf      types.ManagedObjectReference
}

type Datastore struct {
	Type  string `json:"Type"`
	Value string `json:"Value"`
}

type HostSummary struct {
	Host        Host   `json:"Host"`
	Name        string `json:"Name"`
	UsedCPU     int64  `json:"UsedCPU"`
	TotalCPU    int64  `json:"TotalCPU"`
	FreeCPU     int64  `json:"FreeCPU"`
	UsedMemory  int64  `json:"UsedMemory"`
	TotalMemory int64  `json:"TotalMemory"`
	FreeMemory  int64  `json:"FreeMemory"`
	Self        types.ManagedObjectReference
}

type Host struct {
	Type  string `json:"Type"`
	Value string `json:"Value"`
}

type HostVM struct {
	Host map[string][]VMS
}

type VMS struct {
	Name  string
	Value string
}

type DataCenter struct {
	Datacenter      Self
	Name            string
	VmFolder        Self
	HostFolder      Self
	DatastoreFolder Self
}

type ClusterInfo struct {
	Cluster      Self
	Name         string
	Parent       Self
	ResourcePool Self
	Hosts        []types.ManagedObjectReference
	Datastore    []types.ManagedObjectReference
}

type ResourcePoolInfo struct {
	ResourcePool     Self
	Name             string
	Parent           Self
	ResourcePoolList []types.ManagedObjectReference
	Resource         types.ManagedObjectReference
}

type FolderInfo struct {
	Folder      Self
	Name        string
	ChildEntity []types.ManagedObjectReference
	Parent      Self
	FolderSelf  types.ManagedObjectReference
}

type Self struct {
	Type  string
	Value string
}

type CreateMap struct {
	TempName    string
	Datacenter  string
	Cluster     string
	Host        string
	Resources   string
	Storage     string
	VmName      string
	SysHostName string
	Network     string
}

type VmWare struct {
	IP     string
	User   string
	Pwd    string
	client *govmomi.Client
	ctx    context.Context
}

func NewVmWare(IP, User, Pwd string) *VmWare {
	u := &url.URL{
		Scheme: "https",
		Host:   IP,
		Path:   "/sdk",
	}
	ctx := context.Background()
	u.User = url.UserPassword(User, Pwd)
	client, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		panic(err)
	}
	return &VmWare{
		IP:     IP,
		User:   User,
		Pwd:    Pwd,
		client: client,
		ctx:    ctx,
	}
}

func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func (vw *VmWare) getBase(tp string) (v *view.ContainerView, error error) {
	m := view.NewManager(vw.client.Client)

	v, err := m.CreateContainerView(vw.ctx, vw.client.Client.ServiceContent.RootFolder, []string{tp}, true)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// var Perf_dict map[string]interface{}

// Perf_dict = make(map[string]string, 10)

func (vw *VmWare) GetPerfDict() (map[string]int32, error) {
	Perf_dict := make(map[string]int32)
	perfManager := performance.NewManager(vw.client.Client)
	counters, _ := perfManager.CounterInfo(vw.ctx)
	for _, counter := range counters {
		counter_full := fmt.Sprintf("%s.%s.%s", counter.GroupInfo.GetElementDescription().Key,
			counter.NameInfo.GetElementDescription().Key,
			counter.RollupType)
		Perf_dict[counter_full] = counter.Key
	}
	return Perf_dict, nil
}

func (vw *VmWare) GetAllVmClient1() (vmList []mo.VirtualMachine, err error) {
	v, err := vw.getBase("VirtualMachine")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.ctx)
	// var vms []mo.VirtualMachine
	err = v.Retrieve(vw.ctx, []string{"VirtualMachine"}, []string{"summary"}, &vmList)
	if err != nil {
		return nil, err
	}
	return vmList, nil
}

func (vw *VmWare) GetAllVmClient() (vmList []VirtualMachines, templateList []TemplateInfo, err error) {
	v, err := vw.getBase("VirtualMachine")
	if err != nil {
		return nil, nil, err
	}
	defer v.Destroy(vw.ctx)
	var vms []mo.VirtualMachine
	err = v.Retrieve(vw.ctx, []string{"VirtualMachine"}, []string{"summary"}, &vms)
	if err != nil {
		return nil, nil, err
	}
	for _, vm := range vms {
		//if vm.Summary.Config.Name == "测试机器" {
		//	v := object.NewVirtualMachine(vw.client.Client, vm.Self)
		//	vw.setIP(v)
		//}
		if vm.Summary.Config.Template {
			templateList = append(templateList, TemplateInfo{
				Name:   vm.Summary.Config.Name,
				System: vm.Summary.Config.GuestFullName,
				Self: Self{
					Type:  vm.Self.Type,
					Value: vm.Self.Value,
				},
				VM: vm.Self,
			})
		} else {
			vmList = append(vmList, VirtualMachines{
				Name:   vm.Summary.Config.Name,
				System: vm.Summary.Config.GuestFullName,
				Self: Self{
					Type:  vm.Self.Type,
					Value: vm.Self.Value,
				},
				VM: vm.Self,
			})
		}
	}

	// var interval = flag.Int("i", 20, "Interval ID")
	// for _, vm := range vmList {
	// 	vmmm, _ := queryInfo(vw.ctx, vw.client.Client, interval, []string{"net.received.average"}, 148, types.ManagedObjectReference(vm.Self))
	// 	fmt.Println("vmmm", vmmm)
	// }

	return vmList, templateList, nil
}

func (vw *VmWare) GetAllHost() (hostList []*HostSummary, err error) {
	v, err := vw.getBase("HostSystem")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.ctx)
	var hss []mo.HostSystem
	err = v.Retrieve(vw.ctx, []string{"HostSystem"}, []string{"summary"}, &hss)
	if err != nil {
		return nil, err
	}
	/*
			    hostInfo = {
		                "hostname" : summary.config.name,
		                "uptime"   : formatTime(stats.uptime),
		                "uuid"   : dc_name + "_" + host.name,
		                "vendor"   : summary.hardware.vendor,
		                "parentuuid": vcenter_ip,
		                "datacenter": dc_name,
		                "status": summary.runtime.powerState
		                }
	*/
	for _, hs := range hss {
		fmt.Println("hostInfo: ",
			"hostname: ", hs.Summary.Config.Name,
			"uptime: ", hs.Summary.QuickStats.Uptime,
			"uuid:", "dc_name_", hs.Name,
			"vendor: ", hs.Summary.Hardware.Vendor,
			"parentuuid", "vcenter_ip",
			"datacenter: ", "dc_name",
			"status: ", hs.Summary.Runtime.PowerState)
		fmt.Println("--------------------------------")
		fmt.Println("cpuInfo: ",
			"coreNumber: ", int64(hs.Summary.Hardware.NumCpuCores),
			"threadNumber: ", int64(hs.Summary.Hardware.NumCpuThreads),
			"modeName: ", hs.Summary.Hardware.CpuModel,
			"usedPercent",

		// ""

		)

		memInfo := map[string]interface{}{
			"total": hs.Summary.Hardware.MemorySize,
			"used":  int64(hs.Summary.QuickStats.OverallMemoryUsage * 1024 * 1024),
			"free":  hs.Summary.Hardware.MemorySize - int64(hs.Summary.QuickStats.OverallMemoryUsage*1024*1024),
		}
		fmt.Println("memInfo: ", memInfo)

		totalCPU := int64(hs.Summary.Hardware.CpuMhz) * int64(hs.Summary.Hardware.NumCpuCores)
		freeCPU := int64(totalCPU) - int64(hs.Summary.QuickStats.OverallCpuUsage)
		freeMemory := int64(hs.Summary.Hardware.MemorySize) - (int64(hs.Summary.QuickStats.OverallMemoryUsage) * 1024 * 1024)
		fmt.Println("freeMemory: ", freeMemory)
		fmt.Println(int64(hs.Summary.Hardware.MemorySize))
		fmt.Println((int64(hs.Summary.QuickStats.OverallMemoryUsage) * 1024 * 1024))
		fmt.Println(hs.Datastore)

		hostList = append(hostList, &HostSummary{
			Host: Host{
				Type:  hs.Summary.Host.Type,
				Value: hs.Summary.Host.Value,
			},
			Name:        hs.Summary.Config.Name,
			UsedCPU:     int64(hs.Summary.QuickStats.OverallCpuUsage),
			TotalCPU:    totalCPU,
			FreeCPU:     freeCPU,
			UsedMemory:  int64((units.ByteSize(hs.Summary.QuickStats.OverallMemoryUsage)) * 1024 * 1024),
			TotalMemory: int64(units.ByteSize(hs.Summary.Hardware.MemorySize)),
			FreeMemory:  freeMemory,
			Self:        hs.Self,
		})
	}

	dict, _ := vw.GetPerfDict()

	// var interval = flag.Int("i", 20, "Interval ID")
	var interval = 20
	summm, _ := queryInfo(vw.ctx, vw.client.Client, interval, []string{"net.received.average"}, dict["net.received.average"], types.ManagedObjectReference(hostList[0].Self))
	fmt.Println("summm", summm)

	return hostList, err
}

func (vw *VmWare) GetAllNetwork() (networkList []map[string]string, err error) {
	v, err := vw.getBase("Network")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.ctx)
	var networks []mo.Network
	err = v.Retrieve(vw.ctx, []string{"Network"}, nil, &networks)
	if err != nil {
		return nil, err
	}
	for _, net := range networks {
		networkList = append(networkList, map[string]string{
			"Vlan":      net.Name,
			"NetworkID": strings.Split(net.Reference().String(), ":")[1],
		})
	}
	return networkList, nil
}

func (vw *VmWare) GetAllDatastore() (datastoreList []DatastoreSummary, err error) {
	v, err := vw.getBase("Datastore")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.ctx)
	var dss []mo.Datastore
	err = v.Retrieve(vw.ctx, []string{"Datastore"}, []string{"summary"}, &dss)
	if err != nil {
		return nil, err
	}
	for _, ds := range dss {
		datastoreList = append(datastoreList, DatastoreSummary{
			Name: ds.Summary.Name,
			Datastore: Datastore{
				Type:  ds.Summary.Datastore.Type,
				Value: ds.Summary.Datastore.Value,
			},
			Type:          ds.Summary.Type,
			Capacity:      int64(units.ByteSize(ds.Summary.Capacity)),
			FreeSpace:     int64(units.ByteSize(ds.Summary.FreeSpace)),
			DatastoreSelf: ds.Self,
		})
	}
	return
}

func (vw *VmWare) GetHostVm() (hostVm map[string][]VMS, err error) {
	hostList, err := vw.GetAllHost() //
	if err != nil {
		return
	}
	var hostIDList []string
	hostVm = make(map[string][]VMS)
	for _, host := range hostList {
		hostIDList = append(hostIDList, host.Host.Value)
		hostVm[host.Self.Value] = []VMS{}
	}
	v, err := vw.getBase("VirtualMachine")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.ctx)
	var vms []mo.VirtualMachine
	err = v.Retrieve(vw.ctx, []string{"VirtualMachine"}, []string{"summary"}, &vms)
	if err != nil {
		return nil, err
	}
	for _, vm := range vms {
		if IsContain(hostIDList, vm.Summary.Runtime.Host.Value) {
			hostVm[vm.Summary.Runtime.Host.Value] = append(hostVm[vm.Summary.Runtime.Host.Value], VMS{
				Name:  vm.Summary.Config.Name,
				Value: vm.Summary.Vm.Value,
			})
		}
		//s, _ := json.Marshal(vm.Summary)
		//fmt.Println(string(s))
	}
	//fmt.Println(hostVm)
	return
}

func (vw *VmWare) GetAllCluster() (clusterList []ClusterInfo, err error) {
	v, err := vw.getBase("ClusterComputeResource")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.ctx)
	var crs []mo.ClusterComputeResource
	err = v.Retrieve(vw.ctx, []string{"ClusterComputeResource"}, []string{}, &crs)
	if err != nil {
		return nil, err
	}
	for _, cr := range crs {
		clusterList = append(clusterList, ClusterInfo{
			Cluster: Self{
				Type:  cr.Self.Type,
				Value: cr.Self.Value,
			},
			Name: cr.Name,
			Parent: Self{
				Type:  cr.Parent.Type,
				Value: cr.Parent.Value,
			},
			ResourcePool: Self{
				Type:  cr.ResourcePool.Type,
				Value: cr.ResourcePool.Value,
			},
			Hosts:     cr.Host,
			Datastore: cr.Datastore,
		})
	}
	fmt.Println(clusterList)
	return
}

func (vw *VmWare) GetAllDatacenter() (dataCenterList []DataCenter, err error) {
	v, err := vw.getBase("Datacenter")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.ctx)
	var dcs []mo.Datacenter
	err = v.Retrieve(vw.ctx, []string{"Datacenter"}, []string{}, &dcs)
	if err != nil {
		return nil, err
	}
	for _, dc := range dcs {
		dataCenterList = append(dataCenterList, DataCenter{
			Datacenter: Self{
				Type:  dc.Self.Type,
				Value: dc.Self.Value,
			},
			Name: dc.Name,
			VmFolder: Self{
				Type:  dc.VmFolder.Type,
				Value: dc.VmFolder.Value,
			},
			HostFolder: Self{
				Type:  dc.HostFolder.Type,
				Value: dc.HostFolder.Value,
			},
			DatastoreFolder: Self{
				Type:  dc.DatastoreFolder.Type,
				Value: dc.DatastoreFolder.Value,
			},
		})
		fmt.Println("dc_name: ", dc.Name)
	}

	// {{Datacenter datacenter-7}
	// Datacenter
	// {Folder group-v8}
	// {Folder group-h9}
	// {Folder group-s10}}

	fmt.Println(dataCenterList)
	return
}

func (vw *VmWare) GetAllResourcePool() (resourceList []ResourcePoolInfo, err error) {
	v, err := vw.getBase("ResourcePool")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.ctx)
	var rps []mo.ResourcePool
	err = v.Retrieve(vw.ctx, []string{"ResourcePool"}, []string{}, &rps)
	for _, rp := range rps {
		//if rp.Name == "测试虚机" {
		//	s, _ := json.Marshal(rp)
		//	fmt.Println(string(s))
		//}
		resourceList = append(resourceList, ResourcePoolInfo{
			ResourcePool: Self{
				Type:  rp.Self.Type,
				Value: rp.Self.Value,
			},
			Name: rp.Name,
			Parent: Self{
				Type:  rp.Parent.Type,
				Value: rp.Parent.Value,
			},
			ResourcePoolList: rp.ResourcePool,
			Resource:         rp.Self,
		})
	}
	return
}

func (vw *VmWare) GetFolder() (folderList []FolderInfo, err error) {
	v, err := vw.getBase("Folder")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.ctx)
	var folders []mo.Folder
	err = v.Retrieve(vw.ctx, []string{"Folder"}, []string{}, &folders)
	for _, folder := range folders {
		//newFolder := object.NewFolder(vw.client.Client, folder.Self)
		//fmt.Println(newFolder)
		folderList = append(folderList, FolderInfo{
			Folder: Self{
				Type:  folder.Self.Type,
				Value: folder.Self.Value,
			},
			Name:        folder.Name,
			ChildEntity: folder.ChildEntity,
			Parent: Self{
				Type:  folder.Parent.Type,
				Value: folder.Parent.Value,
			},
			FolderSelf: folder.Self,
		})
		//break
	}
	return folderList, nil
}

func (vw *VmWare) getLibraryItem(ctx context.Context, rc *rest.Client) (*library.Item, error) {
	const (
		libraryName     = "模板"
		libraryItemName = "template-rehl7.7"
		libraryItemType = "ovf"
	)

	m := library.NewManager(rc)
	libraries, err := m.FindLibrary(ctx, library.Find{Name: libraryName})
	if err != nil {
		fmt.Printf("Find library by name %s failed, %v", libraryName, err)
		return nil, err
	}

	if len(libraries) == 0 {
		fmt.Printf("Library %s was not found", libraryName)
		return nil, fmt.Errorf("library %s was not found", libraryName)
	}

	if len(libraries) > 1 {
		fmt.Printf("There are multiple libraries with the name %s", libraryName)
		return nil, fmt.Errorf("there are multiple libraries with the name %s", libraryName)
	}

	items, err := m.FindLibraryItems(ctx, library.FindItem{Name: libraryItemName,
		Type: libraryItemType, LibraryID: libraries[0]})

	if err != nil {
		fmt.Printf("Find library item by name %s failed", libraryItemName)
		return nil, fmt.Errorf("find library item by name %s failed", libraryItemName)
	}

	if len(items) == 0 {
		fmt.Printf("Library item %s was not found", libraryItemName)
		return nil, fmt.Errorf("library item %s was not found", libraryItemName)
	}

	if len(items) > 1 {
		fmt.Printf("There are multiple library items with the name %s", libraryItemName)
		return nil, fmt.Errorf("there are multiple library items with the name %s", libraryItemName)
	}

	item, err := m.GetLibraryItem(ctx, items[0])
	if err != nil {
		fmt.Printf("Get library item by %s failed, %v", items[0], err)
		return nil, err
	}
	return item, nil
}

func (vw *VmWare) CreateVM() {
	createData := CreateMap{
		TempName:    "xxx",
		Datacenter:  "xxx",
		Cluster:     "xxx",
		Host:        "xxx",
		Resources:   "xxx",
		Storage:     "xxx",
		VmName:      "xxx",
		SysHostName: "xxx",
		Network:     "xxx",
	}
	_, templateList, err := vw.GetAllVmClient()
	if err != nil {
		panic(err)
	}
	var templateNameList []string
	for _, template := range templateList {
		templateNameList = append(templateNameList, template.Name)
	}
	if !IsContain(templateNameList, createData.TempName) {
		fmt.Fprintf(os.Stderr, "模版不存在，虚拟机创建失败")
		return
	}
	resourceList, err := vw.GetAllResourcePool()
	if err != nil {
		panic(err)
	}
	var resourceStr, resourceID string
	for _, resource := range resourceList {
		if resource.Name == createData.Resources {
			resourceStr = resource.Name
			resourceID = resource.ResourcePool.Value
		}
	}
	if resourceStr == "" {
		fmt.Fprintf(os.Stderr, "资源池不存在，虚拟机创建失败")
		return
	}
	fmt.Println("ResourceID", resourceID)
	datastoreList, err := vw.GetAllDatastore()
	if err != nil {
		panic(err)
	}
	var datastoreID, datastoreStr string
	for _, datastore := range datastoreList {
		if datastore.Name == createData.Storage {
			datastoreID = datastore.Datastore.Value
			datastoreStr = datastore.Name
		}
	}
	if datastoreStr == "" {
		fmt.Fprintf(os.Stderr, "存储中心不存在，虚拟机创建失败")
		return
	}
	fmt.Println("DatastoreID", datastoreID)
	networkList, err := vw.GetAllNetwork()
	if err != nil {
		panic(err)
	}
	var networkID, networkStr string
	for _, network := range networkList {
		if network["Vlan"] == createData.Network {
			networkStr = network["Vlan"]
			networkID = network["NetworkID"]
		}
	}

	if networkStr == "" {
		fmt.Fprintf(os.Stderr, "网络不存在，虚拟机创建失败")
		return
	}
	fmt.Println("NetworkID", networkID)
	finder := find.NewFinder(vw.client.Client)
	//resourcePools, err := finder.DatacenterList(vw.ctx, "*")
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Failed to list resource pool at vc %v", err)
	//	os.Exit(1)
	//}
	//fmt.Println(reflect.TypeOf(resourcePools[0].Reference().Value), resourcePools)
	folders, err := finder.FolderList(vw.ctx, "*")
	var folderID string
	for _, folder := range folders {
		if folder.InventoryPath == "/"+createData.Datacenter+"/vm" {
			folderID = folder.Reference().Value
		}
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list folder at vc  %v", err)
		return
	}
	rc := rest.NewClient(vw.client.Client)
	if err := rc.Login(vw.ctx, url.UserPassword(vw.User, vw.Pwd)); err != nil {
		fmt.Fprintf(os.Stderr, "rc Login filed, %v", err)
		return
	}
	item, err := vw.getLibraryItem(vw.ctx, rc)
	if err != nil {
		panic(err)
	}
	//cloneSpec := &types.VirtualMachineCloneSpec{
	//	PowerOn:  false,
	//	Template: cmd.template,
	//}
	// 7fa9e782-cba2-4061-95fc-4ebb08ec127a
	fmt.Println("Item", item.ID)
	m := vcenter.NewManager(rc)
	fr := vcenter.FilterRequest{
		Target: vcenter.Target{
			ResourcePoolID: resourceID,
			FolderID:       folderID,
		},
	}
	r, err := m.FilterLibraryItem(vw.ctx, item.ID, fr)
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
	fmt.Println(11111111111, r.Networks, r.StorageGroups)
	networkKey := r.Networks[0]
	//storageKey := r.StorageGroups[0]
	deploy := vcenter.Deploy{
		DeploymentSpec: vcenter.DeploymentSpec{
			Name:               createData.VmName,
			DefaultDatastoreID: datastoreID,
			AcceptAllEULA:      true,
			NetworkMappings: []vcenter.NetworkMapping{
				{
					Key:   networkKey,
					Value: networkID,
				},
			},
			StorageMappings: []vcenter.StorageMapping{{
				Key: "",
				Value: vcenter.StorageGroupMapping{
					Type:         "DATASTORE",
					DatastoreID:  datastoreID,
					Provisioning: "thin",
				},
			}},
			StorageProvisioning: "thin",
		},
		Target: vcenter.Target{
			ResourcePoolID: resourceID,
			FolderID:       folderID,
		},
	}
	ref, err := vcenter.NewManager(rc).DeployLibraryItem(vw.ctx, item.ID, deploy)
	if err != nil {
		fmt.Println(4444444444, err)
		panic(err)
	}
	f := find.NewFinder(vw.client.Client)
	obj, err := f.ObjectReference(vw.ctx, *ref)
	if err != nil {
		panic(err)
	}
	_ = obj.(*object.VirtualMachine)

	//datastores, err := finder.VirtualMachineList(vw.ctx, "*/group-v629")
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Failed to list datastore at vc %v", err)
	//	os.Exit(1)
	//}
	//fmt.Println(datastores)
}

func (vw *VmWare) CloneVM() {
	cloneData := CreateMap{
		TempName:    "xxx",
		Datacenter:  "xxx",
		Cluster:     "xxx",
		Host:        "xxx",
		Resources:   "xxx",
		Storage:     "xxx",
		VmName:      "xxx",
		SysHostName: "xxx",
		Network:     "xxx",
	}
	vmList, templateList, err := vw.GetAllVmClient()
	if err != nil {
		panic(err)
	}
	var templateNameList []string
	var vmTemplate types.ManagedObjectReference
	for _, template := range templateList {
		templateNameList = append(templateNameList, template.Name)
		if template.Name == cloneData.TempName {
			vmTemplate = template.VM
		}
	}
	if !IsContain(templateNameList, cloneData.TempName) {
		fmt.Fprintf(os.Stderr, "模版不存在，虚拟机克隆失败")
		return
	}
	dataCenterList, err := vw.GetAllDatacenter()
	if err != nil {
		panic(err)
	}
	var datacenterID, datacenterName string
	for _, datacenter := range dataCenterList {
		if datacenter.Name == cloneData.Datacenter {
			datacenterID = datacenter.Datacenter.Value
			datacenterName = datacenter.Name
		}
	}
	if datacenterName == "" {
		fmt.Fprintf(os.Stderr, "数据中心不存在，虚拟机克隆失败")
		return
	}
	hostList, err := vw.GetAllHost()
	if err != nil {
		panic(err)
	}
	var hostName string
	var hostRef types.ManagedObjectReference
	for _, host := range hostList {
		if host.Name == cloneData.Host {
			hostName = host.Name
			hostRef = host.Self
		}
	}
	if hostName == "" {
		fmt.Fprintf(os.Stderr, "主机不存在，虚拟机克隆失败")
		return
	}
	resourceList, err := vw.GetAllResourcePool()
	if err != nil {
		panic(err)
	}
	var resourceStr, resourceID string
	var poolRef types.ManagedObjectReference
	for _, resource := range resourceList {
		if resource.Name == cloneData.Resources {
			resourceStr = resource.Name
			resourceID = resource.ResourcePool.Value
			poolRef = resource.Resource
		}
	}
	if resourceStr == "" {
		fmt.Fprintf(os.Stderr, "资源池不存在，虚拟机克隆失败")
		return
	}
	fmt.Println("ResourceID", resourceID)
	datastoreList, err := vw.GetAllDatastore()
	if err != nil {
		panic(err)
	}
	var datastoreID, datastoreStr string
	var datastoreRef types.ManagedObjectReference
	for _, datastore := range datastoreList {
		if datastore.Name == cloneData.Storage {
			datastoreID = datastore.Datastore.Value
			datastoreStr = datastore.Name
			datastoreRef = datastore.DatastoreSelf
		}
	}
	if datastoreStr == "" {
		fmt.Fprintf(os.Stderr, "存储中心不存在，虚拟机克隆失败")
		return
	}
	fmt.Println("DatastoreID", datastoreID)
	networkList, err := vw.GetAllNetwork()
	if err != nil {
		panic(err)
	}
	var networkID, networkStr string
	for _, network := range networkList {
		if network["Vlan"] == cloneData.Network {
			networkStr = network["Vlan"]
			networkID = network["NetworkID"]
		}
	}

	if networkStr == "" {
		fmt.Fprintf(os.Stderr, "网络不存在，虚拟机克隆失败")
		return
	}
	fmt.Println("NetworkID", networkID)
	clusterList, err := vw.GetAllCluster()
	if err != nil {
		panic(err)
	}
	var clusterID, clusterName string
	for _, cluster := range clusterList {
		if cluster.Name == cloneData.Cluster {
			clusterID = cluster.Cluster.Value
			clusterName = cluster.Name
		}
	}
	if clusterName == "" {
		fmt.Fprintf(os.Stderr, "集群不存在，虚拟机克隆失败")
		return
	}
	configSpecs := []types.BaseVirtualDeviceConfigSpec{}
	fmt.Println("ClusterID", clusterID)
	for _, vms := range vmList {
		if vms.Name == cloneData.VmName {
			fmt.Fprintf(os.Stderr, "虚机已存在，虚拟机克隆失败")
			return
		}
	}
	finder := find.NewFinder(vw.client.Client)
	folders, err := finder.FolderList(vw.ctx, "*")
	var Folder *object.Folder
	for _, folder := range folders {
		if folder.InventoryPath == "/"+cloneData.Datacenter+"/vm" {
			Folder = folder
		}
	}
	fmt.Println(Folder)
	folderList, err := vw.GetFolder()
	if err != nil {
		panic(err)
	}
	var folderRef types.ManagedObjectReference
	for _, folder := range folderList {
		if folder.Parent.Value == datacenterID && folder.Name == "vm" {
			folderRef = folder.FolderSelf
		}
	}
	fmt.Println("poolRef", poolRef)
	relocateSpec := types.VirtualMachineRelocateSpec{
		DeviceChange: configSpecs,
		Folder:       &folderRef,
		Pool:         &poolRef,
		Host:         &hostRef,
		Datastore:    &datastoreRef,
	}
	vmConf := &types.VirtualMachineConfigSpec{
		NumCPUs:  4,
		MemoryMB: 16 * 1024,
	}
	cloneSpec := &types.VirtualMachineCloneSpec{
		PowerOn:  false,
		Template: false,
		Location: relocateSpec,
		Config:   vmConf,
	}
	t := object.NewVirtualMachine(vw.client.Client, vmTemplate)
	newFolder := object.NewFolder(vw.client.Client, folderRef)
	fmt.Println(newFolder)
	fmt.Println(cloneData.VmName)
	fmt.Println(cloneSpec.Location)
	task, err := t.Clone(vw.ctx, newFolder, cloneData.VmName, *cloneSpec)
	if err != nil {
		panic(err)
	}
	fmt.Println("克隆任务开始，", task.Wait(vw.ctx))
}

func (vw *VmWare) setIP(vm *object.VirtualMachine) error {
	ipAddr := IpAddr{
		ip:       "192.168.80.108",
		netmask:  "255.255.255.0",
		gateway:  "192.168.80.254",
		hostname: "test",
	}
	cam := types.CustomizationAdapterMapping{
		Adapter: types.CustomizationIPSettings{
			Ip:         &types.CustomizationFixedIp{IpAddress: ipAddr.ip},
			SubnetMask: ipAddr.netmask,
			Gateway:    []string{ipAddr.gateway},
		},
	}
	customSpec := types.CustomizationSpec{
		NicSettingMap: []types.CustomizationAdapterMapping{cam},
		Identity:      &types.CustomizationLinuxPrep{HostName: &types.CustomizationFixedName{Name: ipAddr.hostname}},
	}
	task, err := vm.Customize(vw.ctx, customSpec)
	if err != nil {
		return err
	}
	return task.Wait(vw.ctx)
}

type IpAddr struct {
	ip       string
	netmask  string
	gateway  string
	hostname string
}

func (vw *VmWare) MigrateVM() {
	migrateData := "测试虚机"
	v, err := vw.getBase("VirtualMachine")
	if err != nil {
		panic(err)
	}
	defer v.Destroy(vw.ctx)
	var vms []mo.VirtualMachine
	err = v.Retrieve(vw.ctx, []string{"VirtualMachine"}, []string{"summary"}, &vms)
	if err != nil {
		panic(err)
	}
	var vmTarget types.ManagedObjectReference
	for _, vm := range vms {
		if vm.Summary.Config.Name == migrateData {
			vmTarget = vm.Self
		}
	}
	resourceList, err := vw.GetAllResourcePool()
	if err != nil {
		panic(err)
	}
	var resourceStr, resourceID string
	var poolRef types.ManagedObjectReference
	for _, resource := range resourceList {
		if resource.Name == "" {
			resourceStr = resource.Name
			resourceID = resource.ResourcePool.Value
			poolRef = resource.Resource
		}
	}
	if resourceStr == "" {
		fmt.Fprintf(os.Stderr, "资源池不存在，虚拟机迁移失败")
		return
	}
	fmt.Println("ResourceID", resourceID)
	hostList, err := vw.GetAllHost()
	if err != nil {
		panic(err)
	}
	var hostName string
	var hostRef types.ManagedObjectReference
	for _, host := range hostList {
		if host.Name == "xxxx" {
			hostName = host.Name
			hostRef = host.Self
		}
	}
	if hostName == "" {
		fmt.Fprintf(os.Stderr, "主机不存在，虚拟机迁移失败")
		return
	}
	t := object.NewVirtualMachine(vw.client.Client, vmTarget)
	pool := object.NewResourcePool(vw.client.Client, poolRef)
	host := object.NewHostSystem(vw.client.Client, hostRef)
	//var priority types.VirtualMachineMovePriority
	//var state types.VirtualMachinePowerState
	task, err := t.Migrate(vw.ctx, pool, host, "defaultPriority", "poweredOff")
	if err != nil {
		panic(err)
	}
	fmt.Println("虚拟机迁移中......")
	_ = task.Wait(vw.ctx)
	fmt.Println("虚拟机迁移完成.....")
}

func queryInfo(ctx context.Context, c *vim25.Client, interval int, counterName []string, counterId int32, entity types.ManagedObjectReference) (sum int64, err error) {
	// Create a PerfManager
	perfManager := performance.NewManager(c)
	if err != nil {
		return
	}
	// t := time.Now().Add(-600 * time.Second)
	// t1 := time.Now().Add(60 * time.Second)
	// Create PerfQuerySpec
	spec := types.PerfQuerySpec{
		Entity:     types.ManagedObjectReference{Type: entity.Type, Value: entity.Value},
		MaxSample:  1,
		MetricId:   []types.PerfMetricId{{Instance: "*", CounterId: counterId}},
		IntervalId: int32(interval),
		// StartTime:  &t,
		// EndTime:    &t1,
	}
	// Query metrics
	vmsRefss := []types.ManagedObjectReference{{Type: entity.Type, Value: entity.Value}}
	sample, err := perfManager.SampleByName(ctx, spec, counterName, vmsRefss)
	if err != nil {
		return
	}
	result, err := perfManager.ToMetricSeries(ctx, sample)
	if err != nil {
		return
	}
	// Read result
	for _, metric := range result {
		for _, v := range metric.Value {
			sum += v.Value[0]
		}
	}
	return sum, nil
}

func main() {
	vc := os.Getenv("GOVMOMI_URL")
	user := os.Getenv("GOVMOMI_USERNAME")
	pwd := os.Getenv("GOVMOMI_PASSWORD")
	vm := NewVmWare(vc, user, pwd)
	fmt.Println(vm)
	fmt.Println("-------------------------------- vmList")
	// vmList, _, _ := vm.GetAllVmClient()
	vmList, templateList, _ := vm.GetAllVmClient()
	for _, vm := range vmList {
		fmt.Println(vm)
	}
	fmt.Println("-------------------------------- templateList")
	for _, template := range templateList {
		fmt.Println(template)
	}
	fmt.Println("-------------------------------- hostList")
	hostList, _ := vm.GetAllHost()
	for _, hs := range hostList {
		fmt.Println(hs)
		fmt.Println(hs.Name)

	}
	/*
		fmt.Println("-------------------------------- datastore ")
		storeList, _ := vm.GetAllDatastore()
		for _, store := range storeList {
			fmt.Println(store)
		}

		fmt.Println("GetAllDatacenter --------------------------------")
		dataceterList, _ := vm.GetAllDatacenter()
		fmt.Println(len(dataceterList))
		for _, dataceter := range dataceterList {
			fmt.Println(dataceter)
			fmt.Println(dataceter.Name)
		}

			fmt.Println("GetAllFolder --------------------------------")
			folderList, _ := vm.GetFolder()
			fmt.Println(len(folderList))
			for _, folder := range folderList {
				fmt.Println(folder)
				fmt.Println("----")
			}

			fmt.Println("GetAllNetwork --------------------------------")
			networkList, _ := vm.GetAllNetwork()
			fmt.Println(len(networkList))
			for _, network := range networkList {
				fmt.Println(network)
			}
	*/

}
