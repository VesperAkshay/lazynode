class Lazynode < Formula
  desc "A powerful Terminal UI for managing Node.js projects"
  homepage "https://github.com/VesperAkshay/lazynode"
  version "VERSION_PLACEHOLDER"

  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/VesperAkshay/lazynode/releases/download/vVERSION_PLACEHOLDER/lazynode_VERSION_PLACEHOLDER_darwin_arm64.tar.gz"
      sha256 "SHA256_PLACEHOLDER"
    else
      url "https://github.com/VesperAkshay/lazynode/releases/download/vVERSION_PLACEHOLDER/lazynode_VERSION_PLACEHOLDER_darwin_amd64.tar.gz"
      sha256 "SHA256_PLACEHOLDER"
    end
  elsif OS.linux?
    if Hardware::CPU.arm?
      url "https://github.com/VesperAkshay/lazynode/releases/download/vVERSION_PLACEHOLDER/lazynode_VERSION_PLACEHOLDER_linux_arm64.tar.gz"
      sha256 "SHA256_PLACEHOLDER"
    else
      url "https://github.com/VesperAkshay/lazynode/releases/download/vVERSION_PLACEHOLDER/lazynode_VERSION_PLACEHOLDER_linux_amd64.tar.gz"
      sha256 "SHA256_PLACEHOLDER"
    end
  end

  def install
    bin.install "lazynode"
  end

  test do
    system "#{bin}/lazynode", "--version"
  end
end 