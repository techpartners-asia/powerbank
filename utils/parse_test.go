package powerbankUtils

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/techpartners-asia/powerbank/constants"
	powerbankModels "github.com/techpartners-asia/powerbank/models"
)

// fromHex converts a space-delimited hex string like "A8 00 0C 31" to bytes.
func fromHex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(strings.ReplaceAll(s, " ", ""))
	if err != nil {
		t.Fatalf("decode hex: %v", err)
	}
	return b
}

// TestParsePopupPowerBankResponse uses the example frame from
// docs.volinks.com/powerbank-protocol-v1/en/guide/protocol-popupsn.html
func TestParsePopupPowerBankResponse(t *testing.T) {
	payload := fromHex(t, "A8 00 0C 31 60 00 9B D2 10 01 00 3B")

	got, err := ParsePopupPowerBankResponse(payload)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	want := &powerbankModels.PowerBankPopupResponse{
		Head:        0xA8,
		Length:      0x000C,
		Cmd:         0x31,
		HoleIndex:   0x60,
		PowerbankSN: "10211856",
		State:       0x01,
		Reserved:    0x00,
		Verify:      0x3B,
	}

	if *got != *want {
		t.Errorf("popup mismatch:\n got: %+v\nwant: %+v", got, want)
	}
}

// TestParseCheckResponse uses the 137-byte example from
// docs.volinks.com/powerbank-protocol-v1/en/guide/protocol-check.html
func TestParseCheckResponse(t *testing.T) {
	payload := fromHex(t, "A8 00 89 10 01 FF FF 00 04 16 01 01 00 EC 00 05 11 49 F1 64 1F 32 01 0D 00 02 00 00 00 00 00 00 00 00 00 00 00 00 00 80 03 01 00 E8 00 05 11 46 AC 64 20 32 00 0D 00 04 00 00 00 00 00 00 00 00 00 00 00 00 00 80 02 FF FF 00 04 16 05 01 00 D7 00 04 C6 F0 96 64 1F 32 00 1A 00 06 00 00 00 00 00 00 00 00 00 00 00 00 00 80 07 01 00 E9 00 05 11 49 DB 64 1E 32 00 0D 00 08 00 00 00 00 00 00 00 00 00 00 00 00 00 80 D8")

	got, err := ParseCheckResponse(payload)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if got.Head != 0xA8 || got.Cmd != 0x10 || got.Length != 0x0089 {
		t.Errorf("header mismatch: head=%#x cmd=%#x len=%d", got.Head, got.Cmd, got.Length)
	}
	if got.Verify != 0xD8 {
		t.Errorf("verify: got %#x want 0xD8", got.Verify)
	}
	if len(got.ControlBoards) != 2 {
		t.Fatalf("expected 2 control boards, got %d", len(got.ControlBoards))
	}

	// Board 1: 01 FF FF 00 04 16, then 4 holes.
	b1 := got.ControlBoards[0]
	if b1.ControlIndex != 1 || b1.SoftVersion != 4 || b1.HardVersion != 0x16 {
		t.Errorf("board1 header: %+v", b1)
	}
	if len(b1.Holes) != 4 {
		t.Fatalf("board1: expected 4 holes, got %d", len(b1.Holes))
	}

	// Hole 1 of board 1: 01 01 00 EC 00 05 11 49 F1 64 1F 32 01 0D 00
	h := b1.Holes[0]
	if h.HoleIndex != 1 || h.State != 1 {
		t.Errorf("hole1: index=%d state=%d", h.HoleIndex, h.State)
	}
	if h.PowerbankSN != "85019121" {
		t.Errorf("hole1 SN: got %q want %q (doc example 0x05 0x11 0x49 0xF1)", h.PowerbankSN, "85019121")
	}
	if h.SOC != 100 {
		t.Errorf("hole1 SOC: got %d want 100", h.SOC)
	}
	if h.Temperature != 0x1F {
		t.Errorf("hole1 temp: got %d want %d", h.Temperature, 0x1F)
	}
	if h.Sensor != 0x00 {
		t.Errorf("hole1 sensor: got %#x want 0x00", h.Sensor)
	}

	// Hole 2 of board 1: 02 00 00 00 00 00 00 00 00 00 00 00 00 00 80 — empty slot.
	h2 := b1.Holes[1]
	if h2.HoleIndex != 2 || h2.State != 0 || h2.Sensor != 0x80 {
		t.Errorf("hole2 (empty slot): %+v", h2)
	}
}

