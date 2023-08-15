package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"os"
	"strings"
)

func main() {

	metric := flag.String("metric", "", "The metric to search for in monitors")
	flag.Parse()

	if *metric == "" {
		fmt.Println("Usage: go run main.go -metric=<metric>")
		return
	}

	fmt.Println("Retrieving Monitors")
	ctx := datadog.NewDefaultContext(context.Background())
	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)
	api := datadogV1.NewMonitorsApi(apiClient)
	resp, r, err := api.ListMonitors(ctx, *datadogV1.NewListMonitorsOptionalParameters())

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MonitorsApi.ListMonitors`: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	} else {
		fmt.Fprintf(os.Stdout, "Searching metric in %d Monitors: %s\n", len(resp), *metric)
	}

	for _, element := range resp {
		name := element.Name
		query := element.Query
		if strings.Contains(query, *metric) {
			fmt.Fprintf(os.Stderr, "%s is being used in %s \n", *metric, *name)
		}
	}
	fmt.Println("Retrieving Dashboards")

	apidash := datadogV1.NewDashboardsApi(apiClient)
	resp2, r, err := apidash.ListDashboards(ctx, *datadogV1.NewListDashboardsOptionalParameters().WithFilterShared(false))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DashboardsApi.ListDashboards`: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	} else {
		fmt.Fprintf(os.Stdout, "Searching metric in %d Dashboards: %s\n", len(resp2.GetDashboards()), *metric)
	}

	for _, element := range resp2.GetDashboards() {
		resp3, _, _ := apidash.GetDashboard(ctx, *element.Id)
		responseContent, _ := json.MarshalIndent(resp3, "", "")
		if strings.Contains(string(responseContent), *metric) {
			fmt.Fprintf(os.Stderr, "%s is being used in %s \n", *metric, *element.Title)
		}
	}
}
