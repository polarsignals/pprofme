# pprofme

Upload [pprof](https://github.com/google/pprof) formatted profiles with ease from your terminal (no sign up needed!).

## Install

```
# Download binary for your OS and architecture
curl -LO https://github.com/polarsignals/pprofme/releases/latest/download/pprofme_$(uname)_$(uname -m)
# Verify the checksum
curl -sL https://github.com/polarsignals/pprofme/releases/latest/download/pprofme_checksums.txt | shasum --ignore-missing -a 256 --check
# Make the binary executable
chmod a+x pprofme_$(uname)_$(uname -m)
# Move to path
sudo mv pprofme_$(uname)_$(uname -m) /usr/local/bin/pprofme
```

## Usage

Run the `pprofme upload` with a path to a pprof profile, enter a description and the sharing link will be printed to your terminal and it will be opened in your default browser.

```
$ pprofme upload -d "Fibonacci in Go" ./fibonacci.pb.gz
https://pprof.me/779de8f
```

## License

Apache License 2.0, see [LICENSE](./LICENSE).
