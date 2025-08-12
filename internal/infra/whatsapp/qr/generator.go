package qr

import (
	"encoding/base64"
	"os"

	"github.com/mdp/qrterminal/v3"
	"github.com/skip2/go-qrcode"
)

// Generator gera e exibe QR codes de forma otimizada
type Generator struct{}

// NewGenerator cria um novo gerador de QR codes
func NewGenerator() *Generator {
	return &Generator{}
}

// GenerateBase64PNG gera QR code como PNG em base64
func (g *Generator) GenerateBase64PNG(code string) (string, error) {
	// Gerar QR code como PNG
	pngBytes, err := qrcode.Encode(code, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}

	// Converter para base64
	base64String := base64.StdEncoding.EncodeToString(pngBytes)
	return base64String, nil
}

// DisplayTerminal exibe QR code no terminal
func (g *Generator) DisplayTerminal(code string) {
	// Configurar para exibir no terminal
	config := qrterminal.Config{
		Level:     qrterminal.M,
		Writer:    os.Stdout,
		BlackChar: qrterminal.BLACK,
		WhiteChar: qrterminal.WHITE,
		QuietZone: 1,
	}

	// Gerar e exibir QR code no terminal
	qrterminal.GenerateWithConfig(code, config)
}
