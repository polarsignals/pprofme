# pprofme

Upload [pprof](https://github.com/google/pprof) formatted profiles with ease from your terminal (no sign up needed!).

## Install

```
curl -LO https://github.com/polarsignals/pprofme/releases/latest/download/pprofme-$(uname)-$(uname -m)
curl -LO https://github.com/polarsignals/pprofme/releases/latest/download/pprofme_checksums.txt
echo "$(cat pprofme_checksums.txt) pprofme" | shasum -a 256 --check
sudo mv 
```

## Usage

Run the `pprofme upload` with a path to a pprof profile, enter a description and the sharing link will be printed to your terminal and it will be opened in your default browser.

```
$ pprofme upload -d="Fibonacci in Go" ~/pprof/pprof.pprof-example-app-go.samples.cpu.001.pb.gz
https://pprof.me/779de8f
```

## License

Apache License 2.0, see [LICENSE](./LICENSE).
