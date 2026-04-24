package cxxproject

import "fmt"

const (
	stdDefault      string  = "17"
	stdDescrip      string  = "C++ standard to use (e.g., '11', '14', '17', '20', '23'; default is 17)"
	compilerDefault string  = "clang++"
	compilerDescrip string  = "compiler to use (default is clang++)"
	minCmakeDefault float64 = 3.25
	minCmakeDescrip string  = "minimum cmake version required (default is 3.25)"
	libDefault      bool    = true
	libDescrip      string  = "include a library in the project (default = true)"
	libNameDefault  string  = ""
	libNameDescrip  string  = "library name, if included. Defaults to 'lib'+project name"
	libTypeDefault  string  = ""
	libTypeDescrip  string  = "library type: 'static', 'shared', or '' to let CMake decide (default is '')"
	headerOnlyDefault bool  = false
	headerOnlyDescrip string = "make the library header-only (INTERFACE target, no .cpp source file)"
	exeDefault      bool    = true
	exeDescrip      string  = "include an executable in the project (default = true)"
	exeNameDefault  string  = ""
	exeNameDescrip  string  = "executable program's name, if included. Defaults to project name"
	exeFlagsDescrip string  = "compiler flags for targets (comma-separated)"
	helpDefault     bool    = false
	helpDescrip     string  = "display command usage and arguments"
)

var exeFlagsDefault []string = []string{"-Wall", "-Werror", "-pedantic"}

var cxxcmdHelp string = fmt.Sprintf(
	"Scaffold a C++ project.\n\n"+
		"Usage:\n\n"+
		"\tctempl cxx [arguments] <project name>\n\n"+
		"Arguments:\n\n"+
		"\t--std\t\t%s\n"+
		"\t--compiler\t%s\n"+
		"\t--min-cmake\t%s\n"+
		"\t--lib\t\t%s\n"+
		"\t--lib-name\t%s\n"+
		"\t--lib-type\t%s\n"+
		"\t--header-only\t%s\n"+
		"\t--exe\t\t%s\n"+
		"\t--exe-name\t%s\n"+
		"\t-h, --help\t%s\n",
	stdDescrip,
	compilerDescrip,
	minCmakeDescrip,
	libDescrip,
	libNameDescrip,
	libTypeDescrip,
	headerOnlyDescrip,
	exeDescrip,
	exeNameDescrip,
	helpDescrip,
)
