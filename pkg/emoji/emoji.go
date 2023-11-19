package emoji

var (
	Love      Emoji = "❤️"
	Like      Emoji = "👍"
	CatLove   Emoji = "😻"
	LoveIt    Emoji = "😍"
	InLove    Emoji = "🥰"
	Heart     Emoji = Love
	Cheers    Emoji = "🥂"
	Hot       Emoji = "🔥"
	Party     Emoji = "🎉"
	Birthday  Emoji = "🎂️"
	Sparkles  Emoji = "✨"
	Rainbow   Emoji = "🌈"
	Pride     Emoji = "🏳️‍🌈"
	SeeNoEvil Emoji = "🙈"
	Unknown   Emoji = ""
)

// Reactions specifies reaction emojis by name.
var Reactions = map[string]Emoji{
	"love":        Love,
	"+1":          Like,
	"cat-love":    CatLove,
	"love-it":     LoveIt,
	"in-love":     InLove,
	"heart":       Heart,
	"cheers":      Cheers,
	"hot":         Hot,
	"party":       Party,
	"birthday":    Birthday,
	"sparkles":    Sparkles,
	"rainbow":     Rainbow,
	"pride":       Pride,
	"see-no-evil": SeeNoEvil,
}

// Names specifies the reaction names by emoji.
var Names = map[Emoji]string{
	Love:      "love",
	Like:      "+1",
	CatLove:   "cat-love",
	LoveIt:    "love-it",
	InLove:    "in-love",
	Heart:     "heart",
	Cheers:    "cheers",
	Hot:       "hot",
	Party:     "party",
	Birthday:  "birthday",
	Sparkles:  "sparkles",
	Rainbow:   "rainbow",
	Pride:     "pride",
	SeeNoEvil: "see-no-evil",
}

type Emoji string

// Unknown checks if the reaction is unknown.
func (emo Emoji) Unknown() bool {
	if l := len(emo); l < 2 || len(emo) > 64 {
		return true
	}

	return Names[emo] == ""
}

// Name returns the ASCII name of the reaction.
func (emo Emoji) Name() string {
	return Names[emo]
}

// String returns the reaction as string.
func (emo Emoji) String() string {
	return string(emo)
}

// Bytes returns the reaction emoji as a slice with a maximum size of 64 bytes.
func (emo Emoji) Bytes() (b []byte) {
	if b = []byte(emo); len(b) <= 64 {
		return b
	}

	return b[0:64]
}

// Find finds a reaction by name and emoji.
func Find(reaction string) Emoji {
	if reaction == "" {
		return Unknown
	}

	if found := Reactions[reaction]; found != "" {
		return found
	} else if found = Emoji(reaction); found.Unknown() {
		return Unknown
	} else {
		return found
	}
}

// Known checks if the emoji represents a known reaction.
func Known(reaction string) bool {
	return !Find(reaction).Unknown()
}
