package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/vmware/govmomi/performance"
	"github.com/vmware/govmomi/session/cache"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

func vlogin(ctx context.Context, vc, user, pwd string) (*vim25.Client, error) {
	//
	// Create a vSphere/vCenter client
	//
	//    The govmomi client requires a URL object, u, not just a string representation of the vCenter URL.
	// govmomi客户端需要一个URL对象，u，而不仅仅是vCenter URL的一个字符串表示。
	u, err := soap.ParseURL(vc)
	if u == nil {
		fmt.Printf("could not parse URL (environment variables set?)")
	}
	if err != nil {
		fmt.Printf("URL parsing not successful, error %v", err)
		return nil, err
	}
	u.User = url.UserPassword(user, pwd)
	//
	// Ripped from https://github.com/vmware/govmomi/blob/master/examples/examples.go
	//
	// Share session cache
	// 分享会话缓存
	s := &cache.Session{
		URL:      u,
		Insecure: true,
	}

	c := new(vim25.Client)

	err = s.Login(ctx, c, nil)

	if err != nil {
		fmt.Printf("Log in not successful- could not get vCenter client: %v", err)
		return nil, err
	} else {
		fmt.Printf("Log in successful")

		return c, nil
	}
}

func main() {

	vc := os.Getenv("GOVMOMI_URL")
	user := os.Getenv("GOVMOMI_USERNAME")
	pwd := os.Getenv("GOVMOMI_PASSWORD")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, _ := vlogin(ctx, vc, user, pwd)

	fmt.Println("Login successful =================")
	//
	// Call the login function
	//

	m := view.NewManager(client)
	kind := []string{"ClusterComputeResource"}
	v, _ := m.CreateContainerView(ctx, client.ServiceContent.RootFolder, kind, true)

	var clusters []mo.ClusterComputeResource
	err1 := v.Retrieve(ctx, kind, []string{"summary", "name", "network", "datastore", "host"}, &clusters)
	if err1 != nil {
		fmt.Println("err, cluster:", err1)
	}
	for _, cluster := range clusters {
		fmt.Println(cluster.Configuration.DasConfig.VmMonitoring)
		fmt.Println(len(cluster.Recommendation))

		break
	}

	perfManager := performance.NewManager(client)
	counters, _ := perfManager.CounterInfoByName(ctx)
	//进行对其周期和最大的获取的值进行配置
	spec := types.PerfQuerySpec{
		MaxSample:  1,
		MetricId:   []types.PerfMetricId{{Instance: ""}},
		IntervalId: 20,
	}
	//获取想要字段的id在总字段对应的位置和数据
	sample, _ := perfManager.SampleByName(ctx, spec, names, hostList)
	//获取想要的数据全部信息
	res, _ := perfManager.ToMetricSeries(ctx, sample)

	fmt.Println("-------------------------------- \n")

	m2 := view.NewManager(client)
	kind2 := []string{"HostSystem"}
	v2, _ := m2.CreateContainerView(ctx, client.ServiceContent.RootFolder, kind2, true)
	var hosts []mo.HostSystem
	err2 := v2.Retrieve(ctx, kind, []string{"summary", "datastore"}, &hosts)
	if err2 != nil {
		fmt.Println("error getting")
	}
	for _, host := range hosts {
		fmt.Println(host.Name)
	}

	perfManager := performance.NewManager(client)
	counters, _ := perfManager.CounterInfoByName(ctx)
	fmt.Println("counters: ", counters)
	// spec := types.PerfQuerySpec{
	// 	MaxSample:  1,
	// 	MetricId:   []types.PerfMetricId{{Instance: ""}},
	// 	IntervalId: 20,
	// }
	// sample, err := perfManager.SampleByName(ctx, spec, names, hostList)
	// res, _ := perfManager.ToMetricSeries(ctx, sample)

}
