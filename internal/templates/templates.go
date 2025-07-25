package templates

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Template 模板结构
type Template struct {
	ID          string        `yaml:"id"`
	Info        Info          `yaml:"info"`
	HTTP        []HTTPRequest `yaml:"http"`
}

// Info 模板信息
type Info struct {
	Name        string   `yaml:"name"`
	Author      string   `yaml:"author"`
	Severity    string   `yaml:"severity"`
	Description string   `yaml:"description"`
	Tags        []string `yaml:"tags"`
}

// HTTPRequest HTTP请求定义
type HTTPRequest struct {
	Method   string            `yaml:"method"`
	Path     []string          `yaml:"path"`
	Headers  map[string]string `yaml:"headers"`
	Body     string            `yaml:"body"`
	Matchers []Matcher         `yaml:"matchers"`
}

// Matcher 匹配器
type Matcher struct {
	Type     string   `yaml:"type"`
	Part     string   `yaml:"part"`
	Words    []string `yaml:"words"`
	Status   []int    `yaml:"status"`
}

// TemplateEngine 模板引擎
type TemplateEngine struct {
	templates map[string]*Template
}

// NewTemplateEngine 创建新的模板引擎
func NewTemplateEngine() *TemplateEngine {
	return &TemplateEngine{
		templates: make(map[string]*Template),
	}
}

// LoadTemplates 从目录加载模板
func (te *TemplateEngine) LoadTemplates(dir string) error {
	// 查找目录中的所有YAML文件
	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return err
	}

	// 加载每个模板文件
	for _, file := range files {
		err := te.loadTemplate(file)
		if err != nil {
			return err
		}
	}

	return nil
}

// loadTemplate 加载单个模板文件
func (te *TemplateEngine) loadTemplate(file string) error {
	// 读取文件内容
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	// 解析YAML
	var template Template
	err = yaml.Unmarshal(data, &template)
	if err != nil {
		return err
	}

	// 存储模板
	te.templates[template.ID] = &template

	return nil
}

// GetTemplate 获取模板
func (te *TemplateEngine) GetTemplate(id string) *Template {
	return te.templates[id]
}

// GetAllTemplates 获取所有模板
func (te *TemplateEngine) GetAllTemplates() []*Template {
	templates := make([]*Template, 0, len(te.templates))
	for _, template := range te.templates {
		templates = append(templates, template)
	}
	return templates
}

// ProcessPaths 处理路径中的变量
func (te *TemplateEngine) ProcessPaths(paths []string, baseURL string) []string {
	var processedPaths []string
	
	for _, path := range paths {
		// 替换{{BaseURL}}变量
		processedPath := strings.ReplaceAll(path, "{{BaseURL}}", baseURL)
		processedPaths = append(processedPaths, processedPath)
	}
	
	return processedPaths
}