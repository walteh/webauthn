package awselasticloadbalancingv2


// Options for looking up an NetworkLoadBalancer.
//
// Example:
//   // The code below shows an example of how to instantiate this type.
//   // The values are placeholders you should change.
//   import "github.com/aws/aws-cdk-go/awscdk"
//
//   networkLoadBalancerLookupOptions := &NetworkLoadBalancerLookupOptions{
//   	LoadBalancerArn: jsii.String("loadBalancerArn"),
//   	LoadBalancerTags: map[string]*string{
//   		"loadBalancerTagsKey": jsii.String("loadBalancerTags"),
//   	},
//   }
//
type NetworkLoadBalancerLookupOptions struct {
	// Find by load balancer's ARN.
	// Default: - does not search by load balancer arn.
	//
	LoadBalancerArn *string `field:"optional" json:"loadBalancerArn" yaml:"loadBalancerArn"`
	// Match load balancer tags.
	// Default: - does not match load balancers by tags.
	//
	LoadBalancerTags *map[string]*string `field:"optional" json:"loadBalancerTags" yaml:"loadBalancerTags"`
}