func TestParseResponseDispatch(t *testing.T) {
	cases := []struct {
		name    string
		hex     string
		wantTyp constants.PUBLISH_TYPE
	}{
		{"popup_sn", "A8 00 0C 31 60 00 9B D2 10 01 00 3B", constants.PUBLISH_TYPE_POPUP},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			typ, res, err := ParseResponse(fromHex(t, tc.hex))
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			if typ != tc.wantTyp {
				t.Errorf("type: got %q want %q", typ, tc.wantTyp)
			}
			if res == nil {
				t.Errorf("nil response struct")
			}
		})
	}
}

// TestParseHealthCheckResponse uses the example frame from
// docs.volinks.com/powerbank-protocol-v1/en/guide/protocol-heart.html
// "A8 00 11 7A 10 43 53 51 3A 32 37 3B 42 50 3A 30 FC" → signal "CSQ:27;BP:0"
func TestParseHealthCheckResponse(t *testing.T) {
	payload := fromHex(t, "A8 00 11 7A 10 43 53 51 3A 32 37 3B 42 50 3A 30 FC")

	got, err := ParseHealthCheckResponse(payload)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if got.Head != 0xA8 || got.Cmd != 0x7A || got.Length != 0x0011 {
		t.Errorf("header mismatch: head=%#x cmd=%#x len=%d", got.Head, got.Cmd, got.Length)
	}
	if got.ControlIndex != 0x10 {
		t.Errorf("controlIndex: got %d want 16", got.ControlIndex)
	}
	if got.Signal != "CSQ:27;BP:0" {
		t.Errorf("signal: got %q want %q", got.Signal, "CSQ:27;BP:0")
	}
	if got.Verify != 0xFC {
		t.Errorf("verify: got %#x want 0xFC", got.Verify)
	}
	if got.GetBackupPowerStatus() != 0 {
		t.Errorf("BP: got %d want 0", got.GetBackupPowerStatus())
	}
}

// TestParsePopupResponseStates verifies the popup_sn state mapping covers
// every state byte enumerated in protocol-popupsn.html.
func TestParsePopupResponseStates(t *testing.T) {
	cases := []struct {
		state     byte
		wantStat  constants.PowerbankStatus
		wantDescr string
	}{
		{0x00, constants.PowerbankStatus_PopupFailed, "Pop-up failed"},
		{0x01, constants.PowerbankStatus_PopupSuccessful, "Pop-up successful"},
		{0x11, constants.PowerbankStatus_PopupSerialTimeout, "Serial communication timeout"},
		{0x12, constants.PowerbankStatus_PopupBankUnpoppedSnReadable, "Power bank has not popped out, but the SN is readable"},
		{0x87, constants.PowerbankStatus_PopupTimestampRetrievalFailed, "Failed to obtain the timestamp"},
		{0x88, constants.PowerbankStatus_PopupTTLExceeded, "Exceeded the TTL validity period"},
		{0xFB, constants.PowerbankStatus_PopupNoMatchingBattery, "No portable charger meets the rental requirements"},
		{0xFC, constants.PowerbankStatus_PopupTargetSnNotFound, "Target SN not found among charging batteries"},
		{0xFD, constants.PowerbankStatus_PopupAddTaskFailed, "Failed to add task to the thread pool"},
		{0xFE, constants.PowerbankStatus_PopupPreviousRentalIncomplete, "Previous rental not completed; new rental cannot start"},
		{0xFF, constants.PowerbankStatus_PopupCommandParsingFailed, "Lease command parsing failed"},
		{0x7A, constants.PowerbankStatus_UnknownError, "Unknown error"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("state_%#x", tc.state), func(t *testing.T) {
			r := &powerbankModels.PowerBankPopupResponse{State: int(tc.state)}
			if got := r.GetStatus(); got != tc.wantStat {
				t.Errorf("GetStatus: got %q want %q", got, tc.wantStat)
			}
			if got := r.GetDescription(); got != tc.wantDescr {
				t.Errorf("GetDescription: got %q want %q", got, tc.wantDescr)
			}
		})
	}
}

func TestParseResponseShortPayload(t *testing.T) {
	if _, _, err := ParseResponse([]byte{0xA8, 0x00}); err == nil {
		t.Errorf("expected error for short payload")
	}
}

func TestParseResponseUnknownCommand(t *testing.T) {
	// Valid header but unknown cmd byte 0xAB.
	if _, _, err := ParseResponse([]byte{0xA8, 0x00, 0x05, 0xAB, 0x00}); err == nil {
		t.Errorf("expected error for unknown cmd")
	}
}
