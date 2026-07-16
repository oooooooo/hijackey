class Hijackey < Formula
  desc "Hijack single-character keystrokes before they reach a command running under a pty"
  homepage "https://github.com/oooooooo/hijackey"
  url "https://github.com/oooooooo/hijackey.git",
      using:    :git,
      tag:      "v0.1.0",
      revision: "e2f8fd173f77d0f49e196fc0a81c9192a7bd75b8"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", "-ldflags", "-X main.version=v#{version}", "-o", bin/"hijackey", "."
  end

  test do
    output = shell_output("#{bin}/hijackey 2>&1", 1)
    assert_match "usage:", output
  end
end
