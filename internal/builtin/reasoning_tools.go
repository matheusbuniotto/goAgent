package builtin

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/matheusbuniotto/goagent/pkg/toolkit"
)

// ::: Ferramenta: Análise de Raciocínio :::

type AnalyzeReasoningInput struct {
	Problem   string `json:"problem"`
	Approach  string `json:"approach"`
	Confidence int   `json:"confidence"` // 1-10
}

func analyzeReasoning(input json.RawMessage) (string, error) {
	var typedInput AnalyzeReasoningInput
	if err := json.Unmarshal(input, &typedInput); err != nil {
		return "", fmt.Errorf("JSON inválido para argumentos: %w", err)
	}

	if typedInput.Problem == "" {
		return "", fmt.Errorf("argumento inválido. 'problem' é obrigatório")
	}

	// Análise estruturada do raciocínio
	analysis := fmt.Sprintf(`🧠 ANÁLISE DE RACIOCÍNIO
⏰ Timestamp: %s

📋 PROBLEMA ANALISADO:
%s

🎯 ABORDAGEM PROPOSTA:
%s

📊 CONFIANÇA: %d/10 %s

🔍 VALIDAÇÃO AUTOMÁTICA:
✅ Problema claramente definido: %s
✅ Abordagem estruturada: %s
✅ Nível de confiança adequado: %s

💡 RECOMENDAÇÕES:
%s`,
		time.Now().Format("15:04:05"),
		formatText(typedInput.Problem),
		formatText(typedInput.Approach),
		typedInput.Confidence,
		getConfidenceEmoji(typedInput.Confidence),
		boolToCheckmark(len(typedInput.Problem) > 10),
		boolToCheckmark(len(typedInput.Approach) > 20),
		boolToCheckmark(typedInput.Confidence >= 6),
		generateRecommendations(typedInput))

	return analysis, nil
}

func formatText(text string) string {
	if text == "" {
		return "❌ Não fornecido"
	}
	// Adiciona indentação para melhor formatação
	lines := strings.Split(text, "\n")
	var formatted []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			formatted = append(formatted, "  "+strings.TrimSpace(line))
		}
	}
	return strings.Join(formatted, "\n")
}

func getConfidenceEmoji(confidence int) string {
	switch {
	case confidence >= 9:
		return "🔥 (Muito Alta)"
	case confidence >= 7:
		return "✅ (Alta)"
	case confidence >= 5:
		return "⚠️ (Média)"
	case confidence >= 3:
		return "❓ (Baixa)"
	default:
		return "❌ (Muito Baixa)"
	}
}

func boolToCheckmark(condition bool) string {
	if condition {
		return "✅ Sim"
	}
	return "❌ Não"
}

func generateRecommendations(input AnalyzeReasoningInput) string {
	var recommendations []string

	if len(input.Problem) < 10 {
		recommendations = append(recommendations, "• Defina o problema com mais detalhes")
	}

	if len(input.Approach) < 20 {
		recommendations = append(recommendations, "• Elabore melhor a abordagem proposta")
	}

	if input.Confidence < 6 {
		recommendations = append(recommendations, "• Considere buscar mais informações antes de prosseguir")
	}

	if input.Confidence > 9 {
		recommendations = append(recommendations, "• Verifique se não há overconfidence - considere possíveis falhas")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "• Análise está bem estruturada, pode prosseguir com confiança")
	}

	return strings.Join(recommendations, "\n")
}

var AnalyzeReasoningDef = toolkit.ToolDefinition{
	Name:        "analyze_reasoning",
	Description: `Analisa e valida o próprio processo de raciocínio. Use quando quiser verificar se sua abordagem está bem fundamentada. Requer JSON com "problem", "approach" e "confidence" (1-10). Exemplo: {"problem": "Preciso implementar X", "approach": "Vou usar Y porque Z", "confidence": 7}`,
	Function:    analyzeReasoning,
}

// ::: Ferramenta: Revisão de Decisão :::

type ReviewDecisionInput struct {
	Decision    string   `json:"decision"`
	Factors     []string `json:"factors"`
	Alternatives []string `json:"alternatives"`
}

