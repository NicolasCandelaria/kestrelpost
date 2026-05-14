package game

// ActForNight maps nights 1–9 to act / “level” 1–3.
func ActForNight(night int) int {
	switch {
	case night <= 3:
		return 1
	case night <= 6:
		return 2
	default:
		return 3
	}
}

// ActTitle is a short banner for the act (terminal “level”).
func ActTitle(act int) string {
	switch act {
	case 1:
		return "LEVEL I · STOCK AND SURGE"
	case 2:
		return "LEVEL II · STATIC AND SHADOW"
	case 3:
		return "LEVEL III · CONVOY WINDOW"
	default:
		return "LEVEL"
	}
}

// Choice is one keyed response for a night (1–3). Reply is what the player sees;
// numeric fields drive resolution behind the scenes and are not shown in the UI.
type Choice struct {
	Reply        string
	Fuel         int
	Hub          int
	Trust        int
	Kid          int
	SetHarrow    bool
	OseiRelease  bool
	ConvoyBetray bool
}

// NightCard is the incoming beat for a single night.
type NightCard struct {
	Act     int
	Source  string
	Hash    string
	Quote   string
	Choices [3]Choice
}

// NightScript returns the transmission card for the given night (1–9).
func NightScript(night int) NightCard {
	switch night {
	case 1:
		return NightCard{1, "MAREN", "9f2c",
			"Maren breaks squelch twice before she trusts the open carrier.\n\n" +
				"“Operator Seven—Kestrel. I’ve got eight in the hall. Two fevers climbing, one kid who won’t take water. The clinic net is rumor soup.\n" +
				"I need a line on who still has power north of the washout—not churches, not ghost towns. A place that can take blood work if we move at dawn.”\n\n" +
				"Under the hash you hear gym mats dragged, a kettle clicking, someone counting breaths like inventory.",
			[3]Choice{
				{Reply: "Stay with her. Walk the rail spur triangle aloud—triage site, diesel on the wind, the substation that looks fine but isn’t. Make her read the map back before you let go of the key.", Fuel: 18, Hub: 2, Trust: 1},
				{Reply: "Keep it short: one bearing, one landmark, one rule—if the lights die, the fridge wins. Then cut so she can work.", Fuel: 10, Hub: 1, Trust: 1},
				{Reply: "Send a single tone and your callsign only. You’ll patch the rest when the band quiets—even if she thinks you’ve gone cold.", Fuel: 3, Hub: 0, Trust: -1},
			}}
	case 2:
		return NightCard{1, "MAREN", "9f2c",
			"Her voice is thinner. A generator coughs in the background like an animal.\n\n" +
				"“If I run the hall lights past midnight I lose the vaccine fridge. I can’t do poetry right now, Seven—tell me straight. Heat or light. Which one gets the knife?”\n\n" +
				"You picture the hall: breath fogging, tape on the floor, parents trading shifts with their coats still on.",
			[3]Choice{
				{Reply: "Talk her through a heat curfew: where bodies sleep, where doors seal, how long the fridge can hold if nobody opens it for vanity. Make the choice feel owned, not imposed.", Fuel: 20, Hub: 2, Trust: 1},
				{Reply: "One sentence only: fridge first. Then silence—trust that she’ll translate fear into procedure without you narrating it.", Fuel: 10, Hub: 1, Trust: 1},
				{Reply: "Tell her you need morning weather and band noise before you commit. It buys you minutes; it also leaves her alone with the thermostat.", Fuel: 4, Hub: 0, Trust: -1},
			}}
	case 3:
		return NightCard{1, "MAREN", "9f2c",
			"“Someone’s painting our band,” she says. “Not jamming—knocking. Like they’re testing whether Kestrel answers soft or loud.”\n\n" +
				"I can go dark. I can answer. I can flip to your old training net and pretend I’m you for sixty seconds and probably regret it.\n" +
				"What’s the move that doesn’t turn this hall into a lighthouse for every hungry ear between here and Hudson?”",
			[3]Choice{
				{Reply: "Give her a tight listen window and a safe CQ pattern—enough to learn who’s knocking without teaching them your grammar.", Fuel: 22, Hub: 2, Trust: 1},
				{Reply: "Procedure voice: twenty minutes listen-only, log only, no hero keys. If it’s hostile, she’ll feel it without you dramatizing.", Fuel: 10, Hub: 1, Trust: 0},
				{Reply: "Tell her to go dark and stay dark until you pattern-match from Kestrel. No net, no comfort—just survival and blame if it goes wrong.", Fuel: 3, Hub: 0, Trust: -1},
			}}
	case 4:
		return NightCard{2, "KID", "—",
			"A voice you don’t have a file for. Young. Too calm.\n\n" +
				"“Kestrel,” they say, like they’ve been saving the name. “Do the geese still cross the river at dawn? I’m counting.”\n\n" +
				"It sounds childish until you hear the second carrier underneath—clean, stable, wrong—like someone borrowed a nursery rhyme to tune an instrument.",
			[3]Choice{
				{Reply: "Answer honestly about the birds and the light. Then ask, quietly, what they’re actually counting. Let the silence do part of the work.", Fuel: 16, Hub: 0, Trust: 0, Kid: 2},
				{Reply: "Give them the plain seasonal truth—yes or no—nothing else. No warmth, no thread for them to pull.", Fuel: 10, Hub: 0, Trust: 0, Kid: 1},
				{Reply: "Clip the carrier. Refuse the riddle. If they’re fishing, let them fish in empty water—even if it feels cruel from where you sit.", Fuel: 4, Hub: 0, Trust: -1, Kid: -1},
			}}
	case 5:
		return NightCard{2, "MAREN", "9f2c",
			"Maren doesn’t greet you. She reads a string like liturgy.\n\n" +
				"“North road signature yesterday—your log called it civ. I need that to still be true tonight, because I’ve got people asking if they should run toward it.\n" +
				"If you’re wrong, say you’re wrong now. If you’re guessing, say you’re guessing. I can’t sell ‘maybe’ to a room that already buried someone today.”",
			[3]Choice{
				{Reply: "Lay the hashes side by side—what you saw, what could mimic it, what would change your mind. End with what you would bet a life on.", Fuel: 22, Hub: 2, Trust: 1},
				{Reply: "One line: civ pattern, as far as Kestrel knows. No flourishes. Enough for her to steer the room.", Fuel: 10, Hub: 1, Trust: 1},
				{Reply: "Tell her you can’t confirm—not ‘no,’ not ‘yes,’ just fog. It keeps you clean on paper. It also leaves her holding the bag alone.", Fuel: 5, Hub: 1, Trust: -2},
			}}
	case 6:
		return NightCard{2, "HARROW", "b81d",
			"Post Harrow rides in narrow, professional, like a blade slid under a door.\n\n" +
				"“You’re lighting up the commons like a Christmas tree, Kestrel. I can share a roster—two nets, one schedule, one lie we tell the sky together.\n" +
				"Or we can keep pretending charity scales. Pick. I’m not here to be your villain; I’m here so neither of us cooks off the tower.”",
			[3]Choice{
				{Reply: "Take the roster. Split the bands. Promise the ugly coordination work so Harrow’s people stop treating your survivors like interference.", Fuel: 18, Hub: 1, Trust: 0, SetHarrow: true},
				{Reply: "Offer shared listen windows but refuse to shrink your net. Cooperation without surrender—knowing it may cost you both heat and sleep.", Fuel: 10, Hub: 1, Trust: 0},
				{Reply: "Decline cold. No rationale on air. Let the carrier close on silence and live with what Harrow writes about you afterward.", Fuel: 4, Hub: 0, Trust: -1},
			}}
	case 7:
		return NightCard{3, "COLE", "44aa",
			"Cole doesn’t waste callsigns. Twelve voices behind him, one truck idling you can almost smell through static.\n\n" +
				"“Corridor truth,” he says. “Not myth. Not poetry. I get a corridor that still exists, you get me forgetting I ever heard a woman named Maren on your lip.\n" +
				"Fair trade. Unfair world. Speak.”",
			[3]Choice{
				{Reply: "Give him a corridor that was true once—stale enough to buy the hall time, sharp enough he believes you’re not bluffing.", Fuel: 20, Hub: 2, Trust: -1},
				{Reply: "Burn his clock: plausible junk, plausible doubt, a story that moves his truck an hour the wrong way without firing a shot you can name.", Fuel: 10, Hub: 0, Trust: 1},
				{Reply: "Refuse. Let the threat hang in the air. Let the hall hear what silence costs when you won’t negotiate with a loaded engine.", Fuel: 6, Hub: -1, Trust: 0},
			}}
	case 8:
		return NightCard{3, "DR. OSEI", "loop",
			"The loop arrives in shards—tape hiss, classroom cough, a voice trying to stay composed while the world unthreads.\n\n" +
				"“…solar class… not disclosed… cascade… personnel…”\n\n" +
				"You could clean it, normalize it, let it run clean once so people understand what happened—or keep it in the drawer and decide who deserves the weight tonight.",
			[3]Choice{
				{Reply: "Stabilize the audio, strip the worst of the hash noise, and let the whole clip ride every open net you still own until the sky changes color.", Fuel: 24, Hub: 0, Trust: 0, OseiRelease: true},
				{Reply: "Archive the raw locally. Send Maren a redacted slice—enough to change her night, not enough to paint a target on her roof.", Fuel: 10, Hub: 2, Trust: 1},
				{Reply: "Checksum it, log it, do not forward. You become the only spine that remembers—until memory becomes its own kind of violence.", Fuel: 5, Hub: 0, Trust: 0},
			}}
	case 9:
		return NightCard{3, "MAREN / SKY", "9f2c",
			"Maren sounds like someone standing at a window while the horizon misbehaves.\n\n" +
				"“Hall’s full. Cole’s truck hash is on the road. Harrow’s carrier just winked at me like a neighbor who knows my schedule.\n" +
				"Seven—who gets your last good hour? Me, them, or yourself? Don’t be brave. Be accurate.”",
			[3]Choice{
				{Reply: "Give her the last hour: maps, prefixes, the blunt truth about what still answers when Kestrel goes quiet. You spend yourself so the hall can keep a door cracked.", Fuel: 26, Hub: 3, Trust: 2},
				{Reply: "Split the hour into slices—fifteen minutes for her net, fifteen for yours, fifteen for the air you refuse to own. Nobody gets a crown; everybody gets a chance.", Fuel: 10, Hub: 2, Trust: 1},
				{Reply: "Name the hall’s grid and the cache path. Tell her extraction is inbound and not to trust any other voice until morning—then live inside what that promise costs.", Fuel: 8, Hub: -2, Trust: -3, ConvoyBetray: true},
			}}
	default:
		return NightScript(1)
	}
}
