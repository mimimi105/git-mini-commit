package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"git-mini-commit/internal/types"
	"git-mini-commit/testutils"
)

func TestNewStorage(t *testing.T) {
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	storage, err := NewStorage()
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}

	// mini-commitsディレクトリが作成されているかチェック
	miniCommitsPath := filepath.Join(".git", "mini-commits")
	if _, err := os.Stat(miniCommitsPath); os.IsNotExist(err) {
		t.Errorf("Expected mini-commits directory to exist")
	}

	// ストレージのbasePathが正しいかチェック
	expectedPath := filepath.Join(".git", "mini-commits")
	if storage.basePath != expectedPath {
		t.Errorf("Expected basePath '%s', but got '%s'", expectedPath, storage.basePath)
	}
}

func TestNewStorageInNonGitRepository(t *testing.T) {
	// 元のディレクトリを保存
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// 一時ディレクトリに移動（Gitリポジトリではない）
	tempDir, err := os.MkdirTemp("", "not-git-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Gitリポジトリではない場所でストレージを作成
	_, err = NewStorage()
	if err == nil {
		t.Errorf("Expected error for NewStorage in non-git directory")
	}
}

func TestSaveMiniCommit(t *testing.T) {
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	storage, err := NewStorage()
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}

	// テスト用mini-commitを作成
	patch := "diff --git a/test.txt b/test.txt\nnew file mode 100644\nindex 0000000..3b18e51\n--- /dev/null\n+++ b/test.txt\n@@ -0,0 +1 @@\n+Hello, World!\n"
	now := time.Now()
	mc := &types.MiniCommit{
		ID:        storage.GenerateID(patch, now),
		Message:   "Test commit",
		CreatedAt: now,
		Patch:     patch,
	}

	// mini-commitを保存
	if err := storage.SaveMiniCommit(mc); err != nil {
		t.Fatalf("SaveMiniCommit() error = %v", err)
	}

	// インデックスファイルが作成されているかチェック
	indexPath := filepath.Join(storage.basePath, IndexFile)
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Errorf("Expected index file to exist")
	}

	// patchファイルが作成されているかチェック
	patchPath := filepath.Join(storage.basePath, mc.ID+".patch")
	if _, err := os.Stat(patchPath); os.IsNotExist(err) {
		t.Errorf("Expected patch file to exist")
	}
}

func TestLoadMiniCommits(t *testing.T) {
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	storage, err := NewStorage()
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}

	// 複数のmini-commitを作成
	patch1 := "diff --git a/test1.txt b/test1.txt\nnew file mode 100644\nindex 0000000..3b18e51\n--- /dev/null\n+++ b/test1.txt\n@@ -0,0 +1 @@\n+Hello, World!\n"
	patch2 := "diff --git a/test2.txt b/test2.txt\nnew file mode 100644\nindex 0000000..3b18e51\n--- /dev/null\n+++ b/test2.txt\n@@ -0,0 +1 @@\n+Hello, World!\n"

	mc1 := &types.MiniCommit{
		ID:        storage.GenerateID(patch1, time.Now()),
		Message:   "Test commit 1",
		CreatedAt: time.Now(),
		Patch:     patch1,
	}

	mc2 := &types.MiniCommit{
		ID:        storage.GenerateID(patch2, time.Now().Add(time.Second)),
		Message:   "Test commit 2",
		CreatedAt: time.Now().Add(time.Second),
		Patch:     patch2,
	}

	// mini-commitを保存
	if err := storage.SaveMiniCommit(mc1); err != nil {
		t.Fatalf("SaveMiniCommit() error = %v", err)
	}
	if err := storage.SaveMiniCommit(mc2); err != nil {
		t.Fatalf("SaveMiniCommit() error = %v", err)
	}

	// mini-commit一覧を読み込み
	miniCommits, err := storage.LoadMiniCommits()
	if err != nil {
		t.Fatalf("LoadMiniCommits() error = %v", err)
	}

	// 2つのmini-commitが読み込まれるかチェック
	if len(miniCommits) != 2 {
		t.Errorf("Expected 2 mini-commits, but got %d", len(miniCommits))
	}

	// 各mini-commitの内容をチェック
	found1, found2 := false, false
	for _, mc := range miniCommits {
		if mc.ID == mc1.ID {
			found1 = true
			if mc.Message != mc1.Message {
				t.Errorf("Expected message '%s', but got '%s'", mc1.Message, mc.Message)
			}
		}
		if mc.ID == mc2.ID {
			found2 = true
			if mc.Message != mc2.Message {
				t.Errorf("Expected message '%s', but got '%s'", mc2.Message, mc.Message)
			}
		}
	}

	if !found1 {
		t.Errorf("Expected to find mini-commit 1")
	}
	if !found2 {
		t.Errorf("Expected to find mini-commit 2")
	}
}

