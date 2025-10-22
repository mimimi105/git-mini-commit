package testutils

import (
	"fmt"
	"testing"
	"time"

	"git-mini-commit/internal/storage"
	"git-mini-commit/internal/types"
)

// BenchmarkStorageOperations ストレージ操作のベンチマーク
func BenchmarkStorageOperations(b *testing.B) {
	// テスト用Gitリポジトリを作成
	repo := NewTestGitRepo(&testing.T{})
	defer repo.Cleanup()

	// テスト用ストレージを作成
	storage, err := storage.NewStorage()
	if err != nil {
		b.Fatalf("Failed to create storage: %v", err)
	}

	// テスト用patchを作成
	patch := "diff --git a/test.txt b/test.txt\nnew file mode 100644\nindex 0000000..3b18e51\n--- /dev/null\n+++ b/test.txt\n@@ -0,0 +1 @@\n+Hello, World!\n"

	b.Run("SaveMiniCommit", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mc := &types.MiniCommit{
				ID:        storage.GenerateID(patch, time.Now()),
				Message:   "Test commit",
				CreatedAt: time.Now(),
				Patch:     patch,
			}
			if err := storage.SaveMiniCommit(mc); err != nil {
				b.Fatalf("SaveMiniCommit() error = %v", err)
			}
		}
	})

	b.Run("LoadMiniCommits", func(b *testing.B) {
		// 事前にmini-commitを作成
		for i := 0; i < 100; i++ {
			mc := &types.MiniCommit{
				ID:        storage.GenerateID(patch, time.Now()),
				Message:   "Test commit",
				CreatedAt: time.Now(),
				Patch:     patch,
			}
			if err := storage.SaveMiniCommit(mc); err != nil {
				b.Fatalf("SaveMiniCommit() error = %v", err)
			}
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := storage.LoadMiniCommits()
			if err != nil {
				b.Fatalf("LoadMiniCommits() error = %v", err)
			}
		}
	})

	b.Run("GetMiniCommit", func(b *testing.B) {
		// 事前にmini-commitを作成
		mc := &types.MiniCommit{
			ID:        storage.GenerateID(patch, time.Now()),
			Message:   "Test commit",
			CreatedAt: time.Now(),
			Patch:     patch,
		}
		if err := storage.SaveMiniCommit(mc); err != nil {
			b.Fatalf("SaveMiniCommit() error = %v", err)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := storage.GetMiniCommit(mc.ID)
			if err != nil {
				b.Fatalf("GetMiniCommit() error = %v", err)
			}
		}
	})

	b.Run("DeleteMiniCommit", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// 事前にmini-commitを作成
			mc := &types.MiniCommit{
				ID:        storage.GenerateID(patch, time.Now()),
				Message:   "Test commit",
				CreatedAt: time.Now(),
				Patch:     patch,
			}
			if err := storage.SaveMiniCommit(mc); err != nil {
				b.Fatalf("SaveMiniCommit() error = %v", err)
			}

			// mini-commitを削除
			if err := storage.DeleteMiniCommit(mc.ID); err != nil {
				b.Fatalf("DeleteMiniCommit() error = %v", err)
			}
		}
	})
}

// BenchmarkStorageConcurrency ストレージの並行性ベンチマーク
func BenchmarkStorageConcurrency(b *testing.B) {
	// テスト用Gitリポジトリを作成
	repo := NewTestGitRepo(&testing.T{})
	defer repo.Cleanup()

	// テスト用ストレージを作成
	storage, err := storage.NewStorage()
	if err != nil {
		b.Fatalf("Failed to create storage: %v", err)
	}

	// テスト用patchを作成
	patch := "diff --git a/test.txt b/test.txt\nnew file mode 100644\nindex 0000000..3b18e51\n--- /dev/null\n+++ b/test.txt\n@@ -0,0 +1 @@\n+Hello, World!\n"

	b.Run("ConcurrentSave", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				mc := &types.MiniCommit{
					ID:        storage.GenerateID(patch, time.Now()),
					Message:   "Test commit",
					CreatedAt: time.Now(),
					Patch:     patch,
				}
				if err := storage.SaveMiniCommit(mc); err != nil {
					b.Fatalf("SaveMiniCommit() error = %v", err)
				}
			}
		})
	})

	b.Run("ConcurrentLoad", func(b *testing.B) {
		// 事前にmini-commitを作成
		for i := 0; i < 100; i++ {
			mc := &types.MiniCommit{
				ID:        storage.GenerateID(patch, time.Now()),
				Message:   "Test commit",
				CreatedAt: time.Now(),
				Patch:     patch,
			}
			if err := storage.SaveMiniCommit(mc); err != nil {
				b.Fatalf("SaveMiniCommit() error = %v", err)
			}
		}

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := storage.LoadMiniCommits()
				if err != nil {
					b.Fatalf("LoadMiniCommits() error = %v", err)
				}
			}
		})
	})
}

