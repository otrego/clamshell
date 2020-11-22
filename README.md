# Clamshell

Clamshell is a library of Go(lang) functions for the game of go (Aka baduk,
weiqi, igo).

For getting started with Clamshell and the technologies we use (Git, Go), check
out our [Contributing guide](CONTRIBUTING.md).

## Testing

To run the all the unit tests, run:

```shell
go test ./...
```

## Setting up a Dev Environment

To start doing basic Pull Requests (PRs), you'll need the following configured:

1.   [Go(lang)](https://golang.org/doc/install)
2.   [Git](https://git-scm.com/book/en/v2)
3.   Some IDE to edit the code.

To work on katalyze / katago-related KRs, you'll need katago installed:

1.   [Katago](https://github.com/lightvector/KataGo).

To work on the Server codebase, you'll need a couple more tools:

1.   [Docker](https://www.docker.com/get-started). In particular, Docker Desktop.

### Inital OS Setup

*   **OSX**
    *   We recommend using [Homebrew](https://brew.sh/) to install any packages you may need.

*   **Linux**: TODO. Send us a PR with instructions!

*   **Windows**: TODO. Send us a PR with instructions!

### Installing KataGo

[KataGo](https://github.com/lightvector/KataGo) works with either
[CUDA](https://developer.nvidia.com/cuda-zone) or
[OpenCL](https://www.khronos.org/opencl/). With OpenCL, you can use it without
needing a GPU machine.

In addition to the built-binaries, KataGo needs 3 configuration files to run:

1.   Model: You can get that from the [KataGo releases](https://github.com/lightvector/KataGo/releases)
  * *Note: There are a number of different models. If you want some nets that are much faster to run, try any of the "b10c128" or "b15c192" Extended Training Nets here which are 10 block and 15 block networks from earlier in the run that are much weaker but still pro-level-and-beyond. [https://d3dndmfyhecmj0.cloudfront.net/g170/neuralnets/g170-b10c128-s197428736-d67404019.bin.gz](https://d3dndmfyhecmj0.cloudfront.net/g170/neuralnets/g170-b10c128-s197428736-d67404019.bin.gz)*
2.   GTP Config: You can get that from the [KataGo releases](https://github.com/lightvector/KataGo/releases)
3.   Tuning Parameters. This is set via running `katago benchmark -tune`.`

*   **OSX**
    *   OSX generally comes pre-installed with [OpenCL](https://www.khronos.org/opencl/)
    *   Run `brew install katago`
    *   Download the models. Brew comes with a models, which you can get with
        *   `KATAGO_GTP_CONFIG=$(brew list --verbose katago | grep gtp)`
        *   `KATAGO_MODEL_PATH= $(brew list --verbose katago | grep .gz | head -1)`
    *   Tune Katago: `katago benchmark -tune -config $(KATAGO_GTP_CONFIG) -model $(KATAGO_MODEL_PATH)`
        *   This will output configuration to `$HOME/.katago`
    *   KataGo should now be operational!

### Using KataGo via UI

*   **OSX/Linux**
    *   If you want to try it out with a UI, download
        [Lizzie](https://github.com/featurecat/lizzie) and change the engine
        command to: `/path/to/katago gtp -config /path/to/config.cfg -model /path/to/model.gz`
        *   I had better luck editing the `config.txt` file that
            comes with Lizzie directly, rather than trying to set the engine in the Java UI.
