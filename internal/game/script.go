package game

// ActForNight maps nights 1–20 to act 1–4.
func ActForNight(night int) int {
	switch {
	case night <= 5:
		return 1
	case night <= 10:
		return 2
	case night <= 15:
		return 3
	default:
		return 4
	}
}

// ActTitle is a short banner for the act (terminal “level”).
func ActTitle(act int) string {
	switch act {
	case 1:
		return "ACT I · SETTLING IN"
	case 2:
		return "ACT II · PRESSURE BUILDING"
	case 3:
		return "ACT III · CONVERGENCE"
	case 4:
		return "ACT IV · ENDGAME"
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
			"Maren comes through after two failed calls. Her signal is weak, but her voice is steady.\n\n" +
				"“Kestrel, this is Maren at the community hall. I have eight people here. Two have high fevers, and one child cannot keep water down. The clinic channel is full of conflicting reports.\n" +
				"I need to know which site north of the washout still has power and medical supplies. If we move at dawn, I need the safest place to send them.”\n\n" +
				"You hear people coughing behind her and someone dragging tables across the floor.",
			[3]Choice{
				{Reply: "Keep Maren on the radio. Give her the route to the school gym, warn her about the flooded service road, and make her repeat the directions before anyone leaves.", Fuel: 18, Hub: 2, Trust: 1},
				{Reply: "Give her the nearest usable site and one clear landmark, then sign off so she can organize the room.", Fuel: 10, Hub: 1, Trust: 1},
				{Reply: "Send only your callsign and tell her to hold position until you can confirm more. It gives you time, but leaves her without a plan.", Fuel: 3, Hub: 0, Trust: -1},
			}}
	case 2:
		return NightCard{1, "MAREN", "9f2c",
			"Maren calls again near midnight. The hall generator is running rough.\n\n" +
				"“If I keep the heat and lights on all night, the vaccine fridge dies before morning. If I protect the fridge, the hall gets cold fast.\n" +
				"I need a decision I can defend to scared people. Heat for the room, or power for the medicine?”\n\n" +
				"The carrier catches a child crying, then a door slamming shut against the cold.",
			[3]Choice{
				{Reply: "Talk her through a ration plan: seal the gym, move everyone into one corner, run heat in short bursts, and keep the fridge closed unless medicine is needed.", Fuel: 20, Hub: 2, Trust: 1},
				{Reply: "Tell her to protect the medicine first and use coats, blankets, and body heat to get through the night.", Fuel: 10, Hub: 1, Trust: 1},
				{Reply: "Tell her you cannot call it yet and need another weather report. It buys time for Kestrel, but she has to face the room alone.", Fuel: 4, Hub: 0, Trust: -1},
			}}
	case 3:
		return NightCard{1, "MAREN", "9f2c",
			"Maren lowers her voice before she speaks.\n\n" +
				"“Someone is testing our channel. They are not blocking us. They are calling once, waiting, then calling again from a slightly different frequency.\n" +
				"If I answer, they learn we are here. If I go silent, I lose contact with anyone who might still help. What do you want me to do?”",
			[3]Choice{
				{Reply: "Give her a short test script: one reply, no location, no names, then listen for ten minutes and write down every response.", Fuel: 22, Hub: 2, Trust: 1},
				{Reply: "Order her to stay listen-only for twenty minutes and log each call without answering.", Fuel: 10, Hub: 1, Trust: 0},
				{Reply: "Tell her to power down the hall radio until you identify the caller from Kestrel.", Fuel: 3, Hub: 0, Trust: -1},
			}}
	case 4:
		return NightCard{2, "KID", "—",
			"A young voice cuts into the band. You do not recognize the callsign. The voice is calm in a way that makes the room feel smaller.\n\n" +
				"“Kestrel, do geese still cross the river at dawn? I am counting them.”\n\n" +
				"The question sounds harmless, but a second signal sits under the voice. It is too clean, too steady, and it follows every pause.",
			[3]Choice{
				{Reply: "Answer the question plainly, then ask where they are and who taught them to use this channel.", Fuel: 16, Hub: 0, Trust: 0, Kid: 2},
				{Reply: "Say only that geese cross at first light this time of year, then end the exchange.", Fuel: 10, Hub: 0, Trust: 0, Kid: 1},
				{Reply: "Cut the transmission and mark the frequency as suspicious.", Fuel: 4, Hub: 0, Trust: -1, Kid: -1},
			}}
	case 5:
		return NightCard{2, "MAREN", "9f2c",
			"Maren comes back with no greeting.\n\n" +
				"“Yesterday you logged a civilian truck convoy north of the washout. People here want to leave the hall and try to meet it.\n" +
				"If that convoy is real, say so. If you only guessed, say that too. I cannot send families into the dark on a maybe.”",
			[3]Choice{
				{Reply: "Read her the log: time, direction, signal strength, and why you think it was civilian. Then tell her what would prove you wrong.", Fuel: 22, Hub: 2, Trust: 1},
				{Reply: "Tell her the convoy looked civilian based on the last clean signal, but you cannot guarantee what it is now.", Fuel: 10, Hub: 1, Trust: 1},
				{Reply: "Tell her you cannot confirm the convoy and will not recommend moving anyone toward it.", Fuel: 5, Hub: 1, Trust: -2},
			}}
	case 6:
		return NightCard{2, "HARROW", "b81d",
			"Post Harrow breaks in with a clean signal and a trained operator.\n\n" +
				"“Kestrel, your open calls are drawing attention. I can share our contact roster and split the traffic with you.\n" +
				"You keep your hall. We keep ours. We stop broadcasting over each other before the wrong people map both sites.”",
			[3]Choice{
				{Reply: "Accept Harrow’s roster and divide the channels so both shelters can transmit without stepping on each other.", Fuel: 18, Hub: 1, Trust: 0, SetHarrow: true},
				{Reply: "Offer scheduled listening windows, but refuse to give Harrow control of Maren’s traffic.", Fuel: 10, Hub: 1, Trust: 0},
				{Reply: "Decline the offer and close the channel before Harrow can press for details.", Fuel: 4, Hub: 0, Trust: -1},
			}}
	case 7:
		return NightCard{3, "COLE", "44aa",
			"Cole forces his way onto the channel. You hear a truck engine and several people talking behind him.\n\n" +
				"“I know you are guiding people from Kestrel. Give me a passable route through the highway blocks, and I forget I heard the name Maren.\n" +
				"Refuse, and I start looking for the hall myself.”",
			[3]Choice{
				{Reply: "Give Cole an old route that may still be open but sends him away from the hall.", Fuel: 20, Hub: 2, Trust: -1},
				{Reply: "Keep him talking with partial directions and road warnings, wasting his time while Maren prepares to move.", Fuel: 10, Hub: 0, Trust: 1},
				{Reply: "Refuse to bargain and warn Maren that Cole is listening for her location.", Fuel: 6, Hub: -1, Trust: 0},
			}}
	case 8:
		return NightCard{3, "DR. OSEI", "loop",
			"An old recording repeats on an emergency frequency. It is damaged, but the speaker identifies herself as Dr. Osei.\n\n" +
				"“…solar event larger than disclosed… grid cascade already underway… personnel advised not to travel…”\n\n" +
				"If you rebroadcast the file, shelters may finally understand what happened. They may also panic, move too soon, or draw attention by calling for answers.",
			[3]Choice{
				{Reply: "Clean the audio enough to understand and rebroadcast the full recording on every channel you can reach.", Fuel: 24, Hub: 0, Trust: 0, OseiRelease: true},
				{Reply: "Save the full recording, but send Maren only the parts that help her make decisions tonight.", Fuel: 10, Hub: 2, Trust: 1},
				{Reply: "Log the recording and keep it off the air until you know who else is listening.", Fuel: 5, Hub: 0, Trust: 0},
			}}
	case 9:
		return NightCard{3, "MAREN / SKY", "9f2c",
			"Maren’s final call is clear enough that you can hear wind against the hall doors.\n\n" +
				"“The hall is full. Cole is on the road. Harrow just called close enough to know our schedule.\n" +
				"If Kestrel has one more clean hour, I need to know where you are putting it. Do we move the hall, split the network, or trust the convoy message that just came in?”",
			[3]Choice{
				{Reply: "Use the hour for Maren: confirm the safest exit route, mark the hostile frequencies, and give her the last verified shelter list.", Fuel: 26, Hub: 3, Trust: 2},
				{Reply: "Divide the hour between Maren, Harrow, and the open emergency channel so no one group is left completely blind.", Fuel: 10, Hub: 2, Trust: 1},
				{Reply: "Give Maren the convoy’s grid and tell her to wait for extraction, even though you cannot verify who sent the message.", Fuel: 8, Hub: -2, Trust: -3, ConvoyBetray: true},
			}}
	default:
		return NightScript(1)
	}
}
