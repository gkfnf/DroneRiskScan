package scanner

// Config 扫描器配置
type Config struct {
	Threads     int    // 并发线程数
	Timeout     int    // 请求超时时间(秒)
	TemplateDir string // 模板目录路径
	Output      string // 输出文件路径
}