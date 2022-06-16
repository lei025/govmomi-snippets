package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/vmware/govmomi/performance"
	"github.com/vmware/govmomi/view"

	"github.com/vmware/govmomi/session/cache"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/methods"
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

// var Perf_dict map[string]int32

// Perf_dict = make(map[string]string, 10)

func GetPerfDict(c *vim25.Client, ctx context.Context) (map[string]int32, error) {
	Perf_dict := make(map[string]int32)
	perfManager := performance.NewManager(c)
	counters, _ := perfManager.CounterInfo(ctx)
	for _, counter := range counters {
		counter_full := fmt.Sprintf("%s.%s.%s", counter.GroupInfo.GetElementDescription().Key,
			counter.NameInfo.GetElementDescription().Key,
			counter.RollupType)
		Perf_dict[counter_full] = counter.Key
	}
	return Perf_dict, nil
}

// func queryInfo(c *vim25.Client, vchtime, counterId int, instance, entity, interval){
// 	perfManager := performance.NewManager(c)
// 	vim25

// }

func main() {

	vc := os.Getenv("GOVMOMI_URL")
	user := os.Getenv("GOVMOMI_USERNAME")
	pwd := os.Getenv("GOVMOMI_PASSWORD")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c, _ := vlogin(ctx, vc, user, pwd)

	fmt.Println("Login successful =================")
	//
	// Call the login function
	now, _ := methods.GetCurrentTime(ctx, c) // vCenter server time (UTC)
	fmt.Println("Current time: ", now)
	// content := c.ServiceContent
	// fmt.Println("content: ", content)

	dict, _ := GetPerfDict(c, ctx)
	// fmt.Println("Perf_dict", Perf_dict)

	// perfManager := performance.NewManager(c)
	// spec := types.PerfQuerySpec{
	// 	MaxSample:  1,
	// 	MetricId:   []types.PerfMetricId{{Instance: "*"}, {CounterId: 632}},
	// 	IntervalId: int32(*interval),
	// }
	m := view.NewManager(c)
	kind := []string{"HostSystem"}
	v, _ := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, kind, true)

	vmsRefs, err := v.Find(ctx, []string{"VirtualMachine"}, nil)
	fmt.Println("vmsRefs ", vmsRefs, "len ", len(vmsRefs))

	var hosts []mo.HostSystem
	err = v.Retrieve(ctx, kind, []string{"summary", "datastore"}, &hosts)
	if err != nil {
		return
	}
	perfManager := performance.NewManager(c)
	counters, err := perfManager.CounterInfoByName(ctx)
	// fmt.Println("counters:", counters)
	fmt.Printf("datastore.read.average,", *counters["datastore.read.average"])

	spec := types.PerfQuerySpec{
		MaxSample:  1,
		MetricId:   []types.PerfMetricId{{Instance: ""}},
		IntervalId: 20,
	}
	// names := []string{"datastore.read.average"}
	var names []string
	for name := range counters {
		names = append(names, name)
	}
	sample, err := perfManager.SampleByName(ctx, spec, names, vmsRefs)
	res, _ := perfManager.ToMetricSeries(ctx, sample)
	fmt.Println("res: ", res)

}

/*
statDatastoreRead = queryInfo(content, vchtime, perf_dict['datastore.read.average'],"*", vm_instance, interval)
statDatastoreRead = queryInfo(content, vchtime, 180,"*", vm_instance, 60)
queryInfo(content, vchtime, counterId, instance, entity, interval)

def queryInfo(content, vchtime, counterId, instance, entity, interval):
    perfManager = content.perfManager
    metricId = vim.PerformanceManager.MetricId(counterId=counterId, instance=instance)
    startTime = vchtime - timedelta(seconds=(interval + 60))
    endTime = vchtime - timedelta(seconds=60)
    query = vim.PerformanceManager.QuerySpec(intervalId=20, entity=entity, metricId=[metricId], startTime=startTime,
                                             endTime=endTime)
    perfResults = perfManager.QueryPerf(querySpec=[query])
    if perfResults:
        return perfResults
    else:
        return False
*/
