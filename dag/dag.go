package dag

import (
	"fmt"
	"mermaid-dag-go/parser"
	"os"
	"os/exec"
	"strings"

	thedag "github.com/natessilva/dag"
)

// MyDagImpl : natessilva/dagをラップする構造体
type MyDagImpl struct {
	Dag thedag.Runner
}

// NewMyDagImpl : mermaidの構造体を引数に取って、DagmermaidImplを返す
func NewMyDagImpl(NewMermaidImpl *parser.MermaidImpl) (*MyDagImpl, error) {
	var r thedag.Runner

	for _, node := range NewMermaidImpl.Nodes {
		r.AddVertex(node.ID, func() error {
			// コメントの前後の空白を除去
			comment := strings.TrimSpace(node.Comment)

			// コメントが exec で始まる場合はコマンドを実行する
			if strings.HasPrefix(comment, "exec ") {
				cmdStr := strings.TrimPrefix(comment, "exec ")
				cmd := exec.Command("sh", "-c", cmdStr)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				return cmd.Run()
			}

			// それ以外の場合はノード名を出力する
			fmt.Println(node.Name)
			return nil
		})
	}

	for _, relation := range NewMermaidImpl.Relations {
		r.AddEdge(relation.From, relation.To)
	}

	return &MyDagImpl{
		Dag: r,
	}, nil
}

// Run : 関数を実行するメソッド
func (m *MyDagImpl) Run() error {
	return m.Dag.Run()
}
