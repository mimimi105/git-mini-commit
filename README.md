# git-mini-commit

Git CLI extension for "mini-commit" workflow: local, small commits between staging and regular commit. Written in Go.

## Description

git-mini-commit は、Git のステージングエリアと通常コミットの間に「mini-commit」という中間単位を導入し、大規模なリファクタリングや変更作業中の差分管理を容易にするツールです。

- mini-commit はローカル限定で管理され、Gitの push/fetch には影響を与えません
- ステージング中の差分を整理・確認することができます
- 複数の mini-commit をまとめて通常のコミットに統合可能です

---

## Features / 仕様

- ローカルの小さなコミット単位として変更を保存
- 保存した mini-commit を一覧表示
- 指定した mini-commit の差分を表示
- mini-commit をステージングに戻す（pop）
- mini-commit を削除（drop）
- 通常の git commit でまとめて反映

---

## Installation / バイナリ

### macOS (Homebrew)

```
brew install <your-tap>/git-mini-commit
```

### Linux / 手動ビルド

Go 1.20+ が必要です。

```
git clone https://github.com/<your-user>/git-mini-commit.git
cd git-mini-commit
go build -o git-mini-commit ./cmd/git-mini-commit
sudo mv git-mini-commit /usr/local/bin/
```

---

## Usage / コマンド一覧

- **Create mini-commit**

    ```
    git mini-commit -m "Refactor core module"
    ```

- **List mini-commits**

    ```
    git mini-commit list
    ```

- **Show diff of a mini-commit**

    ```
    git mini-commit show <hash>
    ```

- **Pop mini-commit back to staging**

    ```
    git mini-commit pop <hash>
    ```

- **Drop mini-commit**

    ```
    git mini-commit drop <hash>
    ```

- **Integrate mini-commits into a normal commit**

    ```
    git commit -m "まとめコミット"
    ```

---

## Directory Structure / 内部構造

- `.git/mini-commits/` に mini-commit パッチを保存
- 各 mini-commit は以下を保持:
    - ID（SHA1ハッシュ）
    - 作成日時
    - メッセージ
    - ステージング差分（patch形式）

---

## Contributing / 貢献方法

- Issue や Pull Request で提案・修正可能
- コードは Go のフォーマット `gofmt` に従う
- コミットメッセージは conventional commit 形式推奨

---

## Changelog / 変更履歴

- v0.1.0 初期リリース（mini-commit 作成、一覧、表示、pop、drop 機能）

---

## License / ライセンス

MIT License © 2025 Minoru Kinugasa

---

## Author / Contact
- GitHub: https://github.com/minoru-kinugasa-105
