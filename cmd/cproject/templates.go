package cproject

// /////////////////////////////////////////////////////////////////////////////

//	                                                                          //
//						 TOP-LEVEL CMAKELISTS.TXT TEMPLATE                    //
//	                                                                          //

// /////////////////////////////////////////////////////////////////////////////
type ProjCMakeLists struct {
	MinVersion  float64
	Compiler    string
	ProjectName string
}

var projCMakeListTempl string = "cmake_minimum_required(VERSION {{.MinVersion}})\n\n" +
	"set(CMAKE_EXPORT_COMPILE_COMMANDS ON)\n" +
	"set(CMAKE_C_COMPILER {{.Compiler}})\n\n" +
	"enable_testing()\n\n" +
	"project({{.ProjectName}})\n\n" +
	"add_subdirectory(src bin)\n" +
	"add_subdirectory(tests)"

// /////////////////////////////////////////////////////////////////////////////

//	                                                                          //
//						EXECUTABLE CMAKELISTS.TXT TEMPLATE                    //
//	                                                                          //

// /////////////////////////////////////////////////////////////////////////////
type ExeCMakeLists struct {
	Exe      bool
	ExeName  string
	Flags    []string
	Standard string
	Lib      bool
	LibName  string
}

var ExeCMakeListTempl string = "{{ if .Lib }}add_subdirectory({{.LibName}}){{ end }}" +
	"{{ if .Exe }}\n\nadd_executable({{.ExeName}} main.c)\n\n" +
	"set_target_properties({{.ExeName}} PROPERTIES C_STANDARD {{.Standard}})\n\n" +
	"target_compile_options({{.ExeName}}\n" +
	"\tPRIVATE\n" +
	"{{ range .Flags }}" +
	"\t\t{{.}}\n" +
	"{{ end }}" +
	")" +
	"{{ if .Lib }}\n\ntarget_link_libraries({{.ExeName}} PRIVATE {{.LibName}}){{ end }}{{ end }}"

// /////////////////////////////////////////////////////////////////////////////

//	                                                                          //
//						        MAIN.C TEMPLATE                               //
//	                                                                          //

// /////////////////////////////////////////////////////////////////////////////
var MainDotCTempl string = "{{ if .Lib }}#include \"{{.LibName}}.h\"\n\n" +
	"{{ else }}#include <stdlib.h>\n\n{{ end }}" +
	"int main(int argc, char* argv[]) { " +
	"return {{ if .Lib }}lib_func(){{ else }}EXIT_SUCCESS{{ end }}; }"

// /////////////////////////////////////////////////////////////////////////////

//	                                                                          //
//		  		      LIBRARY CMAKELISTS.TXT TEMPLATE                         //
//	                                                                          //

// /////////////////////////////////////////////////////////////////////////////
type LibCMakeLists struct {
	LibName  string
	Flags    []string
	Standard string
}

var LibCMakeListTempl string = "add_library({{.LibName}} {{.LibName}}.c)\n\n" +
	"target_include_directories({{.LibName}} PUBLIC include)\n\n" +
	"set_target_properties({{.LibName}} PROPERTIES C_STANDARD {{.Standard}})\n\n" +
	"target_compile_options({{.LibName}}\n" +
	"\tPRIVATE\n" +
	"{{ range .Flags }}" +
	"\t\t{{.}}\n" +
	"{{ end }}" +
	")"

// /////////////////////////////////////////////////////////////////////////////

//	                                                                          //
//		  		            LIBRARY.C TEMPLATE                                //
//	                                                                          //

// /////////////////////////////////////////////////////////////////////////////
type LibDotC struct {
	LibName string
}

var LibDotCTempl string = "#include \"{{.LibName}}.h\"\n"

// /////////////////////////////////////////////////////////////////////////////

//	                                                                          //
//		  		            LIBRARY.H TEMPLATE                                //
//	                                                                          //

// /////////////////////////////////////////////////////////////////////////////
type LibDotH struct {
	LibName string
}

var LibDotHTempl string = "#ifndef {{.LibName}}_H\n" +
	"#define {{.LibName}}_H\n\n" +
	"int lib_func(void) {\n\treturn 1;\n}\n\n" +
	"#endif //{{.LibName}}_H"

// /////////////////////////////////////////////////////////////////////////////

//	                                                                          //
//		  		      TESTING CMAKELISTS.TXT TEMPLATE                         //
//	                                                                          //

// /////////////////////////////////////////////////////////////////////////////
type TestCMakeLists struct {
	Flags    []string
	Standard string
}

var TestCMakeListTempl string = "add_executable(tests tests.c)\n\n" +
	"set_target_properties(tests PROPERTIES C_STANDARD {{.Standard}})\n\n" +
	"target_link_libraries(tests PRIVATE cunit)\n\n" +
	"target_compile_options(tests\n" +
	"\tPRIVATE\n" +
	"{{ range .Flags }}" +
	"\t\t{{.}}\n" +
	"{{ end }}" +
	")\n\n" +
	"add_test(\n" +
	"\tNAME tests\n" +
	"\tCOMMAND $<TARGET_FILE:tests>\n" +
	")"
