// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

// Package eksiface provides an interface to enable mocking the Amazon Elastic Kubernetes Service service client
// for testing your code.
//
// It is important to note that this interface will have breaking changes
// when the service model is updated and adds new API operations, paginators,
// and waiters.
package eksiface

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/eks"
)

// EKSAPI provides an interface to enable mocking the
// eks.EKS service client's API operation,
// paginators, and waiters. This make unit testing your code that calls out
// to the SDK's service client's calls easier.
//
// The best way to use this interface is so the SDK's service client's calls
// can be stubbed out for unit testing your code with the SDK without needing
// to inject custom request handlers into the SDK's request pipeline.
//
//    // myFunc uses an SDK service client to make a request to
//    // Amazon Elastic Kubernetes Service.
//    func myFunc(svc eksiface.EKSAPI) bool {
//        // Make svc.CreateCluster request
//    }
//
//    func main() {
//        sess := session.New()
//        svc := eks.New(sess)
//
//        myFunc(svc)
//    }
//
// In your _test.go file:
//
//    // Define a mock struct to be used in your unit tests of myFunc.
//    type mockEKSClient struct {
//        eksiface.EKSAPI
//    }
//    func (m *mockEKSClient) CreateCluster(input *eks.CreateClusterInput) (*eks.CreateClusterOutput, error) {
//        // mock response/functionality
//    }
//
//    func TestMyFunc(t *testing.T) {
//        // Setup Test
//        mockSvc := &mockEKSClient{}
//
//        myfunc(mockSvc)
//
//        // Verify myFunc's functionality
//    }
//
// It is important to note that this interface will have breaking changes
// when the service model is updated and adds new API operations, paginators,
// and waiters. Its suggested to use the pattern above for testing, or using
// tooling to generate mocks to satisfy the interfaces.
type EKSAPI interface {
	CreateCluster(*eks.CreateClusterInput) (*eks.CreateClusterOutput, error)
	CreateClusterWithContext(aws.Context, *eks.CreateClusterInput, ...request.Option) (*eks.CreateClusterOutput, error)
	CreateClusterRequest(*eks.CreateClusterInput) (*request.Request, *eks.CreateClusterOutput)

	DeleteCluster(*eks.DeleteClusterInput) (*eks.DeleteClusterOutput, error)
	DeleteClusterWithContext(aws.Context, *eks.DeleteClusterInput, ...request.Option) (*eks.DeleteClusterOutput, error)
	DeleteClusterRequest(*eks.DeleteClusterInput) (*request.Request, *eks.DeleteClusterOutput)

	DescribeCluster(*eks.DescribeClusterInput) (*eks.DescribeClusterOutput, error)
	DescribeClusterWithContext(aws.Context, *eks.DescribeClusterInput, ...request.Option) (*eks.DescribeClusterOutput, error)
	DescribeClusterRequest(*eks.DescribeClusterInput) (*request.Request, *eks.DescribeClusterOutput)

	DescribeUpdate(*eks.DescribeUpdateInput) (*eks.DescribeUpdateOutput, error)
	DescribeUpdateWithContext(aws.Context, *eks.DescribeUpdateInput, ...request.Option) (*eks.DescribeUpdateOutput, error)
	DescribeUpdateRequest(*eks.DescribeUpdateInput) (*request.Request, *eks.DescribeUpdateOutput)

	ListClusters(*eks.ListClustersInput) (*eks.ListClustersOutput, error)
	ListClustersWithContext(aws.Context, *eks.ListClustersInput, ...request.Option) (*eks.ListClustersOutput, error)
	ListClustersRequest(*eks.ListClustersInput) (*request.Request, *eks.ListClustersOutput)

	ListClustersPages(*eks.ListClustersInput, func(*eks.ListClustersOutput, bool) bool) error
	ListClustersPagesWithContext(aws.Context, *eks.ListClustersInput, func(*eks.ListClustersOutput, bool) bool, ...request.Option) error

	ListUpdates(*eks.ListUpdatesInput) (*eks.ListUpdatesOutput, error)
	ListUpdatesWithContext(aws.Context, *eks.ListUpdatesInput, ...request.Option) (*eks.ListUpdatesOutput, error)
	ListUpdatesRequest(*eks.ListUpdatesInput) (*request.Request, *eks.ListUpdatesOutput)

	ListUpdatesPages(*eks.ListUpdatesInput, func(*eks.ListUpdatesOutput, bool) bool) error
	ListUpdatesPagesWithContext(aws.Context, *eks.ListUpdatesInput, func(*eks.ListUpdatesOutput, bool) bool, ...request.Option) error

	UpdateClusterConfig(*eks.UpdateClusterConfigInput) (*eks.UpdateClusterConfigOutput, error)
	UpdateClusterConfigWithContext(aws.Context, *eks.UpdateClusterConfigInput, ...request.Option) (*eks.UpdateClusterConfigOutput, error)
	UpdateClusterConfigRequest(*eks.UpdateClusterConfigInput) (*request.Request, *eks.UpdateClusterConfigOutput)

	UpdateClusterVersion(*eks.UpdateClusterVersionInput) (*eks.UpdateClusterVersionOutput, error)
	UpdateClusterVersionWithContext(aws.Context, *eks.UpdateClusterVersionInput, ...request.Option) (*eks.UpdateClusterVersionOutput, error)
	UpdateClusterVersionRequest(*eks.UpdateClusterVersionInput) (*request.Request, *eks.UpdateClusterVersionOutput)

	WaitUntilClusterActive(*eks.DescribeClusterInput) error
	WaitUntilClusterActiveWithContext(aws.Context, *eks.DescribeClusterInput, ...request.WaiterOption) error

	WaitUntilClusterDeleted(*eks.DescribeClusterInput) error
	WaitUntilClusterDeletedWithContext(aws.Context, *eks.DescribeClusterInput, ...request.WaiterOption) error
}

var _ EKSAPI = (*eks.EKS)(nil)
