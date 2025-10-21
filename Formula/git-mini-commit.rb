class GitMiniCommit < Formula
  desc "Git CLI extension for mini-commit workflow: local, small commits between staging and regular commit"
  homepage "https://github.com/mimimi105/git-mini-commit"
  url "https://github.com/mimimi105/git-mini-commit/archive/v0.1.0.tar.gz"
  sha256 "PLACEHOLDER_SHA256"
  license "MIT"
  head "https://github.com/mimimi105/git-mini-commit.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w")
  end

  test do
    system "#{bin}/git-mini-commit", "--version"
  end
end
