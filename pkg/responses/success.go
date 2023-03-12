package responses

func CreateSuccessResponse(data interface{}) Response {
	return Response{Data: data, IsSuccess: true}
}
