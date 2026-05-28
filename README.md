# Web Proxy com Controle de Conteúdo

Miguel Casarim da Silva  
Guilherme Estrella

Este projeto consiste na implementação de um web proxy didático com controle de conteúdo, desenvolvido para a disciplina de Sistemas para Internet 2. O proxy atua como um intermediário entre o cliente (navegador) e a Internet, interceptando requisições HTTP e HTTPS para repassar, bloquear ou filtrar o conteúdo baseado em regras pré-configuradas.

## Por que escolhemos Go (Golang)?

A escolha da tecnologia foi fundamentada nos seguintes pilares técnicos da linguagem Go:

* **Ausência de Frameworks Pesados:** A biblioteca padrão do Go (`net/http`) é extremamente robusta. Ela permitiu construir o servidor HTTP, manipular requisições e gerenciar sockets TCP de forma nativa, sem a necessidade de instalar pacotes ou frameworks de terceiros.
* **Paralelismo e Concorrência Nativa:** Um servidor proxy precisa lidar com múltiplas requisições simultâneas. Go soluciona isso nativamente através das **Goroutines**, que são funções executadas de forma concorrente. Elas consomem uma fração mínima da memória se comparadas às threads tradicionais de sistemas operacionais.
* **Eficiência no Túnel HTTPS:** Para a implementação do método `CONNECT`, utilizamos Goroutines para ler e copiar o tráfego de dados bidirecionalmente (cliente ↔ servidor de destino) em tempo real, garantindo alta performance e evitando o bloqueio do fluxo principal do servidor.

## Estrutura do Projeto e Lógica dos Arquivos

O projeto está modularizado para separar claramente as responsabilidades de configuração, registro e manipulação de tráfego.

### Diretório Principal e Configurações
* **`config.go`**: Responsável por ler e carregar na memória os arquivos JSON de regras (`blocked.json` e `words.json`). Contém a lógica `IsBlocked` para verificar a presença de um domínio na lista de bloqueio de forma *case-insensitive*.
* **`logger.go`**: Registra todas as requisições em um histórico estruturado, armazenando o *Timestamp*, a *URL* e a *Ação* executada. Como o ambiente é concorrente (múltiplas Goroutines acessando o mesmo recurso), foi implementado um `sync.Mutex` para garantir exclusão mútua nas operações de escrita, evitando condições de corrida (*race conditions*) e a corrupção do arquivo de log.

### Handlers (Manipuladores de Requisição)
* **`proxy_handler.go`**: O roteador principal do proxy. Ele intercepta as requisições, faz o parsing da URL de destino e avalia qual ação tomar com base nas configurações carregadas, direcionando o fluxo para o Handler específico.
* **`pass_handler.go` (Repasse)**: Utilizado quando o site solicitado é livre. Ele clona a requisição original do cliente, limpa cabeçalhos de codificação incompatíveis (como `Accept-Encoding`) para permitir análise subsequente caso necessário, efetua o disparo ao servidor de destino e devolve a resposta integral ao cliente.
* **`block_handler.go` (Bloqueio de Sites)**: Se o domínio estiver listado no arquivo de bloqueios, a requisição externa é abortada. O proxy responde diretamente ao cliente com o status `403 Forbidden` e renderiza um template HTML customizado a partir de `templates/blocked.html`.
* **`filter_handler.go` (Filtro de Conteúdo)**: Realiza a requisição ao site de origem e inspeciona se o `Content-Type` é do tipo `text/html`. Sendo positivo, o corpo da página é interceptado e processado por expressões regulares, que localizam as palavras proibidas de forma *case-insensitive* e as substituem pelos termos equivalentes configurados.
* **`connect_handler.go` (Tunelamento HTTPS)**: Manipula requisições que utilizam o método HTTP `CONNECT`. Ele estabelece uma conexão TCP pura (`net.Dial`) com o servidor de destino, utiliza o recurso `http.Hijacker` para assumir o controle total do socket TCP do cliente e inicia duas Goroutines paralelas usando `io.Copy` para trafegar os dados criptografados de forma transparente entre as pontas.

## Pré-requisitos e Instalação

1. **Instalar o Go:** Certifique-se de ter o Go instalado em sua máquina (versão 1.25 ou mais recentes). Download disponível em [go.dev](https://go.dev/dl/).
2. **Clonar o Repositório:**
   ```bash
   git clone <https://github.com/Miguel-casarin/Web-Proxy->
   cd <Web-Proxy->
3. Dependências Externas: Nenhuma. O projeto utiliza os pacotes nativos da biblioteca padrão do Go.

### Configuração das Listas de Controle
O comportamento do proxy é ditado por dois arquivos JSON localizados na raiz do projeto:

1. `blocked.json` (Lista de Sites Bloqueados):
2. `words.json` (Filtro de Termos):

### Como Executar o Proxy
Abra o terminal na pasta raiz do projeto e execute o comando:

    go run .

O proxy será iniciado localmente e ficará ouvindo novas requisições na porta 5000.

### Configuração do Navegador (Firefox)
Para que o tráfego do navegador seja interceptado e processado pelas regras do proxy, é obrigatório realizar a configuração no browser. Nossos testes e validações foram executados no Mozilla Firefox seguindo as etapas abaixo:

1. Acesse as Configurações.

2. Localize Configuracoes de proxy e acesse a opção "Configurar proxy"

3. Marque a opção "Configuração manual de proxy"

4. Na caixa "Proxy HTTP", defina o endereço como localhost e defina a Porta como 5000.

5. Marque a opção "Usar este proxy também em HTTPS" (esta etapa é essencial para que o método CONNECT no connect_handler.go funcione corretamente em sites seguros).

Clique em OK para aplicar as alterações.

**Nota:** Para desativar o proxy após os testes, basta voltar a essa mesma tela e reverter para a opção "Usar as configurações de proxy do sistema".

### Transparência no Uso de IA
Nota: Declaramos de forma transparente como ferramentas de inteligência artificial apoiaram o desenvolvimento do projeto.

Ferramentas utilizadas: Claude

Como foram usadas: Atuaram de forma consultiva para sanar dúvidas sobre convenções de nomenclatura da linguagem (como o uso de Pascal Case para exportação de funções), estruturação de expressões regulares eficientes no pacote regexp e para compreender o funcionamento prático do http.Hijacker na transição do fluxo HTTP para sockets TCP brutos.