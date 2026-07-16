class Hijackey < Formula
  desc "Hijack keystrokes for CLIs that won't let you rebind them"
  homepage "https://github.com/oooooooo/hijackey"
  url "https://github.com/oooooooo/hijackey.git",
      using:    :git,
      tag:      "v0.1.1",
      revision: "879166c0e2f668f69225e4f23d3f4332f78b773a"
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
