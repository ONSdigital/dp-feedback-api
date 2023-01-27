# dp-feedback-api SDK

## Overview

This SDK contains a Go client for interacting with the Feedback API. The client contains a method for each API endpoint
so that any Go application wanting to interact with the feedback api can do so. Please refer to the [swagger specification](../swagger.yaml)
as the source of truth of how each endpoint works.

## Example use of the API SDK

Initialise new Feedback API client

```go
package main

import (
	"context"
	"github.com/ONSdigital/dp-feedback-api/sdk"
)

func main() {
    ...
	feedbackAPIClient := sdk.New("http://localhost:28600")
    ...
}
```

### Post feedback

Use the PostFeedback method to send a request to send a feedback email via the feedback API. This is a private endpoint and requires authorisation header.

```go
...
    // Create the feedback model you want to send
    f := &models.Feedback{
		IsPageUseful:      &isPageUsefulVal,
		IsGeneralFeedback: &isGeneralFeedbackVal,
        ...
	}

    // Pass the authorisation token (without the 'Bearer ' prefix) as an SDK Option parameter
    opts := sdk.Options{AuthToken: authToken}

    // Call PostFeedback to send the POST request to the feedback API
    err := apiClient.PostFeedback(ctx, f, opts)
    if err != nil {
        // handle error
    }
...
```

### Handling errors

The error returned from the method contains status code that can be accessed via `Status()` method and similar to extracting the error message using `Error()` method; see snippet below:

```go
...
    err := apiClient.PostFeedback(ctx, f, opts)
    if err != nil {
        // Retrieve status code from error
        statusCode := err.Status()
        // Retrieve error message from error
        errorMessage := err.Error()

        // log message, below uses "github.com/ONSdigital/log.go/v2/log" package
        log.Error(ctx, "failed to provide feedback", err, log.Data{"code": statusCode})

        return err
    }
...
```

### Healthcheck

This client extends the default Healthcheck Client. Please view this [README](https://github.com/ONSdigital/dp-api-clients-go/tree/main/health) for more information.

If your app is using the [dp-healthcheck library](https://github.com/ONSdigital/dp-healthcheck), then you may register the checker provided by this sdk, like so:

```go
...
    // assuming hc is an instance of HealthCheck, then you may register the feedback API checker
    err := hc.AddCheck("Feedback API", apiClient.Checker)
    if err != nil {
        // handle error
	}
...
```
