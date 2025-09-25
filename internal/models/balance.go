package models

import (
	"encoding/json"
	"fmt"
)

// Float64 es un tipo personalizado que siempre serializa con 2 decimales
type Float64 float64

// MarshalJSON implementa la serialización JSON con 2 decimales
func (f Float64) MarshalJSON() ([]byte, error) {
	// Forzar siempre 2 decimales, incluso para enteros
	return []byte(fmt.Sprintf("%.2f", float64(f))), nil
}

// UnmarshalJSON implementa la deserialización JSON
func (f *Float64) UnmarshalJSON(data []byte) error {
	var val float64
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	*f = Float64(val)
	return nil
}

// BalanceInfo representa la información de balance de un usuario
type BalanceInfo struct {
	Balance      Float64 `json:"balance"`
	TotalDebits  Float64 `json:"total_debits"`
	TotalCredits Float64 `json:"total_credits"`
}
