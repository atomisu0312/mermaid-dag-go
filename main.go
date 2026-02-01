package main

import (
	"fmt"
	"mermaid-dag-go/dag"
	"mermaid-dag-go/parser"
	"os"
)

func main() {
	// 引数チェック
	if len(os.Args) < 1 {
		fmt.Printf("Usage: %s <filename>\n", os.Args[0])
		return
	}
	fileName := os.Args[1]

	// mermaidファイルをパースして、DAGを表現する構造体を返す
	mermaidGraph, err := parser.NewMermaidImpl(fileName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// natessilva/dagをラップする構造体
	myDag, err := dag.NewMyDagImpl(mermaidGraph)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Dagの実行
	if err := myDag.Run(); err != nil {
		fmt.Printf("Execution failed: %v\n", err)
		return
	}
}
