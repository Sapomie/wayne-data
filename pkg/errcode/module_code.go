package errcode

var (
	ErrorGetEventListFail     = NewError(30010001, "获取Event列表失败")
	ErrorGetEssentialListFail = NewError(30010002, "获取Essential列表失败")
	ErrorGetProgressFail      = NewError(30010003, "获取Progress失败")
	ErrorGetGcRunningFail     = NewError(30010004, "获取Gcrunning失败")
	ErrorUploadFileFail       = NewError(30030009, "上传文件失败")
)
