package ending

// ResolveEnding picks exactly one ending; first matching rule wins (spec order).
func ResolveEnding(cfg Config, s RunState) Ending {
	if s.ConvoyBetrayal {
		return TheConvoy
	}
	if s.OseiFullRelease {
		return FullBroadcast
	}
	if s.HarrowDarkPlan {
		return DarkFrequency
	}
	if s.KidInvestigation >= cfg.KMax && !s.OseiFullRelease {
		return TheKidWasRight
	}
	if s.Fuel <= 0 && s.TerminalDarkNight > 0 && s.TerminalDarkNight < cfg.DeadAirExclusiveMaxTerminalNight && !s.ConvoyBetrayal {
		return DeadAir
	}
	if s.Fuel <= 0 &&
		s.TerminalDarkNight >= cfg.RelayMinTerminalNight &&
		s.MarenHubSupport >= cfg.MThreshold &&
		s.MarenTrust >= cfg.TThreshold {
		return TheRelay
	}
	return Fallback
}
