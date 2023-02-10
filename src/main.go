package main

import (
	"os"
	"scratch-llvm/compiler"

	kingpin "github.com/alecthomas/kingpin/v2"
)

var (
	app = kingpin.New("Scratch Golang LLVM", "Compiler for Scratch to Golang and vice versa")

	toScratch = app.Command("scratch", "Compile Golang to Scratch")

	inputFile  = toScratch.Arg("input", "Input file").Required().String()
	outputFile = toScratch.Arg("output", "Output file").String()
)

func main() {

	kpCmd := kingpin.MustParse(app.Parse(os.Args[1:]))
	compilingScratch := kpCmd == toScratch.FullCommand()

	compiler, err := compiler.NewCompiler(compilingScratch, *inputFile, *outputFile)
	if err != nil {
		panic(err)
	}

	if err := compiler.Compile(); err != nil {
		panic(err)
	}
}
