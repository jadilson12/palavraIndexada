# Palavra Indexada

[![CI](https://github.com/jadilson12/palavraindexada/actions/workflows/ci.yaml/badge.svg)](https://github.com/jadilson12/palavraindexada/actions/workflows/ci.yaml)
[![Netlify Status](https://api.netlify.com/api/v1/badges/6e83fd88-5ffe-4808-9689-c0f3b100bfe3/deploy-status)](https://app.netlify.com/sites/palavraindexada/deploys)
[![License: CC BY-NC-SA 4.0](https://img.shields.io/badge/License-CC%20BY--NC--SA%204.0-lightgrey.svg)](https://creativecommons.org/licenses/by-nc-sa/4.0/)

Projeto com artigos, ministrações, seminários e conhecimento da Palavra de Deus, otimizado para indexação, uso de IA, AIO e SEO. Construído com [Hextra](https://github.com/imfing/hextra).

## Desenvolvimento Local

Pré-requisitos: [Hugo](https://gohugo.io/getting-started/installing/), [Go](https://golang.org/doc/install) e [Git](https://git-scm.com)

```shell
# Clonar o repositório
git clone https://github.com/jadilson12/palavraindexada.git
cd palavraindexada

# Iniciar o servidor
hugo mod tidy
hugo server --logLevel debug --disableFastRender -p 1313
```

### Atualizar o tema

```shell
hugo mod get -u
hugo mod tidy
```

## Deploy

O deploy é feito automaticamente pelo [Netlify](https://palavraindexada.netlify.app/) a cada push na branch `main`.

## Licença

Este trabalho está licenciado sob [Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International](https://creativecommons.org/licenses/by-nc-sa/4.0/).

[![CC BY-NC-SA 4.0](https://licensebuttons.net/l/by-nc-sa/4.0/88x31.png)](https://creativecommons.org/licenses/by-nc-sa/4.0/)
