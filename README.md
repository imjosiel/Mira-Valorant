# Mira Valorant - Lupa Virtual Externa

Aplicativo de ampliação de tela desenvolvido em Go, focado em ser uma ferramenta externa e passiva para auxílio visual em jogos.

## Requisitos de Compilação
- **Go 1.20+**
- **Compilador C (GCC)**: Necessário para a biblioteca gráfica Fyne.
  - Windows: Instale [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) ou [MinGW-w64](https://www.mingw-w64.org/).

## Como Executar
1. Instale as dependências:
   ```bash
   go mod tidy
   ```
2. Execute o projeto:
   ```bash
   go run cmd/mira/main.go
   ```

## Estrutura do Projeto
- `cmd/mira`: Ponto de entrada (Main).
- `internal/config`: Estado compartilhado da aplicação.
- `internal/capture`: Lógica de captura de tela via API do Windows.
- `internal/ui`: Interface gráfica (Controle e Luneta).

## Anti-Cheat e Limitações (Vanguard / Outros)

### O que este app FAZ:
- Captura a imagem da tela usando APIs padrão do Windows (`BitBlt` / `GetCursorPos`).
- Renderiza uma janela "Always On Top" (Sempre no topo) sobre o jogo.
- Funciona de forma passiva, sem ler memória ou injetar código.

### O que este app NÃO FAZ:
- Não injeta DLLs.
- Não lê memória do processo do jogo.
- Não modifica os arquivos do jogo.

### Riscos e Observações:
1. **Captura de Tela**: Anti-cheats modernos como o Vanguard monitoram capturas de tela para evitar "Pixel Bots" (AimBots baseados em imagem). Embora este app seja apenas visual (Lupa), o uso contínuo de captura de tela *pode* ser flagrado heuristicamente se confundido com um bot.
2. **Janelas Sobrepostas (Overlays)**: Jogos em modo "Tela Cheia Exclusiva" podem impedir que janelas externas apareçam por cima. Recomenda-se usar o jogo em modo **"Janela em Tela Cheia" (Borderless Window)**.
3. **Bloqueio de Input**: A janela da luneta precisa ser "transparente" para o mouse (click-through) para não impedir que você atire/mire. O aplicativo tenta aplicar estilos do Windows (`WS_EX_TRANSPARENT`) para isso, mas pode haver conflitos dependendo da configuração do sistema.
4. **Desempenho**: Capturar e renderizar a tela em tempo real consome CPU/GPU. Se o FPS do jogo cair, reduza o tamanho da janela da luneta ou a taxa de atualização.

### Aviso Legal
Este software é uma ferramenta de acessibilidade visual. O uso em partidas competitivas deve ser feito com cautela e sob sua própria responsabilidade. Verifique os Termos de Serviço do jogo.
