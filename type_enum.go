package aliyundrive

type FileType string // 文件类型枚举

const (
	FileTypeFolder FileType = "folder" // 文件夹
	FileTypeFile   FileType = "file"   // 文件
)

const RootFileID = "root"