func reviewDecision(input json.RawMessage) (string, error) {
	var typedInput ReviewDecisionInput
	if err := json.Unmarshal(input, &typedInput); err != nil {
		return "", fmt.Errorf("JSON inválido para argumentos: %w", err)
	}

	if typedInput.Decision == "" {
		return "", fmt.Errorf("argumento inválido. 'decision' é obrigatório")
	}

	review := fmt.Sprintf(`🎯 REVISÃO DE DECISÃO
⏰ Timestamp: %s

💼 DECISÃO:
%s

📊 FATORES CONSIDERADOS:
%s

🔄 ALTERNATIVAS AVALIADAS:
%s

🔍 ANÁLISE CRÍTICA:
• Completude dos fatores: %s
• Consideração de alternativas: %s
• Clareza da decisão: %s

⭐ SCORE GERAL: %s`,
		time.Now().Format("15:04:05"),
		formatText(typedInput.Decision),
		formatFactors(typedInput.Factors),
		formatAlternatives(typedInput.Alternatives),
		evaluateCompleteness(typedInput.Factors),
		evaluateAlternatives(typedInput.Alternatives),
		evaluateClarity(typedInput.Decision),
		calculateOverallScore(typedInput))

	return review, nil
}

func formatFactors(factors []string) string {
	if len(factors) == 0 {
		return "❌ Nenhum fator especificado"
	}
	var formatted []string
	for i, factor := range factors {
		formatted = append(formatted, fmt.Sprintf("  %d. %s", i+1, factor))
	}
	return strings.Join(formatted, "\n")
}

func formatAlternatives(alternatives []string) string {
	if len(alternatives) == 0 {
		return "❌ Nenhuma alternativa considerada"
	}
	var formatted []string
	for i, alt := range alternatives {
		formatted = append(formatted, fmt.Sprintf("  %d. %s", i+1, alt))
	}
	return strings.Join(formatted, "\n")
}

func evaluateCompleteness(factors []string) string {
	count := len(factors)
	switch {
	case count >= 3:
		return "✅ Boa (3+ fatores)"
	case count >= 2:
		return "⚠️ Razoável (2 fatores)"
	case count == 1:
		return "❓ Limitada (1 fator)"
	default:
		return "❌ Insuficiente (0 fatores)"
	}
}

func evaluateAlternatives(alternatives []string) string {
	count := len(alternatives)
	switch {
	case count >= 2:
		return "✅ Boa (2+ alternativas)"
	case count == 1:
		return "⚠️ Limitada (1 alternativa)"
	default:
		return "❌ Insuficiente (0 alternativas)"
	}
}

func evaluateClarity(decision string) string {
	length := len(decision)
	switch {
	case length >= 50:
		return "✅ Clara e detalhada"
	case length >= 20:
		return "⚠️ Razoavelmente clara"
	default:
		return "❌ Muito vaga"
	}
}

func calculateOverallScore(input ReviewDecisionInput) string {
	score := 0
	
	// Pontuação baseada em fatores
	if len(input.Factors) >= 3 {
		score += 3
	} else if len(input.Factors) >= 2 {
		score += 2
	} else if len(input.Factors) == 1 {
		score += 1
	}
	
	// Pontuação baseada em alternativas
	if len(input.Alternatives) >= 2 {
		score += 2
	} else if len(input.Alternatives) == 1 {
		score += 1
	}
	
	// Pontuação baseada na clareza
	if len(input.Decision) >= 50 {
		score += 2
	} else if len(input.Decision) >= 20 {
		score += 1
	}
	
	total := score
	maxScore := 7
	
	percentage := (total * 100) / maxScore
	
	switch {
	case percentage >= 85:
		return fmt.Sprintf("🏆 Excelente (%d/%d - %d%%)", total, maxScore, percentage)
	case percentage >= 70:
		return fmt.Sprintf("✅ Boa (%d/%d - %d%%)", total, maxScore, percentage)
	case percentage >= 50:
		return fmt.Sprintf("⚠️ Razoável (%d/%d - %d%%)", total, maxScore, percentage)
	default:
		return fmt.Sprintf("❌ Precisa melhorar (%d/%d - %d%%)", total, maxScore, percentage)
	}
}

var ReviewDecisionDef = toolkit.ToolDefinition{
	Name:        "review_decision",
	Description: `Revisa criticamente uma decisão tomada, analisando fatores e alternativas. Use para validar decisões importantes. Requer JSON com "decision", "factors" (array) e "alternatives" (array). Exemplo: {"decision": "Vou usar X", "factors": ["performance", "custo"], "alternatives": ["usar Y", "manter Z"]}`,
	Function:    reviewDecision,
}