package deepseek

import (
	"bufio"
	"context"
	"fmt"

	utils "github.com/cohesion-org/deepseek-go/utils"
)

// CreateChatCompletion sends a chat completion request and returns the generated response.
func (c *Client) CreateChatCompletion(
	ctx context.Context,
	request *ChatCompletionRequest,
) (*ChatCompletionResponse, int, error) {
	if request == nil {
		return nil, 0, fmt.Errorf("request cannot be nil")
	}

	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL(c.BaseURL).
		SetPath("chat/completions").
		SetBodyFromStruct(request).
		SetApiType(c.ApiType).
		Build(ctx)

	if err != nil {
		return nil, 0, fmt.Errorf("error building request: %w", err)
	}
	resp, err := HandleSendChatCompletionRequest(*c, req)

	if err != nil {
		return nil, 0, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, 0, HandleAPIError(resp)
	}

	updatedResp, err := HandleChatCompletionResponse(resp)

	if err != nil {
		return nil, 0, fmt.Errorf("error decoding response: %w", err)
	}

	return updatedResp, resp.StatusCode, err
}

// CreateChatCompletionStream sends a chat completion request with stream = true and returns the delta
func (c *Client) CreateChatCompletionStream(
	ctx context.Context,
	request *StreamChatCompletionRequest,
) (ChatCompletionStream, int, error) {

	request.Stream = true
	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL(c.BaseURL).
		SetPath("chat/completions").
		SetBodyFromStruct(request).
		SetApiType(c.ApiType).
		BuildStream(ctx)

	if err != nil {
		return nil, 0, fmt.Errorf("error building request: %w", err)
	}

	resp, err := HandleSendChatCompletionRequest(*c, req)
	if err != nil {
		return nil, 0, err
	}

	if resp.StatusCode >= 400 {
		return nil, resp.StatusCode, HandleAPIError(resp)
	}

	ctx, cancel := context.WithCancel(ctx)
	stream := &chatCompletionStream{
		ctx:    ctx,
		cancel: cancel,
		resp:   resp,
		reader: bufio.NewReader(resp.Body),
	}
	return stream, resp.StatusCode, nil
}

// CreateFIMCompletion is a beta feature. It sends a FIM completion request and returns the generated response.
// the base URL is set to "https://api.deepseek.com/beta/"
func (c *Client) CreateFIMCompletion(
	ctx context.Context,
	request *FIMCompletionRequest,
) (*FIMCompletionResponse, error) {
	if request.MaxTokens > 4000 {
		return nil, fmt.Errorf("max tokens must be <= 4000")
	}
	baseURL := "https://api.deepseek.com/beta/"

	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL(baseURL).
		SetPath("completions").
		SetBodyFromStruct(request).
		Build(ctx)
	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}
	resp, err := HandleSendChatCompletionRequest(*c, req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, HandleAPIError(resp)
	}
	updatedResp, err := HandleFIMCompletionRequest(resp)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return updatedResp, err
}

// CreateFIMStreamCompletion sends a FIM completion request with stream = true and returns the delta
func (c *Client) CreateFIMStreamCompletion(
	ctx context.Context,
	request *FIMStreamCompletionRequest,
) (FIMChatCompletionStream, error) {
	baseURL := "https://api.deepseek.com/beta/"

	request.Stream = true
	req, err := utils.NewRequestBuilder(c.AuthToken).
		SetBaseURL(baseURL).
		SetPath("/completions").
		SetBodyFromStruct(request).
		BuildStream(ctx)

	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}

	resp, err := HandleSendChatCompletionRequest(*c, req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, HandleAPIError(resp)
	}

	ctx, cancel := context.WithCancel(ctx)
	stream := &fimCompletionStream{
		ctx:    ctx,
		cancel: cancel,
		resp:   resp,
		reader: bufio.NewReader(resp.Body),
	}
	return stream, nil
}
