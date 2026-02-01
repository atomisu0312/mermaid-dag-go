package dag

import (
	"mermaid-dag-go/parser"
	"testing"
)

func TestNewMyDagImplAndRun(t *testing.T) {
	// テスト用のデータを作成
	nodes := []parser.Node{
		{ID: "A", Name: "Task A", Comment: "exec echo 'Task A done'"},
		{ID: "B", Name: "Task B", Comment: ""}, // default output
		{ID: "C", Name: "Task C", Comment: "exec echo 'Task C done'"},
	}
	relations := []parser.Relation{
		{From: "A", To: "B"},
		{From: "B", To: "C"},
	}

	mermaidGraph := &parser.MermaidImpl{
		Nodes:     nodes,
		Relations: relations,
	}

	// NewMyDagImplのテスト
	myDag, err := NewMyDagImpl(mermaidGraph)
	if err != nil {
		t.Fatalf("NewMyDagImpl failed: %v", err)
	}
	if myDag == nil {
		t.Fatal("myDag should not be nil")
	}

	// Runメソッドのテスト
	// ここでは実際にコマンドが実行される（echoコマンド）
	if err := myDag.Run(); err != nil {
		t.Errorf("Run() failed: %v", err)
	}
}
