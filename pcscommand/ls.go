package pcscommand

import (
	"fmt"
	"github.com/iikira/BaiduPCS-Go/baidupcs"
	"github.com/iikira/BaiduPCS-Go/pcsutil"
	"os"
	"text/template"
)

// RunLs 执行列目录
func RunLs(path string) {
	path, err := getAbsPath(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	files, err := info.FilesDirectoriesList(path, false)
	if err != nil {
		fmt.Println(err)
		return
	}

	for k := range files {
		if files[k].Isdir {
			files[k].Filename += "/"
		}
	}

	tmpl, err := template.New("ls").Funcs(
		template.FuncMap{
			"convertFileSize": func(size int64) string {
				res := pcsutil.ConvertFileSize(size)
				if res == "0" {
					return "-       "
				}
				return res
			},
			"timeFmt": pcsutil.FormatTime,
			"totalSize": func() string {
				return pcsutil.ConvertFileSize(files.TotalSize())
			},
			"fdCount": func() string {
				fN, dN := files.Count()
				s := fmt.Sprintf("文件总数: %d,\t目录总数: %d", fN, dN)
				if fN+dN >= 50 {
					s += fmt.Sprintf(",\t当前目录: %s", path)
				}
				return s
			},
		},
	).Parse(`
------------------------------------------------------------------------------
当前目录: {{.ThisPath}}
----
文件大小	创建日期		文件(目录)
------------------------------------------------------------------------------
{{range .FilesInfo}}
{{convertFileSize .Size}}	{{timeFmt .Ctime}}	{{.Filename}}{{end}}
------------------------------------------------------------------------------
总大小: {{totalSize}}	{{fdCount}}
------------------------------------------------------------------------------
`)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, struct {
		ThisPath  string
		FilesInfo baidupcs.FileDirectoryList
	}{
		ThisPath:  path,
		FilesInfo: files,
	})
	if err != nil {
		panic(err)
	}
}

type Summary struct {
	TotalSize   int64	`json:"totalSize"`
	TotalSizeStr   string	`json:"totalSizeStr"`
	FileCount   int64	`json:"fileCount"`
	DirCount    int64	`json:"dirCount"`
}

// Ls 执行列目录
func Ls(path string) (interface{}, interface{}, error){
	path, err := getAbsPath(path)
	if err != nil {
		return nil, nil, err
	}

	files, err := info.FilesDirectoriesList(path, false)
	if err != nil {
		return nil, nil, err
	}

	fN, dN := files.Count()
	summary := new(Summary)
	summary.TotalSize = files.TotalSize()
	summary.TotalSizeStr = pcsutil.ConvertFileSize(files.TotalSize())
	summary.FileCount = fN
	summary.DirCount = dN
	// summary := new(Summary)

	return files, summary, nil
}
