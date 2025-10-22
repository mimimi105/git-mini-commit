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

## Usage / コマンド一覧

### 基本的な使用方法

- **Create mini-commit（mini-commitを作成）**

    ```bash
    git mini-commit -m "Refactor core module"
    ```

- **List mini-commits（mini-commit一覧表示）**

    ```bash
    git mini-commit list
    ```

- **Show diff of a mini-commit（mini-commitの差分表示）**

    ```bash
    git mini-commit show <hash>
    ```

- **Pop mini-commit back to staging（mini-commitをステージングに戻す）**

    ```bash
    git mini-commit pop <hash>
    ```

- **Drop mini-commit（mini-commitを削除）**

    ```bash
    git mini-commit drop <hash>
    ```

- **Integrate mini-commits into a normal commit（mini-commitを統合してコミット）**

    ```bash
    git commit -m "まとめコミット"
    ```

### 使用例

```bash
# 1. ファイルをステージング
git add src/main.go

# 2. mini-commitとして保存
git mini-commit -m "メイン関数のリファクタリング"

# 3. 別のファイルをステージング
git add src/utils.go

# 4. 別のmini-commitとして保存
git mini-commit -m "ユーティリティ関数の追加"

# 5. mini-commit一覧を確認
git mini-commit list

# 6. すべてを統合してコミット（標準Gitコマンド）
git commit -m "機能追加とリファクタリング"
```

---

## Directory Structure / 内部構造

- `.git/mini-commits/` に mini-commit パッチを保存
- 各 mini-commit は以下を保持:
    - ID（SHA1ハッシュ）
    - 作成日時
    - メッセージ
    - ステージング差分（patch形式）

### ファイル命名規則

```
.git/mini-commits/
├── index.json           # mini-commit一覧のインデックス
├── <hash>.patch         # 各mini-commitのpatchファイル
└── <hash>.patch         # (例: a1b2c3d4.patch)
```

- **ID生成**: `SHA1(patch内容 + タイムスタンプ)` で生成
- **patchファイル**: `<hash>.patch` 形式で保存
- **インデックス**: `index.json` で一覧管理

## 制約事項 / Limitations

- **GUI表示不可**: VSCode Gitタブ、GitHub Desktop、SourceTreeなどのGUIツールには表示されません
- **ローカル限定**: `git push`や`git fetch`には影響しません
- **標準Gitコマンドとの分離**: `git log`、`git status`などには表示されません
- **統合は標準Gitコマンド**: `git commit`でmini-commitが統合されます
- **統合順序**: 作成順（古いものから新しいものへ）で統合されます

## 差分確認方法 / Diff Inspection

mini-commitの差分をより詳細に確認する方法:

```bash
# 1. mini-commitの一覧を表示
git mini-commit list

# 2. 特定のmini-commitの差分を表示
git mini-commit show <hash>

# 3. patchファイルを直接確認（上級者向け）
cat .git/mini-commits/<hash>.patch

# 4. patchの統計情報を確認
git mini-commit show <hash> | git apply --stat

# 5. patchを一時的に適用して確認
git mini-commit show <hash> | git apply --check
```

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

## Installation / インストール方法

### 1. npm (推奨 - クロスプラットフォーム)

```bash
npm install -g git-mini-commit
```

### 2. Homebrew (macOS/Linux)

```bash
# ワンライナー（推奨）
brew install mimimi105/git-mini-commit/git-mini-commit

# または手動で
brew tap mimimi105/git-mini-commit
brew install git-mini-commit
```

### 3. Scoop (Windows)

```bash
scoop bucket add git-mini-commit https://github.com/mimimi105/git-mini-commit
scoop install git-mini-commit
```

### 4. 直接ダウンロード (GitHub Release)

```bash
# Linux (AMD64)
curl -L https://github.com/mimimi105/git-mini-commit/releases/latest/download/git-mini-commit-linux-amd64 -o git-mini-commit
chmod +x git-mini-commit
sudo mv git-mini-commit /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/mimimi105/git-mini-commit/releases/latest/download/git-mini-commit-darwin-amd64 -o git-mini-commit
chmod +x git-mini-commit
sudo mv git-mini-commit /usr/local/bin/

# macOS (Apple Silicon)
curl -L https://github.com/mimimi105/git-mini-commit/releases/latest/download/git-mini-commit-darwin-arm64 -o git-mini-commit
chmod +x git-mini-commit
sudo mv git-mini-commit /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/mimimi105/git-mini-commit/releases/latest/download/git-mini-commit-windows-amd64.exe" -OutFile "git-mini-commit.exe"
```

### 5. Debian/Ubuntu パッケージ

```bash
# 自動インストールスクリプト
curl -fsSL https://raw.githubusercontent.com/mimimi105/git-mini-commit/main/scripts/install-debian.sh | bash

# または手動で .deb パッケージをインストール
wget https://github.com/mimimi105/git-mini-commit/releases/latest/download/git-mini-commit_0.1.0-1_amd64.deb
sudo dpkg -i git-mini-commit_0.1.0-1_amd64.deb
```

### 6. ソースからビルド

```bash
git clone https://github.com/mimimi105/git-mini-commit.git
cd git-mini-commit
go build -o git-mini-commit .
sudo mv git-mini-commit /usr/local/bin/
```

## 対応プラットフォーム / Supported Platforms

| OS      | Architecture  | npm | Homebrew | Scoop | Direct Download | Debian |
| ------- | ------------- | --- | -------- | ----- | --------------- | ------ |
| Linux   | AMD64         | ✅  | ✅       | ❌    | ✅              | ✅     |
| Linux   | ARM64         | ✅  | ✅       | ❌    | ✅              | ✅     |
| macOS   | Intel         | ✅  | ✅       | ❌    | ✅              | ❌     |
| macOS   | Apple Silicon | ✅  | ✅       | ❌    | ✅              | ❌     |
| Windows | AMD64         | ✅  | ❌       | ✅    | ✅              | ❌     |
| Windows | ARM64         | ✅  | ❌       | ✅    | ✅              | ❌     |

## インストール後の確認

```bash
git mini-commit --version
```

## Author / Contact

- GitHub: https://github.com/mimimi105
