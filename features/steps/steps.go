package steps

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	c.apiFeature.RegisterSteps(ctx)

	ctx.Step(`^I should receive a (\d+) status code with an empty body response`, c.iShouldReceiveAnEmptyResponse)
	ctx.Step(`^I should receive a (\d+) status code with an the following body response$`, c.iShouldReceiveResponse)
	ctx.Step(`^the following email is sent$`, c.theFollowingEmailIsSent)
	ctx.Step(`^no email is sent`, c.noEmailIsSent)
}

func (c *Component) iShouldReceiveAnEmptyResponse(code string) error {
	return c.iShouldReceiveResponse(
		code,
		&godog.DocString{
			Content: "",
		},
	)
}

func (c *Component) iShouldReceiveResponse(code string, documentJSON *godog.DocString) error {
	// Validate status code
	statusCode := c.apiFeature.HttpResponse.StatusCode
	expectedCode, err := strconv.Atoi(code)
	if err != nil {
		return fmt.Errorf("cannot parse expected code: %w", err)
	}
	assert.Equal(c, expectedCode, statusCode)

	// Validate body
	var expectedBody = trimLines(documentJSON.Content)
	responseBody := c.apiFeature.HttpResponse.Body
	body, err := io.ReadAll(responseBody)
	if err != nil {
		return fmt.Errorf("cannot read body from response: %w", err)
	}
	assert.Equal(c, expectedBody, trimLines(string(body)))

	return c.StepError()
}

func (c *Component) theFollowingEmailIsSent(documentJSON *godog.DocString) error {
	assert.Equal(c, len(c.EmailSenderMock.SendCalls()), 1)

	var expectedEmail = trimLines(documentJSON.Content)
	sentEmail := c.EmailSenderMock.SendCalls()[0].Msg
	assert.Equal(c, expectedEmail, trimLines(string(sentEmail)))

	return c.StepError()
}

func (c *Component) noEmailIsSent() error {
	assert.Equal(c, len(c.EmailSenderMock.SendCalls()), 0)
	return c.StepError()
}

func trimLines(in string) string {
	var sb strings.Builder
	for _, line := range strings.Split(strings.TrimSpace(in), "\n") {
		sb.WriteString(strings.TrimSpace(line))
		sb.WriteByte('\n')
	}
	return sb.String()
}
