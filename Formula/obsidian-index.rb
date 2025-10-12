class ObsidianIndex < Formula
  desc "A CLI tool for indexing Obsidian vaults"
  homepage "https://github.com/nzb3/obsidian-index"
  url "https://github.com/nzb3/obsidian-index/archive/v1.0.0.tar.gz"
  sha256 "dce460e21304c3645424966cfac5dd2e58fe8e82286f92e0e5d447eb563f6b40"
  license "MIT"

  depends_on "go" => :build

  def install
    # Get git commit from the extracted source
    git_commit = `cd #{buildpath} && git rev-parse --short HEAD 2>/dev/null || echo "unknown"`.strip
    
    ldflags = %W[
      -s -w
      -X github.com/nzb3/obsidian-index/internal/version.Version=#{version}
      -X github.com/nzb3/obsidian-index/internal/version.GitCommit=#{git_commit}
      -X github.com/nzb3/obsidian-index/internal/version.BuildDate=#{Time.now.utc.iso8601}
    ].join(" ")

    system "go", "build", "-ldflags", ldflags, "-o", bin/"obsidian-index", "./cmd"
  end

  test do
    system "#{bin}/obsidian-index", "--help"
  end
end

