package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// LogEntry representa uma entrada de log no formato JSON
type LogEntry struct {
	Time     string `json:"time"`
	Level    string `json:"level"`
	Message  string `json:"msg"`
	Package  string `json:"package"`
	Function string `json:"function"`
	File     string `json:"file"`
	Line     uint   `json:"line"`
}

// Cores ANSI para terminal
const (
	Reset    = "\033[0m"
	Bold     = "\033[1m"
	Red      = "\033[31m"
	Green    = "\033[32m"
	Yellow   = "\033[33m"
	Blue     = "\033[34m"
	Magenta  = "\033[35m"
	Cyan     = "\033[36m"
	White    = "\033[37m"
	BgRed    = "\033[41m"
	BgGreen  = "\033[42m"
	BgYellow = "\033[43m"
	BgBlue   = "\033[44m"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		processLogLine(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler entrada: %v\n", err)
		os.Exit(1)
	}
}

func processLogLine(line string) {
	var entry LogEntry

	// Verifica se a linha não está vazia
	if strings.TrimSpace(line) == "" {
		return
	}

	// Trata linhas truncadas ou inválidas
	if !strings.HasSuffix(line, "}") {
		line = line + "}"
	}

	err := json.Unmarshal([]byte(line), &entry)
	if err != nil {
		fmt.Printf("%sErro ao processar linha: %s%s\n", Yellow, line, Reset)
		return
	}

	// Formata o timestamp
	timestamp := formatTimestamp(entry.Time)

	// Determina a cor com base no nível de log
	levelColor := getLevelColor(entry.Level)

	// Imprime a linha formatada
	fmt.Printf("%s %s [%s %s] %s (%s)\n",
		Bold+White+timestamp+Reset,
		levelColor+entry.Level+Reset,

		Cyan+entry.Package+Reset,
		formatFunctionName(entry.Function),

		levelColor+entry.Message+Reset,

		fmt.Sprintf("%s:%d%s", Magenta+entry.File, entry.Line, Reset),
	)
}

func formatTimestamp(timeStr string) string {
	t, err := time.Parse(time.RFC3339Nano, timeStr)
	if err != nil {
		return timeStr
	}
	return t.Format("15:04:05.000")
}

func formatFunctionName(function string) string {
	// Extrai apenas o nome da função sem o pacote
	parts := strings.Split(function, ".")
	if len(parts) > 0 {
		return Blue + parts[len(parts)-1] + Reset
	}
	return Blue + function + Reset
}

func getLevelColor(level string) string {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return White
	case "INFO":
		return Green
	case "WARN":
		return Yellow
	case "ERROR":
		return Red
	case "FATAL", "PANIC":
		return BgRed + Bold + White
	case "TRACE":
		return Cyan
	default:
		return Reset
	}
}
