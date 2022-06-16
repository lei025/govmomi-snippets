package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"

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

// var interval int32 = 60
var interval = flag.Int("i", 60, "Interval ID")

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
	//

	// Get virtual machines references
	m := view.NewManager(c)
	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, nil, true)
	if err != nil {
		fmt.Println("err ", err)
	}

	defer v.Destroy(ctx)

	// vmsRefs, err := v.Find(ctx, []string{"VirtualMachine"}, nil)
	// fmt.Println("vmsRefs ", vmsRefs, "len ", len(vmsRefs))
	if err != nil {
		fmt.Println("err ", err)
	}
	perf_dict := make(map[string]interface{})
	// Create a PerfManager
	perfManager := performance.NewManager(c)
	counterinfo, _ := perfManager.CounterInfo(ctx)
	for _, counter := range counterinfo {
		counter_full := fmt.Sprintf("%s.%s.%s", counter.GroupInfo.GetElementDescription().Key, counter.NameInfo.GetElementDescription().Key, counter.RollupType)
		perf_dict[counter_full] = counter.Key
		// fmt.Println("counter_full[key]: ", counter_full)
		// fmt.Println("counter_full[value]: ", counter.Key)
		// perf_dict["ok"] = "ok"
		// break
	}
	// counter_full[key]:  net.throughput.vds.arpTimeout.summation
	// counter_full[value]:  632 counterId

	fmt.Println("datastore.read.average ", perf_dict["datastore.read.average"])

	// m := view.NewManager(c)

	v1, _ := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"HostSystem"}, true)
	defer v1.Destroy(ctx)
	var hss []mo.HostSystem
	err = v.Retrieve(ctx, []string{"HostSystem"}, []string{"summary"}, &hss)
	if err != nil {
		panic(err)
	}
	for _, hs := range hss {
		fmt.Println("hostInfo: ",
			"hostname: ", hs.Summary.Config.Name,
			" uptime: ", hs.Summary.QuickStats.Uptime,
			" uuid:", "dc_name_", hs.Name,
			" vendor: ", hs.Summary.Hardware.Vendor,
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

		// Create PerfQuerySpec
		tt := time.Now()
		tt1 := time.Now().Add(60 * time.Second)
		spec2 := types.PerfQuerySpec{
		Entity: 
			IntervalId: 20,
			MetricId:   []types.PerfMetricId{{Instance: "*"}, {CounterId: 180}},
			StartTime:  &tt,
			EndTime:    &tt1,
		}
		ww, _ := perfManager.Query(ctx, []types.PerfQuerySpec{spec2})
		fmt.Println("ww: ", ww)
		fmt.Println("ww len: ", len(ww))
	}
}

/*
	// Retrieve counters name list
	counters, err := perfManager.CounterInfoByName(ctx)
	if err != nil {
		fmt.Println("err ", err)
	}
	var names []string
	for name := range counters {
		names = append(names, name)
	}
	t := time.Now()
	t1 := time.Now().Add(60 * time.Second)
	query := types.PerfQuerySpec{
		IntervalId: 180,
		Entity:     types.ManagedObjectReference{Type: "VirtualMachine"},
		MetricId:   []types.PerfMetricId{{Instance: "*"}, {CounterId: 632}},
		StartTime:  &t,
		EndTime:    &t1,
	}
	pq, err := perfManager.Query(ctx, []types.PerfQuerySpec{query})
	fmt.Println("qp: ", pq)
	fmt.Println("qp len: ", len(pq))

	for _, q := range pq {
		fmt.Println("q: ", q)
		fmt.Println("q1: ", q.GetPerfEntityMetricBase())

	}

	// Create PerfQuerySpec
	spec := types.PerfQuerySpec{
		MaxSample:  1,
		MetricId:   []types.PerfMetricId{{Instance: "*"}, {CounterId: 632}},
		IntervalId: int32(*interval),
	}
	// perfManager.CounterInfoByKey()
	// perfManager.CounterInfoByName()
	// Query metrics
	sample, err := perfManager.SampleByName(ctx, spec, names, vmsRefs)
	if err != nil {
		fmt.Println("err ", err)

	}

	result, err := perfManager.ToMetricSeries(ctx, sample)
	if err != nil {
		fmt.Println("err ", err)
	}

	// Read result
	for _, metric := range result {
		name := metric.Entity

		for _, v := range metric.Value {
			counter := counters[v.Name]
			units := counter.UnitInfo.GetElementDescription().Label

			instance := v.Instance
			if instance == "" {
				instance = "-"
			}
			if len(v.Value) != 0 {
				fmt.Println("units ", units)
				fmt.Println("name ", name)

			}
			// break

			if len(v.Value) != 0 {
				fmt.Printf("%s\t%s\t%s\t%s\t%s\n",
					name, instance, v.Name, v.ValueCSV(), units)
			}
		}
	}

	fmt.Println("counters: ", len(counters))
	fmt.Println("names: ", len(names))
	// fmt.Println(counters)
	// fmt.Println(names[0])

	spec1 := types.PerfQuerySpec{
		MaxSample:  1,
		MetricId:   []types.PerfMetricId{{Instance: "*"}},
		IntervalId: 20,
	}
	// names := []string{"datastore.read.average"}
	var names1 []string
	for name := range counters {
		names1 = append(names1, name)
	}
	names1 = []string{"datastore.totalReadLatency.average"}
	fmt.Println("names1: ", names1)
	vmsRefs1 := vmsRefs[:1]

	sample1, err := perfManager.SampleByName(ctx, spec1, names1, vmsRefs1)
	res, _ := perfManager.ToMetricSeries(ctx, sample1)
	fmt.Printf("res: %+#v", res, "len: ", len(res))
	fmt.Println("res read:", res[0].Value[0].Value[0])
*/
