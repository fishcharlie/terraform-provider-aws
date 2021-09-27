package aws

import (
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tftags "github.com/hashicorp/terraform-provider-aws/aws/internal/tags"
	tfec2 "github.com/hashicorp/terraform-provider-aws/aws/internal/service/ec2"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func DataSourceNetworkACLs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkACLsRead,
		Schema: map[string]*schema.Schema{
			"filter": ec2CustomFiltersSchema(),

			"tags": tftags.TagsSchemaComputed(),

			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func dataSourceNetworkACLsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	req := &ec2.DescribeNetworkAclsInput{}

	if v, ok := d.GetOk("vpc_id"); ok {
		req.Filters = tfec2.BuildAttributeFilterList(
			map[string]string{
				"vpc-id": v.(string),
			},
		)
	}

	filters, filtersOk := d.GetOk("filter")
	tags, tagsOk := d.GetOk("tags")

	if tagsOk {
		req.Filters = append(req.Filters, buildEC2TagFilterList(
			tftags.New(tags.(map[string]interface{})).Ec2Tags(),
		)...)
	}

	if filtersOk {
		req.Filters = append(req.Filters, buildEC2CustomFilterList(
			filters.(*schema.Set),
		)...)
	}

	if len(req.Filters) == 0 {
		// Don't send an empty filters list; the EC2 API won't accept it.
		req.Filters = nil
	}

	log.Printf("[DEBUG] DescribeNetworkAcls %s\n", req)
	resp, err := conn.DescribeNetworkAcls(req)
	if err != nil {
		return err
	}

	if resp == nil || len(resp.NetworkAcls) == 0 {
		return errors.New("no matching network ACLs found")
	}

	networkAcls := make([]string, 0)

	for _, networkAcl := range resp.NetworkAcls {
		networkAcls = append(networkAcls, aws.StringValue(networkAcl.NetworkAclId))
	}

	d.SetId(meta.(*conns.AWSClient).Region)

	if err := d.Set("ids", networkAcls); err != nil {
		return fmt.Errorf("Error setting network ACL ids: %w", err)
	}

	return nil
}