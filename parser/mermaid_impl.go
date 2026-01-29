package parser

import (
	"fmt"
	"os"

	"github.com/sammcj/go-mermaid"
	"github.com/sammcj/go-mermaid/ast"
)

// ノードを意味する構造体
type Node struct {
	ID      string
	Name    string
	Comment string
}

// 関係を意味する構造体
type Relation struct {
	From string
	To   string
}

// マーメイドから解釈したDAGを表現する構造体
type MermaidImpl struct {
	Nodes     []Node
	Relations []Relation
}

// NewMermaidImplは、マーメイドファイルをパースして、DAGを表現する構造体を返す
func NewMermaidImpl(filePath string) (*MermaidImpl, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	parsedData, err := mermaid.Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse mermaid content: %w", err)
	}

	impl := &MermaidImpl{
		Nodes:     make([]Node, 0),
		Relations: make([]Relation, 0),
	}

	tmpNodeMap := make(map[string]string)    // ID -> Label
	tmpCommentMap := make(map[string]string) // ID -> Comment

	if flowchart, ok := parsedData.(*ast.Flowchart); ok {
		processStatements(flowchart.Statements, impl, tmpNodeMap, tmpCommentMap)
	} else {
		return nil, fmt.Errorf("failed to parse mermaid content: %w", err)
	}

	for id, label := range tmpNodeMap {
		impl.Nodes = append(impl.Nodes, Node{
			ID:      id,
			Name:    label,
			Comment: tmpCommentMap[id],
		})
	}

	return impl, nil
}

// []ast.Statementを解釈して、Nodeの情報を取り出すメソッド
func processStatements(stmts []ast.Statement, impl *MermaidImpl, nodeMap map[string]string, commentMap map[string]string) {
	var currentComment string
	for _, stmt := range stmts {

		switch v := stmt.(type) {
		case *ast.NodeDef:
			// ノードがすでに存在しても上書きする
			nodeMap[v.ID] = v.Label
			if currentComment != "" {
				commentMap[v.ID] = currentComment
			}
			currentComment = ""
		case *ast.Link:
			// Add relation
			impl.Relations = append(impl.Relations, Relation{
				From: v.From,
				To:   v.To,
			})

			// Linkにおけるノードがノードの一時マップに存在しなければ、一時マップに追加する
			if _, exists := nodeMap[v.From]; !exists {
				nodeMap[v.From] = ""

			}
			// Linkにおけるノードがノードの一時マップに存在しなければ、一時マップに追加する
			if _, exists := nodeMap[v.To]; !exists {
				nodeMap[v.To] = ""
			}
			currentComment = ""
		case *ast.Subgraph:
			processStatements(v.Statements, impl, nodeMap, commentMap)
			currentComment = ""
		case *ast.Comment:
			currentComment = v.Text
		}
	}
}

// デバッグ出力のためのプリントメソッド
func (m *MermaidImpl) Print() {
	fmt.Println("--- Nodes ---")
	for _, node := range m.Nodes {
		label := node.Name
		if label == "" {
			label = "(no label)"
		}
		fmt.Printf("ID: %s, Label: %s, Comment: %s\n", node.ID, label, node.Comment)
	}
	fmt.Println("\n--- Relations ---")
	for _, rel := range m.Relations {
		fmt.Printf("%s --> %s\n", rel.From, rel.To)
	}
}
