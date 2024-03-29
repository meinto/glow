# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Glow < Formula
  desc "A cli tool to adapt git-flow"
  homepage "https://github.com/meinto/glow"
  version "4.3.16"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/meinto/glow/releases/download/v4.3.16/glow_4.3.16_darwin_x86_64.tar.gz"
      sha256 "3d23b87890668a064206da86554ddf43353997ccab2c9a9b63ecb8f013b97ce8"

      def install
        bin.install "glow"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/meinto/glow/releases/download/v4.3.16/glow_4.3.16_darwin_arm64.tar.gz"
      sha256 "dc1990505058e7235bea28ed28d7aa4e78e452afec5a810a8a789dc4d83c1610"

      def install
        bin.install "glow"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/meinto/glow/releases/download/v4.3.16/glow_4.3.16_linux_arm64.tar.gz"
      sha256 "ca3b3e51e15f7f5d7322fc17ea26e033ced69202079c9002fb83a47833465048"

      def install
        bin.install "glow"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/meinto/glow/releases/download/v4.3.16/glow_4.3.16_linux_x86_64.tar.gz"
      sha256 "68ccebb7e334f6ca9e6c8634867eac521b093099481a422325e2af3876f906d4"

      def install
        bin.install "glow"
      end
    end
  end

  depends_on "git"
end
