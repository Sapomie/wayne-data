package errcode

var (
	ErrorGetEventListFail = NewError(30010001, "获取Event列表失败")
	ErrorUploadFileFail   = NewError(20030001, "上传文件失败")
)
