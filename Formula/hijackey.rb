class Hijackey < Formula
  desc "Hijack keystrokes for CLIs that won't let you rebind them"
  homepage "https://github.com/oooooooo/hijackey"
  url "https://github.com/oooooooo/hijackey.git",
      using:    :git,
      tag:      "v0.1.2",
      revision: "a0d4634068f02dc0fa6b6fc0952f04a1cd7790d5"
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
