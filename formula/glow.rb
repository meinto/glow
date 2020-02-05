# This file was generated by GoReleaser. DO NOT EDIT.
class Glow < Formula
  desc ""
  homepage ""
  version "4.2.0"
  bottle :unneeded

  if OS.mac?
    url "https://github.com/meinto/glow/releases/download/v4.2.0/glow_4.2.0_darwin_x86_64.tar.gz"
    sha256 "f70bc8cfd04cf136fc7a34c7d187ac1ae42a7ad445cb49b8634b50e078665417"
  elsif OS.linux?
    if Hardware::CPU.intel?
      url "https://github.com/meinto/glow/releases/download/v4.2.0/glow_4.2.0_linux_x86_64.tar.gz"
      sha256 "4579e56ffcdea64993b965aa390b7bc46cb407c30a038b773914e229da0fdf39"
    end
  end
  
  depends_on "git"

  def install
    bin.install "glow"
  end
end
