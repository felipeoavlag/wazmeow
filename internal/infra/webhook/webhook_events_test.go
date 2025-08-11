package webhook

import (
	"testing"
	"time"

	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

func TestEventSerializer_SerializeContactEvent(t *testing.T) {
	serializer := NewEventSerializer()
	sessionID := "test-session"

	// Criar um evento Contact de teste
	contactEvent := &events.Contact{
		JID:          types.NewJID("5511999999999", types.DefaultUserServer),
		Timestamp:    time.Now(),
		Action:       nil, // Pode ser nil para teste básico
		FromFullSync: false,
	}

	// Serializar o evento
	payload, err := serializer.SerializeEvent(sessionID, contactEvent)
	if err != nil {
		t.Fatalf("Erro ao serializar evento Contact: %v", err)
	}

	// Verificar se o evento foi serializado corretamente
	if payload.Event != "contact" {
		t.Errorf("Esperado evento 'contact', obtido '%s'", payload.Event)
	}

	if payload.SessionID != sessionID {
		t.Errorf("Esperado sessionID '%s', obtido '%s'", sessionID, payload.SessionID)
	}

	// Verificar se os dados contêm o evento
	if payload.Data == nil {
		t.Fatal("Dados do payload não devem ser nil")
	}

	eventData, exists := payload.Data["event"]
	if !exists {
		t.Fatal("Campo 'event' não encontrado nos dados do payload")
	}

	// Verificar se o evento é do tipo correto
	if _, ok := eventData.(*events.Contact); !ok {
		t.Errorf("Esperado *events.Contact, obtido %T", eventData)
	}
}

func TestEventFilter_ShouldSendContactEvent(t *testing.T) {
	filter := NewEventFilter()

	// Verificar se "contact" está na lista de eventos suportados
	supportedEvents := filter.GetSupportedEvents()
	found := false
	for _, event := range supportedEvents {
		if event == "contact" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Evento 'contact' não encontrado na lista de eventos suportados")
	}
}

func TestEventFilter_ContactEventInGroup(t *testing.T) {
	filter := NewEventFilter()

	// Verificar se "contact" está no grupo "contacts"
	eventGroups := filter.GetEventGroups()
	contactsGroup, exists := eventGroups["contacts"]
	if !exists {
		t.Fatal("Grupo 'contacts' não encontrado")
	}

	found := false
	for _, event := range contactsGroup {
		if event == "contact" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Evento 'contact' não encontrado no grupo 'contacts'")
	}
}
