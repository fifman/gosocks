class Surlane < Formula
  desc "secure tunnel like shadowsocks, but lightweight"
  homepage ""
  url "https://github.com/fifman/surlane/releases/download/v0.1.0/surlane_0.1.0_Darwin_x86_64.tar.gz"
  version "0.1.0"
  sha256 "253ca20012b4a0f2dfa872d563bf02b09b1fa6719a7075715005fda22595ff48"

  def install
    bin.install "surlane"
  end

  def caveats; <<~EOS
    secure tunnel like shadowsocks, but lightweight
  EOS
  end

  plist_options :startup => false

  def plist; <<~EOS
    <?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>#{plist_name}</string>
    <key>ProgramArguments</key>
    <array>
        <string>#{opt_bin}/surlane</string>
        <string>-c /etc/surlane/config.toml</string>
    </array>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>

  EOS
  end

  test do
    system "#{bin}/surlane --version"
  end
end
