# LocalSend

Transferencia de arquivos entre dispositivos na rede local, construida em Go.

## O que e este projeto

LocalSend le um arquivo, divide-o em chunks de 4 MB e processa cada chunk de forma concorrente usando um pool de workers. O objetivo final e enviar esses chunks para outro dispositivo na rede local, onde serao remontados no arquivo original.

## Etapa atual

A transferencia pela rede ainda nao foi implementada. Hoje o programa apenas:

1. Abre o arquivo informado na linha de comando.
2. Divide o arquivo em chunks de 4 MB (`ChunkSize` em `BuildChunks.go`).
3. Distribui os chunks para 8 workers concorrentes.
4. Valida o processamento: nenhum chunk pode ficar perdido, duplicado ou com erro.

A funcao `UploadChunk` em `Upload.go` e um stub que retorna `nil` e e o ponto de entrada para a proxima etapa: o envio dos chunks pela rede local.

## Como foi construido

- **Linguagem:** Go 1.26
- **Divisao em chunks:** `BuildChunks.go` calcula os offsets e tamanhos de cada chunk a partir do tamanho total do arquivo.
- **Leitura dos chunks:** `ReadChunk.go` usa `ReadAt` para ler cada chunk em uma posicao especifica do arquivo, permitindo leituras paralelas seguras.
- **Concorrencia:** `Workers.go` implementa o pool de workers; `main.go` cria 8 goroutines que consomem chunks do canal `jobs` e publicam o resultado no canal `results`.
- **Validacao:** `Validate.go` confere se todos os chunks foram processados exatamente uma vez, sem erros, e exibe um relatorio final.

## Como executar

```bash
go run ./cmd <caminho-do-arquivo>
```

Exemplo:

```bash
go run ./cmd /home/ivan/Videos/video.mp4
```

## Estrutura do projeto

```
cmd/
  main.go                 Ponto de entrada: orquestra workers e canais
  internal/
    BuildChunks.go        Divide o arquivo em chunks de 4 MB
    ReadChunk.go          Define a struct Transfer e a leitura de chunks
    Upload.go             Stub de upload (proxima etapa)
    Workers.go            Pool de workers e struct Result
    Validate.go           Valida processamento de todos os chunks
    workers_test.go       Testes dos workers
```

## Proximos passos

- Implementar o envio dos chunks pela rede local em `UploadChunk`.
- Definir o protocolo de transferencia (servidor receptor que remonta o arquivo).
- Exibir progresso do envio e tratar reenvio de chunks com falha.
