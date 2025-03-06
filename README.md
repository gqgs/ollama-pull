## ollama-pull


`ollama-pull` is a lightweight alternative to the official [ollama](https://github.com/ollama/ollama) pull command.

It aims to increase the speed of pulling from the [ollama](https://ollama.com/) registry by removing some overheads that might not be necessary for particular use cases. It supports both the HTTP standard library and `aria2c` download backends for enhanced performance and reliability.


### Usage
```bash
$ export OLLAMA_MODELS=/path/to/your/ollama/models
$ go run github.com/gqgs/ollama-pull/cmd/pull@latest deepseek-r1:1.5b
$ ollama serve
$ ollama run deepseek-r1:1.5b
>>> Hello World
<think>

</think>

Hello! How can I assist you today? 😊
```

### Comparison

##### ollama (7e402ebb8cc95e1fd2b59fe6d9ef9baf8972977e)

```bash
$ time ollama pull deepseek-r1:7b

pulling manifest 
pulling 96c415656d37... 100% ▕██████████████████████████████████████████████████▏ 4.7 GB                         
pulling 369ca498f347... 100% ▕██████████████████████████████████████████████████▏  387 B                         
pulling 6e4c38e1172f... 100% ▕██████████████████████████████████████████████████▏ 1.1 KB                         
pulling f4d24e9138dd... 100% ▕██████████████████████████████████████████████████▏  148 B                         
pulling 40fb844194b2... 100% ▕██████████████████████████████████████████████████▏  487 B                         
verifying sha256 digest 
writing manifest 
success 

real	2m40.397s
user	0m0.112s
sys	0m0.131s
```

##### ollama-pull (7b66627ee04a09deda83eef0683fbf1ca5b15775)

```bash
$ time go run github.com/gqgs/ollama-pull/cmd/pull@latest deepseek-r1:7b

2025/02/08 17:11:21 INFO downloading model manifest url=https://registry.ollama.ai/v2/library/deepseek-r1/manifests/7b
2025/02/08 17:11:21 INFO downloading blob blob=sha256:40fb844194b25e429204e5163fb379ab462978a262b86aadd73d8944445c09fd
2025/02/08 17:11:21 INFO downloading blob blob=sha256:369ca498f347f710d068cbb38bf0b8692dd3fa30f30ca2ff755e211c94768150
2025/02/08 17:11:21 INFO downloading blob blob=sha256:96c415656d377afbff962f6cdb2394ab092ccbcbaab4b82525bc4ca800fe8a49
2025/02/08 17:11:21 INFO downloading blob blob=sha256:f4d24e9138dd4603380add165d2b0d970bef471fac194b436ebd50e6147c6588
2025/02/08 17:11:21 INFO downloading blob blob=sha256:6e4c38e1172f42fdbff13edf9a7a017679fb82b0fde415a3e8b3c31c6ed4a4e4
2025/02/08 17:11:21 INFO writing blob to disk blob=sha256:6e4c38e1172f42fdbff13edf9a7a017679fb82b0fde415a3e8b3c31c6ed4a4e4
2025/02/08 17:11:21 INFO writing blob to disk blob=sha256:96c415656d377afbff962f6cdb2394ab092ccbcbaab4b82525bc4ca800fe8a49
2025/02/08 17:11:21 INFO writing blob to disk blob=sha256:369ca498f347f710d068cbb38bf0b8692dd3fa30f30ca2ff755e211c94768150
2025/02/08 17:11:21 INFO writing blob to disk blob=sha256:40fb844194b25e429204e5163fb379ab462978a262b86aadd73d8944445c09fd
2025/02/08 17:11:21 INFO writing blob to disk blob=sha256:f4d24e9138dd4603380add165d2b0d970bef471fac194b436ebd50e6147c6588

real	1m25.763s
user	0m2.437s
sys	0m8.887s
```

