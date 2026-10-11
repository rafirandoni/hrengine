package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gofiber/fiber/v3"
	"github.com/rafirandoni/hrengine/constant"
	"github.com/rafirandoni/hrengine/utility"
)

type RJson struct {
	Metadata   *ROptionMeta       `json:"meta"`
	Data       interface{}        `json:"data"`
	Pagination *ROptionPagination `json:"pagination,omitempty"`
	RawError   interface{}        `json:"raw_error_message,omitempty"`
	Debug      interface{}        `json:"debug,omitempty"`
}

type ROption struct {
	StatusCode int
	Header     map[string]string
	Meta       *ROptionMeta
	Pagination *ROptionPagination
}

type ROptionMeta struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ROptionPagination struct {
	Limit        int     `json:"limit"`
	TotalItem    int     `json:"total_item"`
	CurrentPage  int     `json:"current_page"`
	TotalPage    int     `json:"total_page"`
	NextPageLink *string `json:"next_page_link"`
	PrevPageLink *string `json:"prev_page_link"`
}

func (slf *web) RSuccess(ctx fiber.Ctx, data interface{}, options ...*ROption) error {
	option := &ROption{
		StatusCode: 200,
		Header: map[string]string{
			constant.HeaderContentType: "application/json",
		},
		Meta: &ROptionMeta{
			Success: true,
		},
	}

	if len(options) > 0 {
		for key, val := range options[0].Header {
			option.Header[key] = val
		}

		if options[0].StatusCode > 0 {
			option.StatusCode = options[0].StatusCode
		}

		if options[0].Meta != nil {
			option.Meta = options[0].Meta
			option.Meta.Success = true
		}

		if options[0].Pagination != nil {
			option.Pagination = options[0].Pagination
		}
	}

	// CORS handling
	// w.Header().Set("Access-Control-Allow-Origin", "*")

	for key, val := range option.Header {
		ctx.Response().Header.Set(key, val)
	}

	res := &RJson{
		Metadata: option.Meta,
		Data:     data,
	}

	if option.Pagination != nil {
		res.Pagination = option.Pagination
	}

	ctx.Response().SetStatusCode(option.StatusCode)

	enc := json.NewEncoder(ctx.Response().BodyWriter())
	enc.SetEscapeHTML(false)
	_ = enc.Encode(res)

	return nil
}

func (slf *web) RFail(ctx fiber.Ctx, err error) error {
	if ctx.Response().Body() != nil && string(ctx.Response().Body()) != "" {
		return err
	}

	res := &RJson{
		Metadata: &ROptionMeta{
			Success: false,
			Message: "something went wrong",
		},
	}

	statusCode := http.StatusInternalServerError
	if resErr, ok := err.(*utility.ResError); ok {
		if resErr.StatusCode != nil {
			statusCode = *resErr.StatusCode
		}

		if resErr.Message != nil {
			res.Metadata.Message = *resErr.Message
		}

		if resErr.StatusCode != nil && *resErr.StatusCode != 500 &&
			resErr.Message == nil && resErr.Err != nil {
			res.Metadata.Message = resErr.Error()
		}

		if resErr.StackTrace == nil {
			resErr.StackTrace = utility.ToPointer(string(debug.Stack()))
		}

		if resErr.StackTrace != nil {
			res.Debug = resErr.StackTrace
		}
	}

	if !res.Metadata.Success && res.Debug == nil {
		res.Debug = utility.ToPointer(string(debug.Stack()))
	}

	if res.Debug != nil {
		log.Printf("%s\n%+v", err.Error(), *res.Debug.(*string))
	}

	ctx.Set("Content-Type", fiber.MIMEApplicationJSON)
	ctx.Status(statusCode)

	// enc := json.NewEncoder(ctx.Response().BodyWriter())
	// enc.SetEscapeHTML(false)
	// _ = enc.Encode(res)
	resJson, err := json.Marshal(res)
	if err != nil {
		ctx.WriteString(`{"meta": {"success": false, "message": "failed to build response"}}`)
		return err
	}

	ctx.Write(resJson)

	return fmt.Errorf("%s", string(resJson))
}
