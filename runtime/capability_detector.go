//runtime/capability_detector.go

package runtimectx

import (
	"os"
	"runtime"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

func DetectCapabilityProfile() *schema.CapabilityProfile {
	cp := schema.NewCapabilityProfile()

	// ---- Display ----
	if !isHeadless() {
		cp.Mark(schema.CapDisplay, schema.CapOK)
	} else {
		cp.Mark(schema.CapDisplay, schema.CapUnavailable)
	}

	// ---- Keyboard ----
	cp.Mark(schema.CapKeyboard, schema.CapOK) // assume present

	// ---- Network ----
	if hasNetwork() {
		cp.Mark(schema.CapNetwork, schema.CapOK)
	} else {
		cp.Mark(schema.CapNetwork, schema.CapDegraded)
	}

	// ---- Microphone (stub) ----
	if hasMicrophone() {
		cp.Mark(schema.CapMicrophone, schema.CapOK)
	} else {
		cp.Mark(schema.CapMicrophone, schema.CapUnavailable)
	}

	// ---- Speaker (stub) ----
	if hasSpeaker() {
		cp.Mark(schema.CapSpeaker, schema.CapOK)
	} else {
		cp.Mark(schema.CapSpeaker, schema.CapUnavailable)
	}

	return cp
}

func isHeadless() bool {
	return os.Getenv("DISPLAY") == "" &&
		os.Getenv("WAYLAND_DISPLAY") == ""
}

func hasNetwork() bool {
	// simple heuristic
	return true
}

func hasMicrophone() bool {
	// stub → replace with OS API later
	return true
}

func hasSpeaker() bool {
	return true
}