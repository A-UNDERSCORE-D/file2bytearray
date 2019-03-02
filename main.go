package main

import (
    "bytes"
    "flag"
    "fmt"
    "go/format"
    "io/ioutil"
    "os"
    "strings"
)

func main() {
    target := flag.String("target", "", "sets the target file to dump")
    outName := flag.String("out", "-", "sets the output file path, - will use stdout")
    packageName := flag.String("package", "out", "set the package name for the outputted file")
    varName := flag.String("var", "Data", "set the var name for the outputted file")
    flag.Parse()
    if *target == "" {
        _, _ = fmt.Fprintln(os.Stderr, "a target file is required")
        return
    }
    data, err := readFile(*target)
    if err != nil {
        _, _ = fmt.Fprintf(os.Stderr, "could not read file: %s\n", err)
    }
    out := bytes.Buffer{}
    out.WriteString(fmt.Sprintf("package %s\n", *packageName))
    out.WriteString(fmt.Sprintf("var %s = %s", *varName, pprintSlice(data)))
    outSlice, err := format.Source(out.Bytes())
    if err != nil {
        _, _ = fmt.Fprintf(os.Stderr, "could not format output: %s\n", err)
        fmt.Println(out.String())
        return
    }
    if *outName == "-" {
        fmt.Println(string(outSlice))
    } else {
        f, err := os.Create(*outName)
        if err != nil {
            _, _ = fmt.Fprintf(os.Stderr, "could not open output file: %s\n", err)
            return
        }
        defer f.Close()
        if _, err = f.Write(outSlice); err != nil {
            _, _ = fmt.Fprintf(os.Stdout, "could not write to file: %s\n", err)
            return
        }

    }
}

func readFile(name string) ([]byte, error) {
    f, err := os.Open(name)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    return ioutil.ReadAll(f)
}

func pprintSlice(s []byte) string {
    out := strings.Builder{}
    if len(s) < 80 {
        out.WriteString(fmt.Sprintf("%#v", s))
        return out.String()
    }
    w := 0
    out.WriteString("[]byte{\n")
    for _, v := range s {
        if w == 20 {
            out.WriteRune('\n')
            w = 0
        }
        out.WriteString("0x")
        tw := fmt.Sprintf("%X", v)
        if len(tw) == 1 {
            out.WriteString("0")
        }
        out.WriteString(tw)
        out.WriteString(", ")
        w++
    }
    out.WriteString("\n}")
    return out.String()
}
