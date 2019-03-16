
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:33</date>
//</624450071282323456>

package main

//标准的“mime”包依赖于系统设置，请参阅mime.osinitmime
//Swarm将在许多操作系统/平台/Docker上运行，并且必须表现出类似的行为。
//此命令生成代码以添加基于mime.types文件的常见mime类型
//
//mailcap提供的mime.types文件，遵循https://www.iana.org/assignments/media-types/media-types.xhtml
//
//
//docker run--rm-v$（pwd）：/tmp-alpine:edge/bin/sh-c“apk-add-u mailcap；mv/etc/mime.types/tmp”

import (
	"bufio"
	"bytes"
	"flag"
	"html/template"
	"io/ioutil"
	"strings"

	"log"
)

var (
	typesFlag   = flag.String("types", "", "Input mime.types file")
	packageFlag = flag.String("package", "", "Golang package in output file")
	outFlag     = flag.String("out", "", "Output file name for the generated mime types")
)

type mime struct {
	Name string
	Exts []string
}

type templateParams struct {
	PackageName string
	Mimes       []mime
}

func main() {
//分析并确保指定了所有需要的输入
	flag.Parse()
	if *typesFlag == "" {
		log.Fatalf("--types is required")
	}
	if *packageFlag == "" {
		log.Fatalf("--types is required")
	}
	if *outFlag == "" {
		log.Fatalf("--out is required")
	}

	params := templateParams{
		PackageName: *packageFlag,
	}

	types, err := ioutil.ReadFile(*typesFlag)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(types))
	for scanner.Scan() {
		txt := scanner.Text()
		if strings.HasPrefix(txt, "#") || len(txt) == 0 {
			continue
		}
		parts := strings.Fields(txt)
		if len(parts) == 1 {
			continue
		}
		params.Mimes = append(params.Mimes, mime{parts[0], parts[1:]})
	}

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}

	result := bytes.NewBuffer([]byte{})

	if err := template.Must(template.New("_").Parse(tpl)).Execute(result, params); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(*outFlag, result.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
}

var tpl = `//代码由github.com/ethereum/go-ethereum/cmd/swarm/mimegen生成。不要编辑。

package {{ .PackageName }}

import "mime"
func init() {
	var mimeTypes = map[string]string{
{{- range .Mimes -}}
	{{ $name := .Name -}}
	{{- range .Exts }}
		".{{ . }}": "{{ $name | html }}",
	{{- end }}
{{- end }}
	}
	for ext, name := range mimeTypes {
		if err := mime.AddExtensionType(ext, name); err != nil {
			panic(err)
		}
	}
}
`

