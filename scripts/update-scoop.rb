#!/usr/bin/env ruby

require 'net/http'
require 'json'
require 'digest'

# GitHub API から最新リリース情報を取得
def get_latest_release
  uri = URI('https://api.github.com/repos/minoru-kinugasa-105/git-mini-commit/releases/latest')
  response = Net::HTTP.get_response(uri)
  JSON.parse(response.body)
end

# ファイルのSHA256を計算
def calculate_sha256(url)
  uri = URI(url)
  response = Net::HTTP.get_response(uri)
  Digest::SHA256.hexdigest(response.body)
end

# Scoop レシピを更新
def update_scoop_recipe(release)
  recipe_path = 'scoop/git-mini-commit.json'
  version = release['tag_name'].gsub('v', '')
  
  # 各プラットフォームのバイナリURLとSHA256を取得
  assets = release['assets']
  windows_amd64_url = assets.find { |a| a['name'].include?('windows-amd64') }&.dig('browser_download_url')
  windows_arm64_url = assets.find { |a| a['name'].include?('windows-arm64') }&.dig('browser_download_url')
  
  windows_amd64_sha256 = calculate_sha256(windows_amd64_url) if windows_amd64_url
  windows_arm64_sha256 = calculate_sha256(windows_arm64_url) if windows_arm64_url
  
  # Scoop レシピテンプレート
  recipe_content = {
    "version" => version,
    "description" => "Git CLI extension for mini-commit workflow: local, small commits between staging and regular commit",
    "homepage" => "https://github.com/minoru-kinugasa-105/git-mini-commit",
    "license" => "MIT",
    "architecture" => {
      "64bit" => {
        "url" => windows_amd64_url,
        "hash" => windows_amd64_sha256,
        "bin" => "git-mini-commit.exe"
      },
      "arm64" => {
        "url" => windows_arm64_url,
        "hash" => windows_arm64_sha256,
        "bin" => "git-mini-commit.exe"
      }
    },
    "checkver" => "github",
    "autoupdate" => {
      "architecture" => {
        "64bit" => {
          "url" => "https://github.com/minoru-kinugasa-105/git-mini-commit/releases/download/v$version/git-mini-commit-windows-amd64.exe"
        },
        "arm64" => {
          "url" => "https://github.com/minoru-kinugasa-105/git-mini-commit/releases/download/v$version/git-mini-commit-windows-arm64.exe"
        }
      }
    }
  }
  
  File.write(recipe_path, JSON.pretty_generate(recipe_content))
  puts "Updated #{recipe_path} with version #{version}"
end

# メイン実行
begin
  release = get_latest_release
  update_scoop_recipe(release)
rescue => e
  puts "Error: #{e.message}"
  exit 1
end
