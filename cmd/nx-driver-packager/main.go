package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/NotrixInc/nx-driver-packager/internal/packager"
)

func main() {
	packCmd := flag.NewFlagSet("pack", flag.ExitOnError)
	input := packCmd.String("input", "", "driver directory")
	out := packCmd.String("out", "", "output .nxpkg file")

	verifyCmd := flag.NewFlagSet("verify", flag.ExitOnError)
	pkg := verifyCmd.String("pkg", "", "path to .nxpkg file")
	verifyInput := verifyCmd.String("input", "", "optional driver directory to compare against")

	if len(os.Args) < 2 {
		fmt.Println("expected 'pack' command")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "pack":
		_ = packCmd.Parse(os.Args[2:])
		if *input == "" || *out == "" {
			fmt.Println("input and out are required")
			os.Exit(1)
		}
		if err := packager.Pack(*input, *out); err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		fmt.Println("package created:", *out)
	case "verify":
		_ = verifyCmd.Parse(os.Args[2:])
		if *pkg == "" {
			fmt.Println("pkg is required")
			os.Exit(1)
		}
		if err := packager.VerifyPackage(*pkg, *verifyInput); err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		fmt.Println("package verified:", *pkg)
	default:
		fmt.Println("unknown command")
		os.Exit(1)
	}
}