// BenchmarkStorageMemory ストレージのメモリ使用量ベンチマーク
func BenchmarkStorageMemory(b *testing.B) {
	// テスト用Gitリポジトリを作成
	repo := NewTestGitRepo(&testing.T{})
	defer repo.Cleanup()

	// テスト用ストレージを作成
	storage, err := storage.NewStorage()
	if err != nil {
		b.Fatalf("Failed to create storage: %v", err)
	}

	// 大きなpatchを作成
	largePatch := "diff --git a/test.txt b/test.txt\nnew file mode 100644\nindex 0000000..3b18e51\n--- /dev/null\n+++ b/test.txt\n@@ -0,0 +1 @@\n+Hello, World!\n"
	for i := 0; i < 1000; i++ {
		largePatch += "This is a test line.\n"
	}

	b.Run("LargePatchSave", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mc := &types.MiniCommit{
				ID:        storage.GenerateID(largePatch, time.Now()),
				Message:   "Test commit",
				CreatedAt: time.Now(),
				Patch:     largePatch,
			}
			if err := storage.SaveMiniCommit(mc); err != nil {
				b.Fatalf("SaveMiniCommit() error = %v", err)
			}
		}
	})

	b.Run("LargePatchLoad", func(b *testing.B) {
		// 事前に大きなmini-commitを作成
		mc := &types.MiniCommit{
			ID:        storage.GenerateID(largePatch, time.Now()),
			Message:   "Test commit",
			CreatedAt: time.Now(),
			Patch:     largePatch,
		}
		if err := storage.SaveMiniCommit(mc); err != nil {
			b.Fatalf("SaveMiniCommit() error = %v", err)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := storage.GetMiniCommit(mc.ID)
			if err != nil {
				b.Fatalf("GetMiniCommit() error = %v", err)
			}
		}
	})
}

// BenchmarkStorageDiskIO ストレージのディスクI/Oベンチマーク
func BenchmarkStorageDiskIO(b *testing.B) {
	// テスト用Gitリポジトリを作成
	repo := NewTestGitRepo(&testing.T{})
	defer repo.Cleanup()

	// テスト用ストレージを作成
	storage, err := storage.NewStorage()
	if err != nil {
		b.Fatalf("Failed to create storage: %v", err)
	}

	// テスト用patchを作成
	patch := "diff --git a/test.txt b/test.txt\nnew file mode 100644\nindex 0000000..3b18e51\n--- /dev/null\n+++ b/test.txt\n@@ -0,0 +1 @@\n+Hello, World!\n"

	b.Run("DiskWrite", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mc := &types.MiniCommit{
				ID:        storage.GenerateID(patch, time.Now()),
				Message:   "Test commit",
				CreatedAt: time.Now(),
				Patch:     patch,
			}
			if err := storage.SaveMiniCommit(mc); err != nil {
				b.Fatalf("SaveMiniCommit() error = %v", err)
			}
		}
	})

	b.Run("DiskRead", func(b *testing.B) {
		// 事前にmini-commitを作成
		mc := &types.MiniCommit{
			ID:        storage.GenerateID(patch, time.Now()),
			Message:   "Test commit",
			CreatedAt: time.Now(),
			Patch:     patch,
		}
		if err := storage.SaveMiniCommit(mc); err != nil {
			b.Fatalf("SaveMiniCommit() error = %v", err)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := storage.GetMiniCommit(mc.ID)
			if err != nil {
				b.Fatalf("GetMiniCommit() error = %v", err)
			}
		}
	})
}

// BenchmarkStorageScalability ストレージのスケーラビリティベンチマーク
func BenchmarkStorageScalability(b *testing.B) {
	// テスト用Gitリポジトリを作成
	repo := NewTestGitRepo(&testing.T{})
	defer repo.Cleanup()

	// テスト用ストレージを作成
	storage, err := storage.NewStorage()
	if err != nil {
		b.Fatalf("Failed to create storage: %v", err)
	}

	// テスト用patchを作成
	patch := "diff --git a/test.txt b/test.txt\nnew file mode 100644\nindex 0000000..3b18e51\n--- /dev/null\n+++ b/test.txt\n@@ -0,0 +1 @@\n+Hello, World!\n"

	// 異なる数のmini-commitでテスト（CIでは重いテストをスキップ）
	sizes := []int{10, 100, 1000}
	if testing.Short() {
		sizes = []int{10, 100} // 短縮モードでは軽いテストのみ
	}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
			// 事前にmini-commitを作成
			for i := 0; i < size; i++ {
				mc := &types.MiniCommit{
					ID:        storage.GenerateID(patch, time.Now()),
					Message:   "Test commit",
					CreatedAt: time.Now(),
					Patch:     patch,
				}
				if err := storage.SaveMiniCommit(mc); err != nil {
					b.Fatalf("SaveMiniCommit() error = %v", err)
				}
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := storage.LoadMiniCommits()
				if err != nil {
					b.Fatalf("LoadMiniCommits() error = %v", err)
				}
			}
		})
	}
}