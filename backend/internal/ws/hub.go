// Package ws implementa a "sala ao vivo" da Etapa 5 do SISTEMA_MESTRE.md: um
// Hub em memória por campanha (não por conexão individual), sem Redis/pub-sub
// (uma instância só no Render, seria over-engineering agora — decisão já
// aprovada com o usuário antes de implementar).
//
// Desenho deliberadamente simples: o WebSocket aqui é só um canal de PUSH
// servidor→clientes conectados. Toda escrita (trocar cenário, mandar
// mensagem de chat) continua passando pelas rotas REST normais, autenticadas
// e validadas do jeito de sempre — o handler REST correspondente, depois de
// persistir a mudança, chama Manager.Broadcast pra avisar quem está
// conectado. Isso evita ter que desenhar um protocolo bidirecional por cima
// do WebSocket (autenticação, validação e idempotência de cada tipo de
// evento já existem nas rotas REST, não precisam ser duplicadas aqui) — a
// única coisa que chega pela conexão WS de um cliente é o ping/pong padrão
// e a detecção de desconexão.
package ws

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/coder/websocket"
)

func writeJSON(ctx context.Context, conn *websocket.Conn, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return conn.Write(ctx, websocket.MessageText, data)
}

// Event é o envelope de todo evento — campo Type extensível sem redesenhar
// nada (scene_changed, chat_message, e mais no futuro).
type Event struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type client struct {
	conn *websocket.Conn
	send chan Event
}

type hub struct {
	mu      sync.Mutex
	clients map[*client]bool
}

// Manager mantém um hub por campanha, criado sob demanda na primeira conexão
// e nunca removido explicitamente (o mapa de clients de um hub vazio não
// custa nada relevante de memória — não vale a complexidade de expirar hubs
// vazios nesta escala).
type Manager struct {
	mu   sync.Mutex
	hubs map[uint]*hub
}

func NewManager() *Manager { return &Manager{hubs: make(map[uint]*hub)} }

func (m *Manager) getOrCreateHub(campaignID uint) *hub {
	m.mu.Lock()
	defer m.mu.Unlock()
	h, ok := m.hubs[campaignID]
	if !ok {
		h = &hub{clients: make(map[*client]bool)}
		m.hubs[campaignID] = h
	}
	return h
}

// Broadcast envia um evento pra todos os clientes conectados no hub dessa
// campanha agora. Não-bloqueante: um cliente lento (buffer cheio) perde o
// evento em vez de travar o broadcast pros outros — aceitável pra esse tipo
// de evento (o cliente pega o estado atual de novo via REST se precisar,
// ex: reabrir a página busca o cenário ativo da sessão via GET normal).
func (m *Manager) Broadcast(campaignID uint, event Event) {
	m.mu.Lock()
	h, ok := m.hubs[campaignID]
	m.mu.Unlock()
	if !ok {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	for c := range h.clients {
		select {
		case c.send <- event:
		default:
		}
	}
}

// Join bloqueia até a conexão cair — chamado a partir da goroutine do
// handler HTTP que fez o upgrade.
func (m *Manager) Join(ctx context.Context, campaignID uint, conn *websocket.Conn) {
	h := m.getOrCreateHub(campaignID)
	cl := &client{conn: conn, send: make(chan Event, 16)}

	h.mu.Lock()
	h.clients[cl] = true
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients, cl)
		h.mu.Unlock()
	}()

	writerDone := make(chan struct{})
	go func() {
		defer close(writerDone)
		for {
			select {
			case ev, ok := <-cl.send:
				if !ok {
					return
				}
				if err := writeJSON(ctx, conn, ev); err != nil {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Loop de leitura só serve pra detectar desconexão (o cliente não manda
	// nada de significativo por aqui, ver comentário do pacote) — qualquer
	// erro de leitura (conexão fechada, ping perdido) encerra o Join.
	for {
		_, _, err := conn.Read(ctx)
		if err != nil {
			break
		}
	}
	<-writerDone
}
