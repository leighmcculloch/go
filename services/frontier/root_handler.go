package main

import "github.com/gofiber/fiber/v2"

type RootHandler struct {
	ArchiveURL string
	Data       *Data
}

type rootResponse struct {
	CoreVersion       string `json:"core_version"`
	CoreLatestLedger  uint32 `json:"core_latest_ledger"`
	NetworkPassphrase string `json:"network_passphrase"`
}

func (h *RootHandler) Handler(c *fiber.Ctx) error {
	resp := rootResponse{
		CoreVersion:       h.Data.RootHAS().Server,
		CoreLatestLedger:  h.Data.RootHAS().CurrentLedger,
		NetworkPassphrase: h.Data.RootHAS().NetworkPassphrase,
	}
	return c.JSON(resp)
}
