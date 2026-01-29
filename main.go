package main

import (
	"fmt"
	"mermaid-dag-go/parser"
	"os/exec"
)

func main() {
	// mermaidファイルをパースして、DAGを表現する構造体を返す
	dag, err := parser.NewMermaidImpl("sample.mmd")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	// デバッグ出力
	dag.Print()

	// TODO: DAGファイルからパイプラインを動的に作成して実行する

	// Goからターミナルコマンドを実施するメソッド
	ls, err := exec.Command("ls").Output()
	fmt.Printf("hello ls:\n%s", ls)
}
