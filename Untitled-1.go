package main

import (
	"encoding/json"
	"fmt"
)

type ListZoneResp struct {
	Data []struct {
		ID              string `json:"id"`
		AvailableZoneID string `json:"availableZoneId"`
		LogicalZone     struct {
			AvailableZoneID string `json:"availableZoneId"`
		} `json:"logicalZone"`
		AvailableZone struct {
			CloudAccount struct {
				Endpoint string `json:"endpoint"`
			} `json:"cloudAccount"`
		} `json:"availableZone"`
	} `json:"data"`
}
type Zone struct {
	ZoneID   string `json:"ZoneId"`
	Endpoint string `json:"endpoint"`
}

func main() {
	res := `{"data":[{"id":"1b286cac-ae77-464a-a631-9765ea55af07","ref":"0ee343cc-45f3-4220-aa29-e16575534793","name":"SDN集群","status":"Up","componentCode":"BareMetal","componentName":"BareMetal","version":"V3.0","lastSyncTime":"2022-06-09 00:00:10","createTime":"2022-05-25 13:17:27","updateTime":"2022-05-25 13:17:27","availableZoneId":"0260c681-30da-41be-a38c-7eed622bf2ae","cloudAccountId":"34c6ba30-a801-44bd-9507-d3676e563699","logicalZone":{},"availableZone":{"cloudAccount":{"IsEnabled":true,"apiKey":"bb31e4bb56074c10afcc88ad85f5071c","createTime":"2022-05-24 11:30:18","endpoint":"172.118.17.17","id":"9c3cdb89-e090-440c-b23e-a94bdf81d3eb","password":"","scope":"ArSdn","updateTime":"2022-05-24 11:30:18","username":""},"cloudAccountId":"9c3cdb89-e090-440c-b23e-a94bdf81d3eb","createTime":"2022-05-24 11:28:30","description":"","id":"0260c681-30da-41be-a38c-7eed622bf2ae","isDeleted":false,"name":"SDN可用区","ref":"","regionId":"e1f225fb-21e1-49d8-b56d-b3ea90044bb6","type":"ArSdn","updateTime":"2022-05-25 13:17:27"}},{"id":"47a4c13c-70e7-4801-b084-ff6f555f4d5e","ref":"9db7a695-24fc-4a8d-aedd-c1484425476f","name":"region7370","status":"Up","componentCode":"ArcherOS","componentName":"金山云","version":"1.5","lastSyncTime":"2022-06-09 00:10:28","createTime":"2022-06-06 10:24:36","updateTime":"2022-06-06 10:24:36","availableZoneId":"0d128ac5-4cfe-4160-a5b0-59dc34114a31","arch":"x86","cloudAccountId":"9582c1b1-e20c-4eb4-bfbe-3c706e4ebb17","logicalZone":{},"availableZone":{"cloudAccount":{"IsEnabled":true,"apiKey":"","createTime":"2022-06-02 11:45:16","endpoint":"172.118.73.70","id":"f9dea1a6-700a-4421-9e3a-b2f0b860f417","password":"SHVheXVuQDEyMw==","scope":"ArSdn","updateTime":"2022-06-02 11:45:16","username":"administrator@vsphere.local"},"cloudAccountId":"f9dea1a6-700a-4421-9e3a-b2f0b860f417","createTime":"2022-06-02 11:45:01","description":"","id":"0d128ac5-4cfe-4160-a5b0-59dc34114a31","isDeleted":false,"name":"灾备开发调试","ref":"","regionId":"e1f225fb-21e1-49d8-b56d-b3ea90044bb6","type":"ArSdn","updateTime":"2022-06-06 10:24:36"},"networkType":["VXLAN","VLAN"],"diskTypes":["ARSTOR"],"deployMode":"HA"},{"id":"4a14d88b-6d6e-4c92-a83d-a08222914003","ref":"172.118.69.31#datacenter-7","name":"Datacenter","status":"Up","componentCode":"VMware","componentName":"VMware","version":"6.7.0","lastSyncTime":"2022-06-09 00:02:44","createTime":"2022-06-02 11:45:16","updateTime":"2022-06-02 11:45:16","availableZoneId":"0d128ac5-4cfe-4160-a5b0-59dc34114a31","cloudAccountId":"a5fcff29-4fd7-4f36-b486-74ef510b82b2","logicalZone":{},"availableZone":{"cloudAccount":{"IsEnabled":true,"apiKey":"","createTime":"2022-06-02 11:45:16","endpoint":"172.118.73.70","id":"f9dea1a6-700a-4421-9e3a-b2f0b860f417","password":"SHVheXVuQDEyMw==","scope":"ArSdn","updateTime":"2022-06-02 11:45:16","username":"administrator@vsphere.local"},"cloudAccountId":"f9dea1a6-700a-4421-9e3a-b2f0b860f417","createTime":"2022-06-02 11:45:01","description":"","id":"0d128ac5-4cfe-4160-a5b0-59dc34114a31","isDeleted":false,"name":"灾备开发调试","ref":"","regionId":"e1f225fb-21e1-49d8-b56d-b3ea90044bb6","type":"ArSdn","updateTime":"2022-06-06 10:24:36"},"networkType":["VLAN","VXLAN"]},{"id":"4a6ce62a-ff4e-4572-8329-91b30974a02c","ref":"d1a8ae37-f61b-4dd9-8e2a-ca0d7953cd9e","name":"ArNet集群","status":"Up","componentCode":"BareMetal","componentName":"BareMetal","version":"V2.0","lastSyncTime":"2022-06-09 11:42:41","createTime":"2022-06-07 15:57:24","updateTime":"2022-06-07 15:57:24","availableZoneId":"21942fcb-a7f3-4fe4-977f-01f0e14cfb5f","cloudAccountId":"e5b4a272-3ae3-4e87-9fe1-00228e490f36","logicalZone":{},"availableZone":{"cloudAccountId":"","createTime":"2022-05-24 11:22:58","description":"","id":"21942fcb-a7f3-4fe4-977f-01f0e14cfb5f","isDeleted":false,"name":"ArNet可用区","ref":"","regionId":"e1f225fb-21e1-49d8-b56d-b3ea90044bb6","type":"ArNet","updateTime":"2022-06-07 15:57:24"}},{"id":"697144e8-9ffc-4eb6-9b10-af5b8ee856b3","ref":"8d9954cb-3c37-49ab-a6e3-82a36c30c264","name":"CloudOS0525","status":"Up","componentCode":"ArcherOS","componentName":"金山云","version":"1.5","lastSyncTime":"2022-06-09 00:10:29","createTime":"2022-05-27 09:26:18","updateTime":"2022-05-27 09:26:18","availableZoneId":"21942fcb-a7f3-4fe4-977f-01f0e14cfb5f","arch":"x86","cloudAccountId":"50f87179-cbee-4f94-b06e-d0d4b2f53bf3","logicalZone":{},"availableZone":{"cloudAccountId":"","createTime":"2022-05-24 11:22:58","description":"","id":"21942fcb-a7f3-4fe4-977f-01f0e14cfb5f","isDeleted":false,"name":"ArNet可用区","ref":"","regionId":"e1f225fb-21e1-49d8-b56d-b3ea90044bb6","type":"ArNet","updateTime":"2022-06-07 15:57:24"},"networkType":["VLAN"],"diskTypes":["ARSTOR"],"deployMode":"HA"},{"id":"7fad0335-9f62-411d-8403-88adf01e5b00","ref":"1e314c6a-c46f-49bb-92ea-2d398fac5820","name":"CloudOS","status":"Up","componentCode":"ArcherOS","componentName":"金山云","version":"1.5","lastSyncTime":"2022-06-09 00:10:31","createTime":"2022-05-24 11:30:18","updateTime":"2022-05-24 11:30:18","availableZoneId":"0260c681-30da-41be-a38c-7eed622bf2ae","arch":"x86","cloudAccountId":"e0a4be6d-e046-4e23-8245-fc1fe5149f0e","logicalZone":{},"availableZone":{"cloudAccount":{"IsEnabled":true,"apiKey":"bb31e4bb56074c10afcc88ad85f5071c","createTime":"2022-05-24 11:30:18","endpoint":"172.118.17.17","id":"9c3cdb89-e090-440c-b23e-a94bdf81d3eb","password":"","scope":"ArSdn","updateTime":"2022-05-24 11:30:18","username":""},"cloudAccountId":"9c3cdb89-e090-440c-b23e-a94bdf81d3eb","createTime":"2022-05-24 11:28:30","description":"","id":"0260c681-30da-41be-a38c-7eed622bf2ae","isDeleted":false,"name":"SDN可用区","ref":"","regionId":"e1f225fb-21e1-49d8-b56d-b3ea90044bb6","type":"ArSdn","updateTime":"2022-05-25 13:17:27"},"networkType":["VXLAN","VLAN"],"diskTypes":["ARSTOR"],"deployMode":"HA"},{"id":"b2cd8e29-40cd-42f0-bad9-ab751d71713a","ref":"6ee4c4e1-91fa-4ca4-b427-27cd22814b20","name":"CloudOS152","status":"Up","componentCode":"ArcherOS","componentName":"金山云","version":"1.5","lastSyncTime":"2022-06-09 15:23:37","createTime":"2022-05-16 14:40:07","updateTime":"2022-05-16 14:40:07","availableZoneId":"0aaf9c71-e651-48c9-989b-a6cfcb02125c","arch":"x86","cloudAccountId":"52ab4691-0ef5-4af2-a0f6-c72cff1fe434","logicalZone":{},"availableZone":{"cloudAccount":{"IsEnabled":true,"apiKey":"ef9723f1013341f8a28b712f0f3466e1","createTime":"2022-05-16 14:40:07","endpoint":"172.118.69.30","id":"9b8f8398-8ef2-4d01-ba43-c1b3ac39b5fc","password":"","scope":"ArSdn","updateTime":"2022-05-16 14:40:07","username":""},"cloudAccountId":"9b8f8398-8ef2-4d01-ba43-c1b3ac39b5fc","createTime":"2022-05-16 14:30:32","description":"","id":"0aaf9c71-e651-48c9-989b-a6cfcb02125c","isDeleted":false,"name":"可用区","ref":"","regionId":"e1f225fb-21e1-49d8-b56d-b3ea90044bb6","type":"ArSdn","updateTime":"2022-05-16 14:40:07"},"networkType":["VXLAN","VLAN"],"diskTypes":["ARSTOR"],"deployMode":"HA"}],"requestId":"b2924dc9-a399-4112-b44c-77bce6868ec2","timestamp":1654767531000}`
	var resp ListZoneResp
	_ = json.Unmarshal([]byte(res), &resp)

	var zs []*Zone

	for _, data := range resp.Data {
		if data.AvailableZone.CloudAccount.Endpoint != "" && data.AvailableZoneID != "" {
			zs = append(zs, &Zone{
				ZoneID:   data.ID,
				Endpoint: data.AvailableZone.CloudAccount.Endpoint,
			})
		}
	}

	for _, zone := range zs {
		fmt.Println(zone.ZoneID, zone.Endpoint)
	}

}
