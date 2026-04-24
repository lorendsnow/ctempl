package cxxproject

import (
	"embed"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"text/template"

	"github.com/lorendsnow/ctempl/internal/folder"
)

// formatCMakeVersion formats a float64 cmake version so that whole numbers
// always include a decimal point (e.g. 4 → "4.0", 3.25 → "3.25").
func formatCMakeVersion(v float64) string {
	if v == math.Trunc(v) {
		return fmt.Sprintf("%.1f", v)
	}
	return fmt.Sprintf("%g", v)
}

//go:embed files/root/* files/tests/*
var f embed.FS

type CXXProject struct {
	std        *string
	compiler   *string
	cmakeVer   *float64
	lib        *bool
	libName    *string
	libType    *string
	headerOnly *bool
	exe        *bool
	exeName    *string
	exeFlags   []string
	projName   string
	flagSet    *flag.FlagSet
	rootFolder string
}

func (p *CXXProject) Run() error {
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

func (p *CXXProject) createRoot() error {
	projData := &ProjData{
		MinVersion:  formatCMakeVersion(*p.cmakeVer),
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
		return fmt.Errorf("failed to get working directory for project root: %w", err)
	}

	return nil
}

func (p *CXXProject) createSrc() error {
	if !*p.exe && !*p.lib {
		if err := folder.CreateFolder("src", nil, "", nil); err != nil {
			return err
		}

		return nil
	}

	templates := make([]*folder.TemplateFile, 0, 2)

	srcData := &SrcData{
		Exe:        *p.exe,
		ExeName:    *p.exeName,
		Flags:      p.exeFlags,
		Standard:   *p.std,
		Lib:        *p.lib,
		LibName:    *p.libName,
		HeaderOnly: *p.headerOnly,
	}

	cmakeTmpl, err := template.New("Exe CMake CMakeLists.txt").Parse(SrcCMakeListTempl)
	if err != nil {
		return fmt.Errorf("failed to create CMakeLists.txt template for src folder: %w", err)
	}

	cmakeTmplFile := &folder.TemplateFile{
		Filename: "CMakeLists.txt",
		Tmpl:     cmakeTmpl,
		TmplData: srcData,
	}

	templates = append(templates, cmakeTmplFile)

	if *p.exe {
		mainTmpl, err := template.New("Main.cpp").Parse(MainDotCppTempl)
		if err != nil {
			return fmt.Errorf("failed to create main.cpp template for src folder: %w", err)
		}

		mainTmplFile := &folder.TemplateFile{
			Filename: "main.cpp",
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

func (p *CXXProject) createLib() error {
	libData := &LibData{
		LibName:    *p.libName,
		LibType:    *p.libType,
		Flags:      p.exeFlags,
		Standard:   *p.std,
		HeaderOnly: *p.headerOnly,
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

	// Choose the right .hpp template: header-only gets the inline impl,
	// regular libs get a declaration-only header paired with a .cpp.
	hppTemplStr := LibDotHppTempl
	if *p.headerOnly {
		hppTemplStr = LibDotHppImplTempl
	}

	dotHppTmpl, err := template.New("Lib .hpp file").Parse(hppTemplStr)
	if err != nil {
		return fmt.Errorf("failed to create .hpp template for library: %w", err)
	}

	dotHppTmplFile := &folder.TemplateFile{
		Filename: *p.libName + ".hpp",
		Tmpl:     dotHppTmpl,
		TmplData: libData,
	}

	if *p.headerOnly {
		// Header-only: lib dir has only CMakeLists.txt; include/ has the .hpp.
		if err := folder.CreateFolder(*p.libName, []*folder.TemplateFile{cmakeTmplFile}, "", nil); err != nil {
			return err
		}

		if err := folder.CreateFolder("include", []*folder.TemplateFile{dotHppTmplFile}, "", nil); err != nil {
			return err
		}

		return nil
	}

	// Regular lib: lib dir gets CMakeLists.txt + .cpp; include/ gets .hpp.
	dotCppTmpl, err := template.New("Lib .cpp file").Parse(LibDotCppTempl)
	if err != nil {
		return fmt.Errorf("failed to create .cpp template for library: %w", err)
	}

	dotCppTmplFile := &folder.TemplateFile{
		Filename: *p.libName + ".cpp",
		Tmpl:     dotCppTmpl,
		TmplData: libData,
	}

	if err := folder.CreateFolder(*p.libName, []*folder.TemplateFile{cmakeTmplFile, dotCppTmplFile}, "", nil); err != nil {
		return err
	}

	if err := folder.CreateFolder("include", []*folder.TemplateFile{dotHppTmplFile}, "", nil); err != nil {
		return err
	}

	return nil
}

func (p *CXXProject) createTests() error {
	if err := os.Chdir(p.rootFolder); err != nil {
		return fmt.Errorf("couldn't find project root directory: %w", err)
	}

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

	if err := folder.CreateFolder("tests", []*folder.TemplateFile{tmplFile}, "files/tests", &f); err != nil {
		return err
	}

	return nil
}

func NewCXXProject(args []string) *CXXProject {
	var project CXXProject

	project.flagSet = flag.NewFlagSet("cxx", flag.ExitOnError)

	project.std = project.flagSet.String("std", stdDefault, stdDescrip)
	project.compiler = project.flagSet.String("compiler", compilerDefault, compilerDescrip)
	project.cmakeVer = project.flagSet.Float64("min-cmake", minCmakeDefault, minCmakeDescrip)
	project.lib = project.flagSet.Bool("lib", libDefault, libDescrip)
	project.libName = project.flagSet.String("lib-name", libNameDefault, libNameDescrip)
	project.libType = project.flagSet.String("lib-type", libTypeDefault, libTypeDescrip)
	project.headerOnly = project.flagSet.Bool("header-only", headerOnlyDefault, headerOnlyDescrip)
	project.exe = project.flagSet.Bool("exe", exeDefault, exeDescrip)
	project.exeName = project.flagSet.String("exe-name", exeNameDefault, exeNameDescrip)
	project.flagSet.Var(&StringSliceValue{&project.exeFlags}, "exe-flags", exeFlagsDescrip)

	var help bool
	project.flagSet.BoolVar(&help, "help", helpDefault, helpDescrip)
	project.flagSet.BoolVar(&help, "h", helpDefault, helpDescrip)

	project.flagSet.Parse(args)

	if help {
		fmt.Println(cxxcmdHelp)
		os.Exit(0)
	}

	project.projName = args[len(args)-1]

	if *project.lib && *project.libName == "" {
		libName := "lib" + project.projName
		project.libName = &libName
	}

	if *project.exe && *project.exeName == "" {
		project.exeName = &project.projName
	}

	if len(project.exeFlags) == 0 {
		project.exeFlags = exeFlagsDefault
	}

	return &project
}

// StringSliceValue is a flag.Value that populates a []string via a pointer,
// fixing the value-receiver bug present in the C version.
type StringSliceValue struct {
	strs *[]string
}

func (v StringSliceValue) String() string {
	if v.strs == nil {
		return "[]"
	}
	return fmt.Sprint(*v.strs)
}

func (v StringSliceValue) Set(s string) error {
	if len(s) == 0 {
		return nil
	}

	*v.strs = strings.Split(s, ",")

	return nil
}
