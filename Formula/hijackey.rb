class Hijackey < Formula
  desc "Hijack single-character keystrokes before they reach a command running under a pty"
  homepage "https://github.com/oooooooo/hijackey"
  url "https://github.com/oooooooo/hijackey.git",
      using:    :git,
      tag:      "v0.1.5",
      revision: "67e965c1d4d19bf70bc6ba73afcfbd1994a90a78"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", "-o", bin/"hijackey", "."
  end

  test do
    output = shell_output("#{bin}/hijackey 2>&1", 1)
    assert_match "usage:", output
  end
end
