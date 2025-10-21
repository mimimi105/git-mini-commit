#!/usr/bin/env ruby

require 'net/http'
require 'json'
require 'digest'

# GitHub API から最新リリース情報を取得
def get_latest_release
  uri = URI('https://api.github.com/repos/mimimi105/git-mini-commit/releases/latest')
  response = Net::HTTP.get_response(uri)
  JSON.parse(response.body)
end

# ファイルのSHA256を計算
def calculate_sha256(url)
  uri = URI(url)
  response = Net::HTTP.get_response(uri)
  Digest::SHA256.hexdigest(response.body)
end

# Formula を更新
def update_formula(release)
  formula_path = 'Formula/git-mini-commit-binary.rb'
  version = release['tag_name'].gsub('v', '')
  
  # 各プラットフォームのバイナリURLとSHA256を取得
  assets = release['assets']
  darwin_amd64_url = assets.find { |a| a['name'].include?('darwin-amd64') }&.dig('browser_download_url')
  darwin_arm64_url = assets.find { |a| a['name'].include?('darwin-arm64') }&.dig('browser_download_url')
  
  darwin_amd64_sha256 = calculate_sha256(darwin_amd64_url) if darwin_amd64_url
  darwin_arm64_sha256 = calculate_sha256(darwin_arm64_url) if darwin_arm64_url
  
  # Formula テンプレート
  formula_content = <<~RUBY
    class GitMiniCommitBinary < Formula
      desc "Git CLI extension for mini-commit workflow: local, small commits between staging and regular commit"
      homepage "https://github.com/mimimi105/git-mini-commit"
      url "#{darwin_amd64_url}"
      sha256 "#{darwin_amd64_sha256}"
      license "MIT"

      if Hardware::CPU.arm?
        url "#{darwin_arm64_url}"
        sha256 "#{darwin_arm64_sha256}"
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
  RUBY
  
  File.write(formula_path, formula_content)
  puts "Updated #{formula_path} with version #{version}"
end

# メイン実行
begin
  release = get_latest_release
  update_formula(release)
rescue => e
  puts "Error: #{e.message}"
  exit 1
end
