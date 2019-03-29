package inits

import "github.com/zekroTJA/yuri2/internal/logger"

// InitLogger will initialize the logger.
func InitLogger() {
	logger.Setup(`%{color}▶  %{level:.4s} %{id:03d}%{color:reset} %{message}`, 5)
}
