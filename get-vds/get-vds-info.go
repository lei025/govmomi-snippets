////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// Description:		Go code to connect to vSphere via environment
// 			variables and retrieve the VDS and VDS PortGroup Information
// 			并检索VDS和VDS端口组信息。
//
// Author:		   	Cormac J. Hogan (VMware)
//
// Date:			04 Jul 2021
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//

package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/vmware/govmomi/session/cache"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

func main() {
	//
	// 3 environment variables are required in order to connect to the vSphere infra
	//
	// Set these in your shell to reflect your vSphere infra:
	//
	// GOVMOMI_URL
	// GOVMOMI_USERNAME
	// GOVMOMI_PASSWORD
	//

	vc := os.Getenv("GOVMOMI_URL")

	if len(vc) > 0 {
		fmt.Printf("DEBUG: vc is %s\n", vc)
	} else {
		fmt.Printf("Unable to find env var GOVMOMI_URL, has it been set?\n", vc)
		return
	}

	user := os.Getenv("GOVMOMI_USERNAME")
	if len(user) > 0 {
		fmt.Printf("DEBUG: user is %s\n", user)
	} else {
		fmt.Printf("Unable to find env var GOVMOMI_USERNAME, has it been set?\n", vc)
		return
	}
	pwd := os.Getenv("GOVMOMI_PASSWORD")

	if len(pwd) > 0 {
		fmt.Printf("DEBUG: password is %s\n", pwd)
	} else {
		fmt.Printf("Unable to find env GOVMOMI_PASSWORD, has it been set?\n", vc)
		return
	}

	//
	// This section allows for insecure vSphere logins
	// 此部分允许不安全的vSphere登录

	var insecure bool
	flag.BoolVar(&insecure, "insecure", true, "ignore any vCenter TLS cert validation error")

	//
	// Imagine that there were multiple operations taking place such as processing some data, logging into vCenter, etc.
	// If one of the operations failed, the context would be used to share the fact that all of the other operations
	// sharing that context needs cancelling.
	// 想象一下，有多个操作正在进行，如处理一些数据、登录到vCenter等。
	// 如果其中一个操作失败了，该上下文将被用来分享所有其他的操作
	// 共享该上下文需要取消。

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//
	// Create a vSphere/vCenter client
	//
	//    The client requires a URL object not just a string representation of the vCenter URL
	//	客户端需要一个URL对象，而不仅仅是vCenter URL的字符串表示。

	u, err := soap.ParseURL(vc)

	if u == nil {
		fmt.Println("could not parse URL (environment variables set?)")
	}

	if err != nil {
		fmt.Println("URL parsing not successful, error ", err)
		return
	}

	u.User = url.UserPassword(user, pwd)

	//
	// Share session cache
	//
	s := &cache.Session{
		URL:      u,
		Insecure: true,
	}

	//-------------------------------------------------------------------
	//
	//     vim25.Client - Call the function from the govmomi package
	//
	//     c, err - Return the client object c and an error object err
	//
	//     ctx - Pass in the shared context
	//
	// vim25.Client - 从govmomi包中调用函数
	//
	// c, err - 返回客户端对象 c 和错误对象 err
	//
	// ctx - 传递共享上下文
	//-------------------------------------------------------------------

	//
	//  A lot of GO functions return more than one variable/object
	//  The majority also return an object of type error.
	//
	//  If the function call is successful it returns nil in the place of an error object.
	//
	//  If something goes wrong the function should create a new error object with the appropriate messaging.
	//
	// 很多GO函数都会返回一个以上的变量/对象
	// 大多数还返回一个错误类型的对象。
	//
	// 如果函数调用成功，它将返回nil以代替错误对象。
	//
	// 如果出了问题，该函数应该创建一个新的错误对象，并给出相应的信息。

	c := new(vim25.Client)

	err = s.Login(ctx, c, nil)

	if err != nil {
		fmt.Println("")
		fmt.Println("Log in not successful (govmomi) - could not get vCenter client: ", err)
		fmt.Println("")
		return
	} else {
		fmt.Println("")
		fmt.Println("Log in successful (govmomi)")
		fmt.Println("")
	}

	// Create a view of DVS Network objects
	// 创建一个DVS网络对象的视图
	// DistributedVirtualSwitch 分布式虚拟交换机

	m := view.NewManager(c)

	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"DistributedVirtualSwitch"}, true)
	if err != nil {
		fmt.Println("Error : could not create DVS container view: ", err)
	}

	defer v.Destroy(ctx)

	// Retrieve summary property for all DVS
	// 检索所有DVS的摘要属性

	var vds []mo.DistributedVirtualSwitch
	err = v.Retrieve(ctx, []string{"DistributedVirtualSwitch"}, nil, &vds)
	if err != nil {
		fmt.Println("Error : could not retrieve DVS info: ", err)
	}

	fmt.Printf("\n")

	// Print details per DVS
	// Use 'govc object.collect network/DVS-Name' to see available fields to retrieve
	// 打印每个DVS的详细信息
	// 使用'govc object.collect network/DVS-Name'来查看可检索的字段

	for _, s := range vds {

		// gomvomi interface provides access to the underlying base type (VMwareDVSConfigInfo)
		// gomvomi接口提供对底层基本类型（VMwareDVSConfigInfo）的访问。

		config := s.Config.(*types.VMwareDVSConfigInfo)

		fmt.Printf("DVS Name is %s\n", config.Name)
		fmt.Printf("DVS Config Status is %s\n", s.ConfigStatus)
		fmt.Printf("DVS Overall Status is %s\n", s.OverallStatus)
		fmt.Printf("DVS Config Version is %s\n", config.ConfigVersion)
		fmt.Printf("DVS IP Address is %s\n", config.SwitchIpAddress)

		// gomvomi interface provides access to the underlying base type (VMwareDVSPortSetting)
		// gomvomi接口提供对底层基本类型（VMwareDVSPortSetting）的访问。

		portConfig := config.DefaultPortConfig.(*types.VMwareDVSPortSetting)

		// gomvomi interface provides access to the underlying base type (VmwareDistributedVirtualSwitchVlanIdSpec)
		// gomvomi接口提供对底层基本类型（VmwareDistributedVirtualSwitchVlanIdSpec）的访问。

		vlan := portConfig.Vlan.(*types.VmwareDistributedVirtualSwitchVlanIdSpec)

		// Display distributed switch vlan id,  if any 
		// 显示分布式交换机的vlan id，如果有的话
		fmt.Printf("vlan id type = %T\n", vlan.VlanId)
		fmt.Printf("vlan id = %v\n", vlan.VlanId)
		fmt.Printf("policy inheritable = %t\n", vlan.InheritablePolicy.Inherited)

		fmt.Printf("\n")
	}

	//
	// Turning our attention to the distributed port groups, create a view of DVS PG Network objects
	// 将我们的注意力转向分布式端口组，创建一个DVS PG网络对象的视图

	v1, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"DistributedVirtualPortgroup"}, true)
	if err != nil {
		fmt.Println("Error : could not create DVS PG container view: ", err)
	}

	defer v1.Destroy(ctx)

	// Retrieve summary property for all DVS
	// 检索所有DVS的摘要属性

	var vdspg []mo.DistributedVirtualPortgroup
	err = v1.Retrieve(ctx, []string{"DistributedVirtualPortgroup"}, nil, &vdspg)
	if err != nil {
		fmt.Println("Error : could not retrieve DVS PG info: ", err)
	}

	fmt.Printf("\n")

	// Print details per DVS-PG
	// Use 'govc object.collect /DC/network/DVPG-Name' to see available fields to retrieve
	// 打印每个 DVS-PG 的详细信息
	// 使用'govc object.collect /DC/network/DVPG-Name'来查看可用的字段来检索。

	for _, pg := range vdspg {
		fmt.Printf("Name of PG is %s\n", pg.Name)
		dpgPortConfig := pg.Config.DefaultPortConfig.(*types.VMwareDVSPortSetting)

		// Need the switch/case to avoid picking up uplinks which have trunk vlans and are a different type
		// 需要交换机/机箱来避免接收有主干网和不同类型的上行链路。

		switch dpgVlan := dpgPortConfig.Vlan.(type) {
		case *types.VmwareDistributedVirtualSwitchVlanIdSpec:
			fmt.Printf("Vlan id=%d\n", dpgVlan.VlanId)
		default:
			fmt.Printf("\tIgnoring VLAN %s of type %T\n", pg.Name, dpgVlan)
		}
	}
}
