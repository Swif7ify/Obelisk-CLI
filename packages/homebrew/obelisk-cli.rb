class ObeliskCli < Formula
  desc "AI-Powered Automated Tech Lead for Modern Codebases"
  homepage "https://github.com/Swif7ify/Obelisk-CLI"
  url "https://github.com/Swif7ify/Obelisk-CLI/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "INSERT_SHA256_HERE"
  license "MIT"
  head "https://github.com/Swif7ify/Obelisk-CLI.git", branch: "main"

  depends_on "go" => :build

  def install
    # Build the binary
    system "go", "build", *std_go_args(ldflags: "-s -w -X main.Version=#{version}"), "-o", bin/"obelisk"
  end

  test do
    # Test that the binary runs
    assert_match version.to_s, shell_output("#{bin}/obelisk version")
    
    # Test help command
    assert_match "AI-Powered Automated Tech Lead", shell_output("#{bin}/obelisk --help")
  end

  def caveats
    <<~EOS
      Obelisk CLI has been installed!

      To get started:
        1. Set your Gemini API key:
           obelisk config set api-key YOUR_API_KEY

        2. Run the interactive TUI:
           obelisk

        3. Or run a quick check:
           obelisk check --path /path/to/project

      Documentation: https://github.com/Swif7ify/Obelisk-CLI
    EOS
  end
end

# Made with Bob