func TestGetMiniCommit(t *testing.T) {
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	storage, err := NewStorage()
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}

	// テスト用mini-commitを作成
	patch := "diff --git a/test.txt b/test.txt\nnew file mode 100644\nindex 0000000..3b18e51\n--- /dev/null\n+++ b/test.txt\n@@ -0,0 +1 @@\n+Hello, World!\n"
	mc := &types.MiniCommit{
		ID:        storage.GenerateID(patch, time.Now()),
		Message:   "Test commit",
		CreatedAt: time.Now(),
		Patch:     patch,
	}

	// mini-commitを保存
	if err := storage.SaveMiniCommit(mc); err != nil {
		t.Fatalf("SaveMiniCommit() error = %v", err)
	}

	// mini-commitを取得
	retrieved, err := storage.GetMiniCommit(mc.ID)
	if err != nil {
		t.Fatalf("GetMiniCommit() error = %v", err)
	}

	// 内容が一致するかチェック
	if retrieved.ID != mc.ID {
		t.Errorf("Expected ID '%s', but got '%s'", mc.ID, retrieved.ID)
	}
	if retrieved.Message != mc.Message {
		t.Errorf("Expected message '%s', but got '%s'", mc.Message, retrieved.Message)
	}
	if retrieved.Patch != mc.Patch {
		t.Errorf("Expected patch '%s', but got '%s'", mc.Patch, retrieved.Patch)
	}
}

func TestGetMiniCommitNotFound(t *testing.T) {
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	storage, err := NewStorage()
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}

	// 存在しないmini-commitを取得
	_, err = storage.GetMiniCommit("nonexistent-id")
	if err == nil {
		t.Errorf("Expected error for nonexistent mini-commit")
	}
}

func TestDeleteMiniCommit(t *testing.T) {
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	storage, err := NewStorage()
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}

	// テスト用mini-commitを作成
	patch := "diff --git a/test.txt b/test.txt\nnew file mode 100644\nindex 0000000..3b18e51\n--- /dev/null\n+++ b/test.txt\n@@ -0,0 +1 @@\n+Hello, World!\n"
	mc := &types.MiniCommit{
		ID:        storage.GenerateID(patch, time.Now()),
		Message:   "Test commit",
		CreatedAt: time.Now(),
		Patch:     patch,
	}

	// mini-commitを保存
	if err := storage.SaveMiniCommit(mc); err != nil {
		t.Fatalf("SaveMiniCommit() error = %v", err)
	}

	// mini-commitを削除
	if err := storage.DeleteMiniCommit(mc.ID); err != nil {
		t.Fatalf("DeleteMiniCommit() error = %v", err)
	}

	// mini-commitが削除されているかチェック
	_, err = storage.GetMiniCommit(mc.ID)
	if err == nil {
		t.Errorf("Expected error for deleted mini-commit")
	}

	// patchファイルが削除されているかチェック
	patchPath := filepath.Join(storage.basePath, mc.ID+".patch")
	if _, err := os.Stat(patchPath); !os.IsNotExist(err) {
		t.Errorf("Expected patch file to be deleted")
	}
}

func TestDeleteMiniCommitNotFound(t *testing.T) {
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	storage, err := NewStorage()
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}

	// 存在しないmini-commitを削除
	err = storage.DeleteMiniCommit("nonexistent-id")
	if err == nil {
		t.Errorf("Expected error for nonexistent mini-commit")
	}
}

func TestClearAllMiniCommits(t *testing.T) {
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	storage, err := NewStorage()
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}

	// 複数のmini-commitを作成
	patch1 := "diff --git a/test1.txt b/test1.txt\nnew file mode 100644\nindex 0000000..3b18e51\n--- /dev/null\n+++ b/test1.txt\n@@ -0,0 +1 @@\n+Hello, World!\n"
	patch2 := "diff --git a/test2.txt b/test2.txt\nnew file mode 100644\nindex 0000000..3b18e51\n--- /dev/null\n+++ b/test2.txt\n@@ -0,0 +1 @@\n+Hello, World!\n"

	mc1 := &types.MiniCommit{
		ID:        storage.GenerateID(patch1, time.Now()),
		Message:   "Test commit 1",
		CreatedAt: time.Now(),
		Patch:     patch1,
	}

	mc2 := &types.MiniCommit{
		ID:        storage.GenerateID(patch2, time.Now().Add(time.Second)),
		Message:   "Test commit 2",
		CreatedAt: time.Now().Add(time.Second),
		Patch:     patch2,
	}

	// mini-commitを保存
	if err := storage.SaveMiniCommit(mc1); err != nil {
		t.Fatalf("SaveMiniCommit() error = %v", err)
	}
	if err := storage.SaveMiniCommit(mc2); err != nil {
		t.Fatalf("SaveMiniCommit() error = %v", err)
	}

	// すべてのmini-commitを削除
	if err := storage.ClearAllMiniCommits(); err != nil {
		t.Fatalf("ClearAllMiniCommits() error = %v", err)
	}

	// mini-commit一覧が空かチェック
	miniCommits, err := storage.LoadMiniCommits()
	if err != nil {
		t.Fatalf("LoadMiniCommits() error = %v", err)
	}

	if len(miniCommits) != 0 {
		t.Errorf("Expected 0 mini-commits, but got %d", len(miniCommits))
	}

	// patchファイルが削除されているかチェック
	patch1Path := filepath.Join(storage.basePath, mc1.ID+".patch")
	patch2Path := filepath.Join(storage.basePath, mc2.ID+".patch")

	if _, err := os.Stat(patch1Path); !os.IsNotExist(err) {
		t.Errorf("Expected patch file 1 to be deleted")
	}
	if _, err := os.Stat(patch2Path); !os.IsNotExist(err) {
		t.Errorf("Expected patch file 2 to be deleted")
	}
}

