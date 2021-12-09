package errcode

var (
	ErrorGetEventListFail     = NewError(30010001, "获取Event列表失败")
	ErrorGetEssentialListFail = NewError(30010002, "获取Essential列表失败")
	ErrorGetProgressFail      = NewError(30010003, "获取Progress失败")
	ErrorUploadFail           = NewError(30010004, "上传文件错误")
	ErrorSaveUploadingFile    = NewError(30010005, "存储上传文件错误")
	ErrorImportCsvFile        = NewError(30010006, "处理Csv文件错误")
	ErrorProcess              = NewError(30010007, "ProcessAll错误")
	ErrorExportCsv            = NewError(30010008, "CSV导出错误")
	ErrorGetEventField        = NewError(30010009, "Event Field 获取失败")
	ErrorGetBook              = NewError(30010010, "Book 获取失败")
	ErrorGetSeries            = NewError(30010011, "Series 获取失败")
	ErrorGetRuns              = NewError(30010012, "Runs 获取失败")

	ErrorGetGcRunningFail = NewError(40010004, "获取Gcrunning失败")
)
