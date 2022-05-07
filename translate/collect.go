package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/vmware/govmomi/session/cache"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

var vcenter1 = map[string]string{
	"vcenterIp":       os.Getenv("GOVMOMI_URL"),
	"vcenterUser":     os.Getenv("GOVMOMI_USERNAME"),
	"vcenterPassword": os.Getenv("GOVMOMI_PASSWORD"),
}
var Perf_dict map[string]string

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
func collect_host_vm(vcenter map[string]string, start_ts int64) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, _ := vlogin(ctx, vcenter["vcenterIp"], vcenter["vcenterUser"], vcenter["vcenterPassword"])
	fmt.Println("Login successful =================")

	content := client.ServiceContent

	// global payload,perf_dict
	// content.PerfManager
	Perf_dict = make(map[string]string)

	perfList := content.PerfManager

	Perf_dict[perfList.Type] = perfList.Value

	var timeout = 300
	vchtime := time.Now().UnixNano()
	interval := 1000
	do_collect_host_vm_parallel(vcenter, content, timeout, vchtime, interval, start_ts)

}

func do_collect_host_vm_parallel(vcenter map[string]string, content types.ServiceContent, timeout int, vchtime int64, interval int, start_ts int64) {
	fmt.Println("do_collect_host_vm_parallel -------------------------------- 1")
	collectSendHost(vcenter, content, timeout, vchtime, interval, start_ts)
	fmt.Println("do_collect_host_vm_parallel -------------------------------- 2")
	collectVm(vcenter, content, timeout, vchtime, interval)
	fmt.Println("do_collect_host_vm_parallel -------------------------------- 3")

}

func collectVm(vcenter map[string]string, content types.ServiceContent, timeout int, vchtime int64, interval int) []interface{} {
	fmt.Println("collectVm -------------------------------- ")

	vms := getVmInfo(vcenter["vcenterIp"], content, timeout, vchtime, interval)
	return vms
}

func getVmInfo(s string, content types.ServiceContent, timeout int, vchtime int64, interval int) []interface{} {
	fmt.Println("getVmInfo -------------------------------- ")

	fmt.Println(content.RootFolder)
	// vms := make([]interface{})

	return nil
}

func collectSendHost(vcenter map[string]string, content types.ServiceContent, timeout int, vchtime int64, interval int, start_ts int64) {
	fmt.Println("collectSendHost -------------------------------- 1")

}

/*
def collect_host_vm(vcenter, start_ts):
    si = None
    try:
        logger.debug("----------collect host vm: start collect,vcenterIp is:%s", vcenter["vcenterIp"])
        interval = cfgData["transfer"]["interval"]
        context = ssl._create_unverified_context()
        si = SmartConnect(host=vcenter["vcenterIp"],
                          user=vcenter["vcenterUser"],
                          pwd=vcenter["vcenterPassword"],
                          port=int(vcenter["vcenterPort"]),
                          connectionPoolTimeout = 30,
                          sslContext=context)
        content = si.RetrieveContent()
        vchtime = si.CurrentTime()

        logger.debug("----------collect host vm: after connection,vcenterIp is:%s", vcenter["vcenterIp"])
        vcenterInfo.initPerf(content)
        logger.debug("----------collect host vm: after init perf,vcenterIp is:%s", vcenter["vcenterIp"])

        timeout = 300
        # Get Host and Vm perf
        do_collect_host_vm_parallel(vcenter, content, timeout, vchtime, interval, start_ts)

        logger.debug("----------collect host vm: finish collect,vcenterIp is:%s", vcenter["vcenterIp"])
    except Exception, e:
        exstr = traceback.format_exc()
        logger.error("[collect_host_vm] vcenterIp: %s, err: %s" % (vcenter["vcenterIp"], exstr))
    finally:
        if si:
            Disconnect(si)

# 执行并行采集vm host
def do_collect_host_vm_parallel(vcenter, content, timeout, vchtime, interval, start_ts):
    # storage = gevent.spawn(storage_collector.collectStorage, vcenter, content)
    hosts = gevent.spawn(host_collector.collectSendHost, vcenter, content, timeout, vchtime, interval, start_ts)
    vms = gevent.spawn(vm_collector.collectVm, vcenter, content, timeout, vchtime, interval)
    vsan_perf = gevent.spawn(etcd_access.fetch_vsan_perf_from_etcd, vcenter["vcenterIp"])

    try:
        gevent.joinall([hosts, vms], timeout=60, raise_error=True)
    except Exception, e:
        exstr = traceback.format_exc()
        logger.error("[do_collect_host_vm_parallel] timeout, error: %s" % exstr)

    # 最后发送需要vms和vsan_perf同时结束
    for vm in vms.value:
        vm_collector.sendVm(start_ts, vm, vsan_perf.value)
    logger.debug("----------after sendVm,vcenterIp is:%s", vcenter["vcenterIp"])
*/
func main() {
	start_ts := time.Now().UnixNano()
	collect_host_vm(vcenter1, start_ts)
}
