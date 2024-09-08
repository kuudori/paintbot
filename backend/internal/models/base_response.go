package models

type BaseResponse struct {
	Code           int         `json:"code" binding:"required"`
	Message        *string     `json:"message"`
	MessageDetails *string     `json:"messageDetails"`
	ResponseData   interface{} `json:"responseData" binding:"required"`
}
