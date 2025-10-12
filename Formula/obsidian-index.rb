class ObsidianIndex < Formula
  desc "CLI tool for indexing Obsidian vaults"
  homepage "https://github.com/nzb3/obsidian-index"
  url "https://github.com/nzb3/obsidian-index/archive/v1.0.0.tar.gz"
  sha256 "0d019aa76d7de087538fbe68efc4e3ca3fe623cdb51d2378a8f81872b9456a74"
  license "MIT"
  head "https://github.com/nzb3/obsidian-index.git", branch: "main"

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

