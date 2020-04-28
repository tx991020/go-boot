package generator

import (
	"errors"
	"fmt"
	"sort"

	"github.com/flosch/pongo2"
	"github.com/iancoleman/strcase"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

var config *Configuration

// InputParams 输入参数
type InputParams struct {
	Name              string   `json:"name,omitempty"`
	HTTPPort          string   `json:"http_port,omitempty"`
	DBType            string   `json:"db_type,omitempty"`
	AppPath           string   `json:"app_path"`
	TemplatePath      string   `json:"template_path"`
	Host              string   `json:"host,omitempty"`
	Port              string   `json:"port,omitempty"`
	User              string   `json:"user,omitempty"`
	Password          string   `json:"password,omitempty"`
	Database          string   `json:"database,omitempty"`
	Sslmode           string   `json:"sslmode,omitempty"`
	Table             []string `json:"table,omitempty"`
	IsOpenMemoryCache bool     `json:"is_open_memory_cache,omitempty"`
}

type Configuration struct {
	driver            string `json:"driver,omitempty"`
	AppPath           string `json:"app_path,omitempty"`
	TemplatePath      string `json:"template_path,omitempty"`
	DbAddress         string `json:"db_address,omitempty"`
	DbName            string `json:"db_name,omitempty"`
	AppName           string `json:"app_name,omitempty"`
	HTTPPort          string `json:"http_port,omitempty"`
	TagLabel          string `json:"tag_label,omitempty"`
	IsOpenMemoryCache bool   `json:"is_open_memory_cache,omitempty"`
}

// Write 写入文件,写入模板文件
func (c *Configuration) Write(templateName, filePath string, params map[string]interface{}) error {

	fPath := c.AppPath + c.AppName + filePath
	var tpl *pongo2.Template
	var err error
	defer func() {
		if err == nil {
			fmt.Printf("%s \U00002714\n", c.AppName+filePath)
		} else {
			PrintErr(err)
		}
	}()
	tplPath := fmt.Sprintf("%s%s/template/%s.tpl", c.AppPath, c.AppName, templateName)
	PrintInfo(tplPath)
	// check 模板是否存在; 不存在则使用data数据
	if IsExists(tplPath) {
		tpl, err = pongo2.FromFile(tplPath)
		if err != nil {

			return err
		}
	} else {

	}

	if params == nil {
		params = map[string]interface{}{}
	}
	b, err := tpl.ExecuteBytes(pongo2.Context(params))
	if err != nil {
		return err
	}
	return CreateAndFormat(fPath, b, 0664)
}

func (c *Configuration) IsExistFile(path string) bool {
	return IsExists(fmt.Sprintf("%s%s%s", c.AppPath, c.AppName, path))
}

// MkdirAll 创建文件夹结构
func (c *Configuration) MkdirAll(filePath string) error {
	return MkdirAll(c.AppPath+c.AppName+filePath, 0755)
}

//生成目录
func initDirs() {
	MkdirAll(config.AppPath+config.AppName+"/", 0755)
	MkdirAll(config.AppPath+config.AppName+"/dao", 0755)
	MkdirAll(config.AppPath+config.AppName+"/cache", 0755)
	MkdirAll(config.AppPath+config.AppName+"/"+"handler/curd", 0755)
	MkdirAll(config.AppPath+config.AppName+"/"+"service", 0755)
	MkdirAll(config.AppPath+config.AppName+"/"+"config", 0755)
	MkdirAll(config.AppPath+config.AppName+"/"+"middleware", 0755)
}

func writeStruct(tables map[string][]*ColumnSchema, dbType string, tableInput []string) error {
	var cacheParam = []string{}
	var handlerParam = ""
	var keys []string
	for k := range tables {
		keys = append(keys, k)
	}
	tableInputs := map[string]int{}
	for _, v := range tableInput {
		tableInputs[v] = 1
	}
	sort.Strings(keys) // 按表名排序
	for _, tableName := range keys {
		if _, exists := tableInputs[tableName]; len(tableInput) > 0 && !exists {
			continue
		}
		strcutName, tagName := FormatName(tableName)
		cacheParam = append(cacheParam, tagName)
		handlerParam += "\n" + strcutName + "Init(r)"
	}
	generateCacheInit(cacheParam)
	generateHandlerInit(handlerParam)

	for _, tableName := range keys {
		if _, exists := tableInputs[tableName]; len(tableInput) > 0 && !exists {
			continue
		}
		var columns = tables[tableName]
		generateModel(tableName, dbType, columns)
		//generateCacheCRUD(tableName, dbType, columns[0])
		generateHandler(tableName, dbType, columns[0])
		generateService(tableName)

	}
	return nil
}

func convertType(column *ColumnSchema, dbType string) (string, string) {
	var s string
	switch dbType {
	case "mysql":
		s, _ = MySQLType(column)
	case "postgres":
		s, _ = PgType(column)

	}
	return s, ""
}

//生成 dao.Model
func generateModel(tableName string, dbType string, columns []*ColumnSchema) {

	s, _ := convertType(columns[0], dbType)
	ftn, ctn := FormatName(tableName)
	Pk, pk := FormatName(columns[0].ColumnName)
	structArr := []map[string]string{}

	for _, column := range columns {
		str := map[string]string{}
		fcn, ccn := FormatName(column.ColumnName)
		goType, _ := convertType(column, dbType)
		str["key"] = fcn
		str["type"] = goType
		if goType == "*time.Time" {
			str["jsonTag"] = "-"
		} else {
			str["jsonTag"] = ccn
		}
		if goType == "*time.Time" {
			str["sql"] = "-"
		}
		structArr = append(structArr, str)
	}

	config.Write("model_template", "/dao/"+tableName+".go", map[string]interface{}{
		"package":   "dao",
		"Table":     ftn,
		"table":     ctn,
		"snake":     tableName,
		"Pk":        Pk,
		"pk":        pk,
		"pk_snake":  strcase.ToSnake(Pk),
		"type":      s,
		"structArr": structArr,
	})
}

//生成cache CRUD
func generateCacheCRUD(tableName string, dbType string, column *ColumnSchema) {
	s, _ := convertType(column, dbType)
	ftn, ctn := FormatName(tableName)
	Pk, pk := FormatName(column.ColumnName)
	cacheTemplateName := "cache_template"
	if config.IsOpenMemoryCache {
		cacheTemplateName = "cache_template_memory"
	}
	if config.IsOpenMemoryCache == false {
		return
	}
	config.Write(cacheTemplateName, "/cache/"+tableName+".go", map[string]interface{}{
		"package": "cache",
		"table":   ctn,
		"Table":   ftn,
		"appPath":config.AppPath,
		"AppName": config.AppName,
		"Pk":      Pk,
		"pk":      pk,
		"type":    s,
	})
}

// 生成handler
func generateHandler(tableName string, dbType string, column *ColumnSchema) {
	s, _ := convertType(column, dbType)

	conventType1 := func(s string) string {
		switch s {
		case "int64":
			return "Int"
		case "string":
			return "String"
		}
		return ""
	}(s)
	ftn, ctn := FormatName(tableName)
	Pk, pk := FormatName(column.ColumnName)
	config.Write("handler_template", "/handler/curd/"+tableName+".go", map[string]interface{}{
		"package": "curd",
		"appPath":config.AppPath,
		"AppName": config.AppName,
		"Table":   ftn,
		"table":   ctn,
		"Pk":      Pk,
		"pk":      pk,
		"Type":    conventType1,
		"type":    s,
	})
}

//生成service
func generateService(tableName string) {
	isExist := config.IsExistFile(fmt.Sprintf("/%s.go", tableName))
	if isExist {
		return
	}
	config.Write("service_service", "/service/service.go", nil)
}

//生成cache Init 列表
func generateCacheInit(cacheParam []string) {
	config.Write("cache", "/cache/cache.go", map[string]interface{}{"namelist": cacheParam})
}

//生成handler Init 列表
func generateHandlerInit(handlerParam string) {
	config.Write("handlerInit", "/handler/curd/handler.go", map[string]interface{}{"namelist": handlerParam, "AppName": config.AppName})

}

//生成database
func generateDatabase() {
	config.Write("database", "/dao/database.go", map[string]interface{}{"package": "dao"})
}

//生成redis
func generateRedis() {
	config.Write("cache_redis", "/cache/redis.go", nil)
}

// 生成config.yml
func generateConfigExample(i *InputParams) {
	m, _ := Struct2Map(config)
	m2, _ := Struct2Map(i)
	if config.IsExistFile("/config/config.yml") {

	} else {
		config.Write("config", "/config/config.yml", MergeMap(m, m2))
	}

	if config.IsExistFile("/config/config_code.go") {

	} else {
		config.Write("config.code", "/config/config_code.go", nil)
	}
	if config.IsExistFile("/config/config.go") {

	} else {
		config.Write("config_config", "/config/config.go", nil)
	}

}

// 生成中间件
func generateMiddle() {
	config.Write("middleware_curd_verify", "/middleware/curd_verify.go", map[string]interface{}{"AppName": config.AppName})
	config.Write("middleware_error_trace", "/middleware/error_trace.go", nil)
}

//生成main
func generateMain() {
	isExist := config.IsExistFile("/main.go")
	if isExist {
		return
	}
	config.Write("main", "/main.go", map[string]interface{}{"AppName": config.AppName})
}

//部署
func generateDeploy() {
	config.Write("dockerfile", "/Dockerfile", map[string]interface{}{"AppName": config.AppName, "http_port": config.HTTPPort})
	config.Write("gitignore", "/.gitignore", map[string]interface{}{"AppName": config.AppName})
	config.Write("deploy", "/deploy.sh", map[string]interface{}{"AppName": config.AppName, "http_port": config.HTTPPort})

}

// Run 开始运行
func Run(i *InputParams) {
	appPath, templatePath, err := GetEnvPath(i.AppPath, i.TemplatePath)
	if err != nil {
		PrintErr(errors.New("查询golang路径失败"))
	}
	addr := ""
	if i.DBType == "postgres" {
		addr = fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=%s", i.Host, i.Port, i.User, i.Password, i.Sslmode)
	} else {
		addr = fmt.Sprintf("%s:%s@tcp(%s:%s)", i.User, i.Password, i.Host, i.Port)
	}

	config = &Configuration{
		driver:            i.DBType,
		AppPath:           appPath,
		TemplatePath:      templatePath,
		DbAddress:         addr,
		DbName:            i.Database,
		AppName:           i.Name,
		HTTPPort:          i.HTTPPort,
		TagLabel:          "json",
		IsOpenMemoryCache: i.IsOpenMemoryCache,
	}

	initDirs()
	generateDeploy()
	generateDatabase()
	generateConfigExample(i)
	generateRedis()
	generateMain()
	schema := getSchema(config.driver, config.DbAddress, config.DbName)
	err = writeStruct(schema, i.DBType, i.Table)
	if err != nil {
		PrintErr(err)
		panic(err)
	}
	fmt.Println("Create Success! \U0001F4AF")
}
