class GitMiniCommitBinary < Formula
  desc "Git CLI extension for mini-commit workflow: local, small commits between staging and regular commit"
  homepage "https://github.com/minoru-kinugasa-105/git-mini-commit"
  url "https://github.com/minoru-kinugasa-105/git-mini-commit/releases/download/v0.1.0/git-mini-commit-darwin-amd64"
  sha256 "PLACEHOLDER_SHA256_DARWIN_AMD64"
  license "MIT"

  if Hardware::CPU.arm?
    url "https://github.com/minoru-kinugasa-105/git-mini-commit/releases/download/v0.1.0/git-mini-commit-darwin-arm64"
    sha256 "PLACEHOLDER_SHA256_DARWIN_ARM64"
  end

  def install
    if Hardware::CPU.arm?
      bin.install "git-mini-commit-darwin-arm64" => "git-mini-commit"
    else
      bin.install "git-mini-commit-darwin-amd64" => "git-mini-commit"
    end
  end

  test do
    system "#{bin}/git-mini-commit", "--version"
  end
end
