////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// Description:		Go code to connect to vSphere via environment
// 					variables and retrieve the defautl datacenter
//
// Author:			Cormac J. Hogan (VMware)
//
// Date:			25 Jan 2021
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// client information from Doug MacEachern:
// ----------------------------------------
//
// govmomi.Client extends vim25.Client
// govmomi.Client does nothing extra aside from automatic login
//
// In the early days (2015), govmomi.Client did much more, but we moved most of it to vim25.Client.
// govmomi.Client remained for compatibility and minor convenience.
//
// Using soap.Client and vim25.Client directly allows apps to use other authentication methods,
// session caching, session keepalive, retries, fine grained TLS configuration, etc.
//
// For the inventory, ContainerView is a vSphere primitive.
// Compared to Finder, ContainerView tends to use less round trip calls to vCenter.
// It may generate more response data however.
//
// Finder was written for govc, where we treat the vSphere inventory as a virtual filesystem.
// The inventory path as input to `govc` behaves similar to the `ls` command, with support for relative paths,
// wildcard matching, etc.
//
// Use govc commands as a reference, and "godoc" for examples that can be run against `vcsim`:
// See: https://godoc.org/github.com/vmware/govmomi/view#pkg-examples
//

// ----------------------------------------
//
// govmomi.Client扩展了vim25.Client
// govmomi.Client除了自动登录外，不做任何额外的事情。
//
// 在早期（2015年），govmomi.Client做得更多，但我们把大部分移到vim25.Client。
// govmomi.Client保留了兼容性和小便利。
//
// 直接使用soap.Client和vim25.Client允许应用程序使用其他认证方法。
// 会话缓存、会话保持、重试、精细的TLS配置等。
//
// 对于清单来说，ContainerView是一个vSphere的原件。
// 与Finder相比，ContainerView倾向于减少对vCenter的往返调用。
// 但是，它可能会产生更多的响应数据。
//
// Finder是为govc编写的，我们将vSphere库存视为一个虚拟文件系统。
// 输入到 "govc "的清单路径与 "ls "命令的行为类似，支持相对路径。
// 通配符匹配等。

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/soap"
)

func main() {

	// We need to get 3 environment variables in order to connect to the vSphere infra
	//
	// Change these to reflect your vSphere infra:
	//
	// GOVMOMI_URL=vcsa-06.rainpole.com/sdk
	// GOVMOMI_USERNAME=administrator@vsphere.local
	// GOVMOMI_PASSWORD=VMware123!
	// GOVMOMI_INSECURE=true

	var insecure bool

	flag.BoolVar(&insecure, "insecure", true, "ignore any vCenter TLS cert validation error")

	vc := os.Getenv("GOVMOMI_URL")
	user := os.Getenv("GOVMOMI_USERNAME")
	pwd := os.Getenv("GOVMOMI_PASSWORD")

	fmt.Printf("DEBUG: vc is %s\n", vc)
	fmt.Printf("DEBUG: user is %s\n", user)
	fmt.Printf("DEBUG: password is %s\n", pwd)

	//
	// Imagine that there were multiple operations taking place such as processing some data, logging into vCenter, etc.
	// If one of the operations failed, the context would be used to share the fact that all of the other operations sharing that context needs cancelling.
	//

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//
	// Create a vSphere/vCenter client
	//
	//    The govmomi client requires a URL object not just a string representation of the vCenter URL.
	//    c, err - Return the client object c and an error object err
	//    govmomi.NewClient - Call the function from the govmomi package
	//    ctx - Pass in the shared context

	u, err := soap.ParseURL(vc)

	if u == nil {
		fmt.Println("could not parse URL (environment variables set?)")
	}

	if err != nil {
		fmt.Printf("URL parsing not successful, error %v", err)
		return
	}

	u.User = url.UserPassword(user, pwd)

	c, err := govmomi.NewClient(ctx, u, insecure)

	if err != nil {
		fmt.Printf("Log in not successful- could not get vCenter client: %v", err)
		return
	} else {
		fmt.Println("Log in successful")
		c.Logout(ctx)
	}
}
