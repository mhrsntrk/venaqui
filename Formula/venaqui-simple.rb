# Simple Homebrew formula for Venaqui
# This version builds from source
class Venaqui < Formula
  desc "CLI tool with TUI for downloading files via Real-Debrid and aria2"
  homepage "https://github.com/mhrsntrk/venaqui"
  url "https://github.com/mhrsntrk/venaqui.git",
      tag:      "v0.1.0",
      revision: "HEAD"
  license "MIT"
  head "https://github.com/mhrsntrk/venaqui.git", branch: "main"

  depends_on "go" => :build
  depends_on "aria2"

  def install
    system "go", "build", "-o", bin/"venaqui", "./cmd/venaqui"
  end

  test do
    system "#{bin}/venaqui", "version"
  end
end
