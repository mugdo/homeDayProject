package backEnd

type LanguagePack struct {
	LangValue string
	LangName  string
}

func getLenguage(oj string) []LanguagePack {
	var languagePack []LanguagePack

	if oj == "CodeForces" {
		languagePack = []LanguagePack{
			LanguagePack{LangValue: "14", LangName: "ActiveTcl 8.5"},
			LanguagePack{LangValue: "33", LangName: "Ada GNAT 4"},
			LanguagePack{LangValue: "18", LangName: "Befunge"},
			LanguagePack{LangValue: "52", LangName: "Clang++17 Diagnostics"},
			LanguagePack{LangValue: "9", LangName: "C# Mono 5.18"},
			LanguagePack{LangValue: "28", LangName: "D DMD32 v2.091.0"},
			LanguagePack{LangValue: "3", LangName: "Delphi 7"},
			LanguagePack{LangValue: "25", LangName: "Factor"},
			LanguagePack{LangValue: "39", LangName: "FALSE"},
			LanguagePack{LangValue: "4", LangName: "Free Pascal 3.0.2"},
			LanguagePack{LangValue: "43", LangName: "GNU GCC C11 5.1.0"},
			LanguagePack{LangValue: "45", LangName: "GNU C++11 5 ZIP"},
			LanguagePack{LangValue: "42", LangName: "GNU G++11 5.1.0"},
			LanguagePack{LangValue: "50", LangName: "GNU G++14 6.4.0"},
			LanguagePack{LangValue: "54", LangName: "GNU G++17 7.3.0"},
			LanguagePack{LangValue: "61", LangName: "GNU G++17 9.2.0 (64 bit, msys 2)"},
			LanguagePack{LangValue: "32", LangName: "Go 1.14"},
			LanguagePack{LangValue: "12", LangName: "Haskell GHC 8.6.3"},
			LanguagePack{LangValue: "15", LangName: "Io-2008-01-07 (Win32)"},
			LanguagePack{LangValue: "47", LangName: "J"},
			LanguagePack{LangValue: "36", LangName: "Java 1.8.0_162"},
			LanguagePack{LangValue: "60", LangName: "Java 11.0.5"},
			LanguagePack{LangValue: "46", LangName: "Java 8 ZIP"},
			LanguagePack{LangValue: "34", LangName: "JavaScript V8 4.8.0"},
			LanguagePack{LangValue: "48", LangName: "Kotlin 1.3.70"},
			LanguagePack{LangValue: "56", LangName: "Microsoft Q#"},
			LanguagePack{LangValue: "2", LangName: "Microsoft Visual C++ 2010"},
			LanguagePack{LangValue: "59", LangName: "Microsoft Visual C++ 2017"},
			LanguagePack{LangValue: "38", LangName: "Mysterious Language"},
			LanguagePack{LangValue: "55", LangName: "Node.js 9.4.0"},
			LanguagePack{LangValue: "19", LangName: "OCaml 4.02.1"},
			LanguagePack{LangValue: "22", LangName: "OpenCobol 1.0"},
			LanguagePack{LangValue: "51", LangName: "PascalABC.NET 3.4.2"},
			LanguagePack{LangValue: "13", LangName: "Perl 5.20.1"},
			LanguagePack{LangValue: "6", LangName: "PHP 7.2.13"},
			LanguagePack{LangValue: "44", LangName: "Picat 0.9"},
			LanguagePack{LangValue: "17", LangName: "Pike 7.8"},
			LanguagePack{LangValue: "40", LangName: "PyPy 2.7 (7.2.0)"},
			LanguagePack{LangValue: "41", LangName: "PyPy 3.6 (7.2.0)"},
			LanguagePack{LangValue: "7", LangName: "Python 2.7.15"},
			LanguagePack{LangValue: "31", LangName: "Python 3.7.2"},
			LanguagePack{LangValue: "27", LangName: "Roco"},
			LanguagePack{LangValue: "8", LangName: "Ruby 2.0.0p645"},
			LanguagePack{LangValue: "49", LangName: "Rust 1.42.0"},
			LanguagePack{LangValue: "20", LangName: "Scala 2.12.8"},
			LanguagePack{LangValue: "26", LangName: "Secret_171"},
			LanguagePack{LangValue: "57", LangName: "Text"},
		}
	} else if oj == "LightOJ" {
		languagePack = []LanguagePack{
			LanguagePack{LangValue: "C", LangName: "C"},
			LanguagePack{LangValue: "C++", LangName: "C++"},
			LanguagePack{LangValue: "JAVA", LangName: "JAVA"},
			LanguagePack{LangValue: "PASCAL", LangName: "PASCAL"},
		}
	}

	return languagePack
}
