package providers

import (
	"context"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"log"
)

func init() {
	Register("qcloud", &QcloudRegistration{})
}

type QcloudRegistration struct {
}

func (q *QcloudRegistration) CreateProvider(opts Options) (provider Provider, err error) {
	qcloudProvider := &QcloudProvider{}

	{ // CVM
		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
		if qcloudProvider.CVMClient, err = cvm.NewClient(
			common.NewCredential(opts.TokenID, opts.TokenSecret),
			opts.Region,
			cpf,
		); err != nil {
			return
		}
	}

	provider = qcloudProvider
	return
}

type QcloudProvider struct {
	CVMClient *cvm.Client
}

func (q *QcloudProvider) ListRecords(ctx context.Context, network string, service string) (records []Record, err error) {
	switch service {
	case "cvm":
		var offset int64
		for {
			req := cvm.NewDescribeInstancesRequest()
			req.Filters = []*cvm.Filter{
				{
					Name:   common.StringPtr("vpc-id"),
					Values: common.StringPtrs([]string{network}),
				},
			}
			req.Offset = common.Int64Ptr(offset)
			req.Limit = common.Int64Ptr(100)

			var res *cvm.DescribeInstancesResponse
			if res, err = q.CVMClient.DescribeInstances(cvm.NewDescribeInstancesRequest()); err != nil {
				return
			}
			log.Printf("%+v", res.Response)
			if len(res.Response.InstanceSet) == 0 {
				log.Println("break")
				break
			}
			offset += int64(len(res.Response.InstanceSet))
			log.Println("offset", offset)
			for _, item := range res.Response.InstanceSet {
				if item.InstanceName == nil {
					log.Println("missing InstanceName:", item.InstanceId)
					continue
				}
				if !regexpRecordName.MatchString(*item.InstanceName) {
					log.Printf("invalid InstanceName: %s, not a DNS name", *item.InstanceName)
					continue
				}
				if len(item.PrivateIpAddresses) == 0 {
					log.Println("missing PrivateIpAddresses:", item.InstanceId)
					continue
				}
				records = append(records, Record{
					Name:    *item.InstanceName,
					IP:      *item.PrivateIpAddresses[0],
					Comment: *item.InstanceId,
				})
			}
		}
	default:
		err = fmt.Errorf("unknown service: %s", service)
		return
	}
	return
}
