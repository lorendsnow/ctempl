package cxxproject

type ProjData struct {
	MinVersion  string
	Compiler    string
	ProjectName string
}

var projCMakeListTempl string = "cmake_minimum_required(VERSION {{.MinVersion}})\n\n" +
	"set(CMAKE_EXPORT_COMPILE_COMMANDS ON)\n" +
	"set(CMAKE_CXX_COMPILER {{.Compiler}})\n\n" +
	"enable_testing()\n\n" +
	"project({{.ProjectName}})\n\n" +
	"add_subdirectory(src bin)\n" +
	"add_subdirectory(tests)"

type SrcData struct {
	Exe        bool
	ExeName    string
	Flags      []string
	Standard   string
	Lib        bool
	LibName    string
	HeaderOnly bool
}

var SrcCMakeListTempl string = "{{ if .Lib }}add_subdirectory({{.LibName}}){{ end }}" +
	"{{ if .Exe }}\n\nadd_executable({{.ExeName}} main.cpp)\n\n" +
	"set_target_properties({{.ExeName}} PROPERTIES CXX_STANDARD {{.Standard}})\n\n" +
	"target_compile_options({{.ExeName}}\n" +
	"\tPRIVATE\n" +
	"{{ range .Flags }}" +
	"\t\t{{.}}\n" +
	"{{ end }}" +
	")" +
	"{{ if .Lib }}\n\ntarget_link_libraries({{.ExeName}} PRIVATE {{.LibName}}){{ end }}{{ end }}"

var MainDotCppTempl string = "{{ if .Lib }}#include \"{{.LibName}}.hpp\"\n\n" +
	"{{ else }}#include <cstdlib>\n\n{{ end }}" +
	"int main(int argc, char* argv[]) { " +
	"return {{ if .Lib }}lib_func(){{ else }}EXIT_SUCCESS{{ end }}; }"

type LibData struct {
	LibName    string
	LibType    string
	Flags      []string
	Standard   string
	HeaderOnly bool
}

// LibCMakeListTempl handles both regular and header-only (INTERFACE) libraries.
var LibCMakeListTempl string = "{{ if .HeaderOnly }}" +
	"add_library({{.LibName}} INTERFACE)\n\n" +
	"target_include_directories({{.LibName}} INTERFACE include)" +
	"{{ else }}" +
	"add_library({{.LibName}}{{ if .LibType }} {{.LibType}}{{ end }} {{.LibName}}.cpp)\n\n" +
	"target_include_directories({{.LibName}} PUBLIC include)\n\n" +
	"set_target_properties({{.LibName}} PROPERTIES CXX_STANDARD {{.Standard}})\n\n" +
	"target_compile_options({{.LibName}}\n" +
	"\tPRIVATE\n" +
	"{{ range .Flags }}" +
	"\t\t{{.}}\n" +
	"{{ end }}" +
	")" +
	"{{ end }}"

var LibDotCppTempl string = "#include \"{{.LibName}}.hpp\"\n\n" +
	"int lib_func() { return 1; }"

var LibDotHppTempl string = "#pragma once\n\n" +
	"int lib_func();\n"

var LibDotHppImplTempl string = "#pragma once\n\n" +
	"inline int lib_func() {\n\treturn 1;\n}\n"

type TestData struct {
	Flags    []string
	Standard string
}

var TestCMakeListTempl string = "find_package(GTest REQUIRED)\n\n" +
	"add_executable(tests tests.cpp)\n\n" +
	"set_target_properties(tests PROPERTIES CXX_STANDARD {{.Standard}})\n\n" +
	"target_link_libraries(tests PRIVATE GTest::gtest_main)\n\n" +
	"target_compile_options(tests\n" +
	"\tPRIVATE\n" +
	"{{ range .Flags }}" +
	"\t\t{{.}}\n" +
	"{{ end }}" +
	")\n\n" +
	"include(GoogleTest)\n" +
	"gtest_discover_tests(tests)"
