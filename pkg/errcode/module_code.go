package errcode

var (
	ErrorGetEventListFail     = NewError(30010001, "获取Event列表失败")
	ErrorGetEssentialListFail = NewError(30010002, "获取Essential列表失败")
	ErrorUploadFileFail       = NewError(30030003, "上传文件失败")
)
