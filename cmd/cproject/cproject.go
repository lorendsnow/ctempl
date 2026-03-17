package cproject

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/lorendsnow/ctempl/internal/folder"
)

//go:embed files/root/* files/tests/*
var f embed.FS

type CProject struct {
	std        *string
	compiler   *string
	cmakeVer   *float64
	lib        *bool
	libName    *string
	exe        *bool
	exeName    *string
	exeFlags   []string
	projName   string
	flagSet    *flag.FlagSet
	rootFolder string
}

func (p *CProject) Run() error {
	if err := p.createRoot(); err != nil {
		return err
	}

	if err := p.createSrc(); err != nil {
		return err
	}

	if *p.lib {
		if err := p.createLib(); err != nil {
			return err
		}
	}

	if err := p.createTests(); err != nil {
		return err
	}

	return nil
}

func (p *CProject) createRoot() error {
	projData := &ProjData{
		MinVersion:  *p.cmakeVer,
		Compiler:    *p.compiler,
		ProjectName: p.projName,
	}

	tmpl, err := template.New("Project CMake CMakeLists.txt").Parse(projCMakeListTempl)
	if err != nil {
		return fmt.Errorf("failed to create CMakeLists.txt template for project folder: %w", err)
	}

	tmplFile := &folder.TemplateFile{
		Filename: "CMakeLists.txt",
		Tmpl:     tmpl,
		TmplData: projData,
	}

	if err := folder.CreateFolder(
		p.projName,
		[]*folder.TemplateFile{tmplFile},
		"files/root",
		&f,
	); err != nil {
		return err
	}

	p.rootFolder, err = os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working director for project root: %w", err)
	}

	return nil
}

func (p *CProject) createSrc() error {
	if !*p.exe && !*p.lib {
		if err := folder.CreateFolder("src", nil, "", nil); err != nil {
			return err
		}

		return nil
	}

	templates := make([]*folder.TemplateFile, 0, 2)

	srcData := &SrcData{
		Exe:      *p.exe,
		ExeName:  *p.exeName,
		Flags:    p.exeFlags,
		Standard: *p.std,
		Lib:      *p.lib,
		LibName:  *p.libName,
	}

	cmakeTmpl, err := template.New("Exe CMake CMakeLists.txt").Parse(SrcCMakeListTempl)
	if err != nil {
		return fmt.Errorf("failed to create CMakeLists.txt template for project folder: %w", err)
	}

	cmakeTmplFile := &folder.TemplateFile{
		Filename: "CMakeLists.txt",
		Tmpl:     cmakeTmpl,
		TmplData: srcData,
	}

	templates = append(templates, cmakeTmplFile)

	if *p.exe {
		mainTmpl, err := template.New("Main.c").Parse(MainDotCTempl)
		if err != nil {
			return fmt.Errorf("failed to create main.c template for src folder: %w", err)
		}

		mainTmplFile := &folder.TemplateFile{
			Filename: "main.c",
			Tmpl:     mainTmpl,
			TmplData: srcData,
		}

		templates = append(templates, mainTmplFile)
	}

	if err := folder.CreateFolder("src", templates, "", nil); err != nil {
		return err
	}

	return nil
}

func (p *CProject) createLib() error {
	templates := make([]*folder.TemplateFile, 0, 3)

	libData := &LibData{
		LibName:  *p.libName,
		Flags:    p.exeFlags,
		Standard: *p.std,
	}

	cmakeTmpl, err := template.New("Lib CMake CMakeLists.txt").Parse(LibCMakeListTempl)
	if err != nil {
		return fmt.Errorf("failed to create CMakeLists.txt template for library folder: %w", err)
	}

	cmakeTmplFile := &folder.TemplateFile{
		Filename: "CMakeLists.txt",
		Tmpl:     cmakeTmpl,
		TmplData: libData,
	}

	templates = append(templates, cmakeTmplFile)

	dotCTmpl, err := template.New("Lib .C file").Parse(LibDotCTempl)
	if err != nil {
		return fmt.Errorf("failed to create .C file for library: %w", err)
	}

	dotCTmplFile := &folder.TemplateFile{
		Filename: *p.libName + ".c",
		Tmpl:     dotCTmpl,
		TmplData: libData,
	}

	templates = append(templates, dotCTmplFile)

	dotHTmpl, err := template.New("Lib .H file").Parse(LibDotHTempl)
	if err != nil {
		return fmt.Errorf("failed to create .H file for library: %w", err)
	}

	dotHTmplFile := &folder.TemplateFile{
		Filename: *p.libName + ".h",
		Tmpl:     dotHTmpl,
		TmplData: libData,
	}

	templates = append(templates, dotHTmplFile)

	if err := folder.CreateFolder(*p.libName, templates[0:2], "", nil); err != nil {
		return err
	}

	if err := folder.CreateFolder("include", templates[2:], "", nil); err != nil {
		return err
	}

	return nil
}

func (p *CProject) createTests() error {
	if err := os.Chdir(p.rootFolder); err != nil {
		return fmt.Errorf("couldn't find project root directory: %w", err)
	}

	templates := make([]*folder.TemplateFile, 0, 1)

	testData := &TestData{
		Flags:    p.exeFlags,
		Standard: *p.std,
	}

	tmpl, err := template.New("Test CMakeLists.txt").Parse(TestCMakeListTempl)
	if err != nil {
		return fmt.Errorf("failed to create CMakeLists.txt template for test folder: %w", err)
	}

	tmplFile := &folder.TemplateFile{
		Filename: "CMakeLists.txt",
		Tmpl:     tmpl,
		TmplData: testData,
	}

	templates = append(templates, tmplFile)

	if err := folder.CreateFolder("tests", templates, "files/tests", &f); err != nil {
		return err
	}

	return nil
}

func NewCProject(args []string) *CProject {
	var project CProject

	project.flagSet = flag.NewFlagSet("c", flag.ExitOnError)

	project.std = project.flagSet.String("std", stdDefault, stdDescrip)
	project.compiler = project.flagSet.String("compiler", compilerDefault, compilerDescrip)
	project.cmakeVer = project.flagSet.Float64("min-cmake", minCmakeDefault, minCmakeDescrip)
	project.lib = project.flagSet.Bool("lib", libDefault, libDescrip)
	project.libName = project.flagSet.String("lib-name", libNameDefault, libNameDescrip)
	project.exe = project.flagSet.Bool("exe", exeDefault, exeDescrip)
	project.exeName = project.flagSet.String("exe-name", exeNameDefault, exeNameDescrip)
	project.flagSet.Var(&StringSliceValue{project.exeFlags}, "exe-flags", exeFlagsDescrip)

	var help bool
	project.flagSet.BoolVar(&help, "help", helpDefault, helpDescrip)
	project.flagSet.BoolVar(&help, "h", helpDefault, helpDescrip)

	project.flagSet.Parse(args)

	if help {
		fmt.Println(ccmdHelp)
		os.Exit(0)
	}

	project.projName = args[len(args)-1]

	if *(project.lib) && *(project.libName) == "" {
		libName := "lib" + project.projName
		project.libName = &libName
	}

	if *(project.exe) && *(project.exeName) == "" {
		project.exeName = &project.projName
	}

	if len(project.exeFlags) == 0 {
		project.exeFlags = exeFlagsDefault
	}

	return &project
}

type StringSliceValue struct {
	strs []string
}

func (v StringSliceValue) String() string {
	return fmt.Sprint(v.strs)
}

func (v StringSliceValue) Set(s string) error {
	if len(s) == 0 {
		return nil
	}

	v.strs = strings.Split(s, ",")

	return nil
}