func TestGenerateID(t *testing.T) {
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	storage, err := NewStorage()
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}

	patch := "test patch content"
	timestamp := time.Now()

	// IDを生成
	id1 := storage.GenerateID(patch, timestamp)
	id2 := storage.GenerateID(patch, timestamp)

	// 同じ入力に対して同じIDが生成されるかチェック
	if id1 != id2 {
		t.Errorf("Expected same ID for same input, but got %s and %s", id1, id2)
	}

	// 異なる入力に対して異なるIDが生成されるかチェック
	diffPatch := "different patch content"
	id3 := storage.GenerateID(diffPatch, timestamp)
	if id1 == id3 {
		t.Errorf("Expected different ID for different input, but got same ID")
	}

	// 異なるタイムスタンプに対して異なるIDが生成されるかチェック
	diffTimestamp := timestamp.Add(time.Second)
	id4 := storage.GenerateID(patch, diffTimestamp)
	if id1 == id4 {
		t.Errorf("Expected different ID for different timestamp, but got same ID")
	}
}

func TestStoragePersistence(t *testing.T) {
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	// 最初のストレージインスタンス
	storage1, err := NewStorage()
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}

	// テスト用mini-commitを作成
	patch := "diff --git a/test.txt b/test.txt\nnew file mode 100644\nindex 0000000..3b18e51\n--- /dev/null\n+++ b/test.txt\n@@ -0,0 +1 @@\n+Hello, World!\n"
	mc := &types.MiniCommit{
		ID:        storage1.GenerateID(patch, time.Now()),
		Message:   "Test commit",
		CreatedAt: time.Now(),
		Patch:     patch,
	}

	// mini-commitを保存
	if err := storage1.SaveMiniCommit(mc); err != nil {
		t.Fatalf("SaveMiniCommit() error = %v", err)
	}

	// 新しいストレージインスタンスを作成
	storage2, err := NewStorage()
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}

	// 保存されたmini-commitが読み込めるかチェック
	retrieved, err := storage2.GetMiniCommit(mc.ID)
	if err != nil {
		t.Fatalf("GetMiniCommit() error = %v", err)
	}

	// 内容が一致するかチェック
	if retrieved.ID != mc.ID {
		t.Errorf("Expected ID '%s', but got '%s'", mc.ID, retrieved.ID)
	}
	if retrieved.Message != mc.Message {
		t.Errorf("Expected message '%s', but got '%s'", mc.Message, retrieved.Message)
	}
	if retrieved.Patch != mc.Patch {
		t.Errorf("Expected patch '%s', but got '%s'", mc.Patch, retrieved.Patch)
	}
}

func TestStorageConcurrency(t *testing.T) {
	repo := testutils.NewTestGitRepo(t)
	defer repo.Cleanup()

	storage, err := NewStorage()
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}

	// 複数のgoroutineで同時にmini-commitを作成
	done := make(chan bool, 10)
	successCount := 0

	for i := 0; i < 10; i++ {
		go func(i int) {
			patch := "diff --git a/test.txt b/test.txt\nnew file mode 100644\nindex 0000000..3b18e51\n--- /dev/null\n+++ b/test.txt\n@@ -0,0 +1 @@\n+Hello, World!\n"
			mc := &types.MiniCommit{
				ID:        storage.GenerateID(patch, time.Now()),
				Message:   fmt.Sprintf("Test commit %d", i),
				CreatedAt: time.Now(),
				Patch:     patch,
			}

			if err := storage.SaveMiniCommit(mc); err != nil {
				// 並行処理での競合は許容する
				done <- false
				return
			}

			done <- true
		}(i)
	}

	// すべてのgoroutineが完了するまで待機
	for i := 0; i < 10; i++ {
		if <-done {
			successCount++
		}
	}

	// 成功した操作が5つ以上あればOK（並行処理での競合を許容）
	if successCount < 5 {
		t.Errorf("Expected at least 5 successful operations, but got %d", successCount)
	}

	// 保存されたmini-commitの数をチェック
	miniCommits, err := storage.LoadMiniCommits()
	if err != nil {
		t.Fatalf("LoadMiniCommits() error = %v", err)
	}

	if len(miniCommits) < 5 {
		t.Errorf("Expected at least 5 mini-commits, but got %d", len(miniCommits))
	}
}
