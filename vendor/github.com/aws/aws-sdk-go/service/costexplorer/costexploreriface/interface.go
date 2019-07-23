// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

// Package costexploreriface provides an interface to enable mocking the AWS Cost Explorer Service service client
// for testing your code.
//
// It is important to note that this interface will have breaking changes
// when the service model is updated and adds new API operations, paginators,
// and waiters.
package costexploreriface

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

// CostExplorerAPI provides an interface to enable mocking the
// costexplorer.CostExplorer service client's API operation,
// paginators, and waiters. This make unit testing your code that calls out
// to the SDK's service client's calls easier.
//
// The best way to use this interface is so the SDK's service client's calls
// can be stubbed out for unit testing your code with the SDK without needing
// to inject custom request handlers into the SDK's request pipeline.
//
//    // myFunc uses an SDK service client to make a request to
//    // AWS Cost Explorer Service.
//    func myFunc(svc costexploreriface.CostExplorerAPI) bool {
//        // Make svc.GetCostAndUsage request
//    }
//
//    func main() {
//        sess := session.New()
//        svc := costexplorer.New(sess)
//
//        myFunc(svc)
//    }
//
// In your _test.go file:
//
//    // Define a mock struct to be used in your unit tests of myFunc.
//    type mockCostExplorerClient struct {
//        costexploreriface.CostExplorerAPI
//    }
//    func (m *mockCostExplorerClient) GetCostAndUsage(input *costexplorer.GetCostAndUsageInput) (*costexplorer.GetCostAndUsageOutput, error) {
//        // mock response/functionality
//    }
//
//    func TestMyFunc(t *testing.T) {
//        // Setup Test
//        mockSvc := &mockCostExplorerClient{}
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
type CostExplorerAPI interface {
	GetCostAndUsage(*costexplorer.GetCostAndUsageInput) (*costexplorer.GetCostAndUsageOutput, error)
	GetCostAndUsageWithContext(aws.Context, *costexplorer.GetCostAndUsageInput, ...request.Option) (*costexplorer.GetCostAndUsageOutput, error)
	GetCostAndUsageRequest(*costexplorer.GetCostAndUsageInput) (*request.Request, *costexplorer.GetCostAndUsageOutput)

	GetCostForecast(*costexplorer.GetCostForecastInput) (*costexplorer.GetCostForecastOutput, error)
	GetCostForecastWithContext(aws.Context, *costexplorer.GetCostForecastInput, ...request.Option) (*costexplorer.GetCostForecastOutput, error)
	GetCostForecastRequest(*costexplorer.GetCostForecastInput) (*request.Request, *costexplorer.GetCostForecastOutput)

	GetDimensionValues(*costexplorer.GetDimensionValuesInput) (*costexplorer.GetDimensionValuesOutput, error)
	GetDimensionValuesWithContext(aws.Context, *costexplorer.GetDimensionValuesInput, ...request.Option) (*costexplorer.GetDimensionValuesOutput, error)
	GetDimensionValuesRequest(*costexplorer.GetDimensionValuesInput) (*request.Request, *costexplorer.GetDimensionValuesOutput)

	GetReservationCoverage(*costexplorer.GetReservationCoverageInput) (*costexplorer.GetReservationCoverageOutput, error)
	GetReservationCoverageWithContext(aws.Context, *costexplorer.GetReservationCoverageInput, ...request.Option) (*costexplorer.GetReservationCoverageOutput, error)
	GetReservationCoverageRequest(*costexplorer.GetReservationCoverageInput) (*request.Request, *costexplorer.GetReservationCoverageOutput)

	GetReservationPurchaseRecommendation(*costexplorer.GetReservationPurchaseRecommendationInput) (*costexplorer.GetReservationPurchaseRecommendationOutput, error)
	GetReservationPurchaseRecommendationWithContext(aws.Context, *costexplorer.GetReservationPurchaseRecommendationInput, ...request.Option) (*costexplorer.GetReservationPurchaseRecommendationOutput, error)
	GetReservationPurchaseRecommendationRequest(*costexplorer.GetReservationPurchaseRecommendationInput) (*request.Request, *costexplorer.GetReservationPurchaseRecommendationOutput)

	GetReservationUtilization(*costexplorer.GetReservationUtilizationInput) (*costexplorer.GetReservationUtilizationOutput, error)
	GetReservationUtilizationWithContext(aws.Context, *costexplorer.GetReservationUtilizationInput, ...request.Option) (*costexplorer.GetReservationUtilizationOutput, error)
	GetReservationUtilizationRequest(*costexplorer.GetReservationUtilizationInput) (*request.Request, *costexplorer.GetReservationUtilizationOutput)

	GetTags(*costexplorer.GetTagsInput) (*costexplorer.GetTagsOutput, error)
	GetTagsWithContext(aws.Context, *costexplorer.GetTagsInput, ...request.Option) (*costexplorer.GetTagsOutput, error)
	GetTagsRequest(*costexplorer.GetTagsInput) (*request.Request, *costexplorer.GetTagsOutput)

	GetUsageForecast(*costexplorer.GetUsageForecastInput) (*costexplorer.GetUsageForecastOutput, error)
	GetUsageForecastWithContext(aws.Context, *costexplorer.GetUsageForecastInput, ...request.Option) (*costexplorer.GetUsageForecastOutput, error)
	GetUsageForecastRequest(*costexplorer.GetUsageForecastInput) (*request.Request, *costexplorer.GetUsageForecastOutput)
}

var _ CostExplorerAPI = (*costexplorer.CostExplorer)(nil)
