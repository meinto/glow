# This file was generated by GoReleaser. DO NOT EDIT.
class Glow < Formula
  desc ""
  homepage ""
  url "https://github.com/meinto/glow/releases/download/v1.10.5/glow_1.10.5_darwin_x86_64.tar.gz"
  version "1.10.5"
  sha256 "53654ad8333c247697d8809221780582d9e6e0aa60ea143aa2db5ae6e088d934"
  
  depends_on "git"

  def install
    bin.install "glow"
  end
end
