package cproject

import "fmt"

const (
	stdDefault      string  = "99"
	stdDescrip      string  = "C standard to use (e.g., '99', '11'; default is 99)"
	compilerDefault string  = "clang"
	compilerDescrip string  = "compiler to use. Default is clang"
	minCmakeDefault float64 = 3.25
	minCmakeDescrip string  = "minimum cmake version required (default is 3.25)"
	libDefault      bool    = true
	libDescrip      string  = "include a library in the project (default = true)"
	libNameDefault  string  = ""
	libNameDescrip          = "library name, if included. Defaults to 'lib'+project name"
	exeDefault      bool    = true
	exeDescrip      string  = "include an executable in the project (default = true)"
	exeNameDefault  string  = ""
	exeNameDescrip  string  = "executable program's name, if included. Defaults to project name"
	exeFlagsDescrip string  = "compiler flags for the exe target"
	helpDefault     bool    = false
	helpDescrip     string  = "display command usage and arguments"
)

var exeFlagsDefault []string = []string{"-Wall", "-Werror", "-pedantic"}

var ccmdHelp string = fmt.Sprintf(
	"Scaffold a C project.\n\n"+
		"Usage:\n\n"+
		"\tctempl c [argmuents] <project name>\n\n"+
		"Arguments:\n\n"+
		"\t--std\t\t%s\n"+
		"\t--compiler\t%s\n"+
		"\t--min-cmake\t%s\n"+
		"\t--lib\t\t%s\n"+
		"\t--lib-name\t%s\n"+
		"\t--exe\t\t%s\n"+
		"\t--exe-name\t%s\n"+
		"\t-h, --help\t%s\n",
	stdDescrip,
	compilerDescrip,
	minCmakeDescrip,
	libDescrip,
	libNameDescrip,
	exeDescrip,
	exeNameDescrip,
	helpDescrip,
)
