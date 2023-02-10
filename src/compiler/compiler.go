package compiler

import (
	"os"
	"path"
	"scratch-llvm/globalutil"
)

type Compiler struct {
	compilingScratch  bool
	currentOutputFile *os.File
	currentInputPath  string
}

func NewCompiler(compilingScratch bool, scratchPath, golangPath string) (*Compiler, error) {
	inputFilePath := scratchPath
	outputFilePath := golangPath
	if compilingScratch {
		inputFilePath = scratchPath
		outputFilePath = golangPath
	}

	if outputFilePath == "" {
		if compilingScratch {
			outputFilePath = path.Join("./output", "output.go")
		} else {
			outputFilePath = path.Join("./output", "output.sb3")
		}
	}

	outputFile, err := globalutil.CreateDirAndFile(outputFilePath)
	if err != nil {
		return nil, globalutil.Errorf(err, "Failed to create output file %q", outputFilePath)
	}

	outputFile.WriteString("package main\n\nfunc main() {\n\n}\n")

	return &Compiler{
		compilingScratch:  compilingScratch,
		currentOutputFile: outputFile,
		currentInputPath:  inputFilePath,
	}, nil
}

func (c *Compiler) Compile() error {
	cwd, err := os.Getwd()
	if err != nil {
		return globalutil.Errorf(err, "Failed to get current working directory")
	}

	if c.compilingScratch {
		outDir := path.Join(cwd, "uncompressed")

		err = globalutil.Unzip(c.currentInputPath, outDir)
		if err != nil {
			return globalutil.Errorf(err, "Failed to unzip file %q", c.currentInputPath)
		}
	}

	return nil
}
