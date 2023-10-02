package jokeResources

import (
	"math/rand"
	"strings"
)

var nouns = []string{
	"lion", "elephant", "giraffe", "monkey", "tiger", "zebra",
	"penguin", "koala", "kangaroo", "parrot", "cheetah", "hippo",
	"crocodile", "rhino", "jaguar", "koala", "panda", "gorilla",
	"lemur", "lemming", "lynx", "meerkat", "moose", "otter",
	"pangolin", "panther", "quokka", "quoll", "rabbit", "raccoon",
	"salamander", "seagull", "sheep", "sloth", "snake", "spider",
	"starfish", "toucan", "turtle", "vulture", "wallaby", "walrus",
	"whale", "wombat", "yak", "armadillo", "axolotl", "baboon",
	"badger", "bison", "bluejay", "bobcat", "buffalo", "butterfly",
	"camel", "capybara", "chameleon", "chinchilla", "chipmunk", "cobra",
	"cockatoo", "coyote", "dingo", "dolphin", "donkey", "dragonfly",
	"duck", "eagle", "falcon", "ferret", "flamingo", "fox",
	"gecko", "gazelle", "gnu", "goldfish", "gopher", "hedgehog",
	"hornet", "hyena", "ibex", "iguana", "jackal", "jellyfish",
	"kookaburra", "lemur", "lizard", "llama", "lobster", "lynx",
	"macaw", "magpie", "manatee", "mongoose", "narwhal", "newt",
	"ocelot", "octopus", "opossum", "orangutan", "ostrich", "owl",
	"panda", "pangolin", "panther", "peacock", "pelican", "penguin",
}

var verbs = []string{
	"walks", "runs", "jogs", "sprints", "strolls",
	"hops", "skips", "jumps", "leaps", "bounds",
	"gallops", "travels", "wanders", "ambles", "saunters",
	"strides", "tiptoes", "marches", "dances", "prances",
	"glides", "slides", "skates", "coasts", "drifts",
	"flows", "sways", "swings", "flits", "flutters",
	"soars", "flies", "hovers", "ascends", "descends",
	"swoops", "dives", "plunges", "plods", "shuffles",
	"treads", "moseys", "crawls", "slithers", "creeps",
	"scuttles", "scampers", "scrambles", "climbs", "ascends",
	"descends", "descends", "vaults", "hurdles", "jumps",
	"bounces", "bounds", "traverses", "crosses", "navigates",
	"prowls", "stalks", "tracks", "chases", "pursues",
	"retreats", "withdraws", "approaches", "advances", "proceeds",
	"retreats", "returns", "regresses", "proceeds", "progresses",
	"snakes", "slaloms", "zigzags", "weaves", "threads",
	"drives", "cruises", "rides", "sails", "glides",
	"surfs", "tumbles", "rolls", "spins", "twirls",
	"rotates", "turns", "pivots", "swivels", "twists",
	"gyrates", "oscillates", "vibrates", "quivers", "trembles",
}

var adjectives = []string{
	"swift", "majestic", "playful", "graceful", "curious",
	"fierce", "gentle", "cunning", "elegant", "vibrant",
	"sleek", "intelligent", "sociable", "resourceful", "adaptable",
	"mysterious", "energetic", "agile", "resilient", "sensitive",
	"charming", "loyal", "furry", "colorful", "alert",
	"regal", "dexterous", "sly", "wise", "endearing",
	"vigilant", "dainty", "bold", "shy", "inquisitive",
	"tenacious", "mellow", "ambitious", "robust", "dynamic",
	"gracious", "spirited", "cooperative", "gregarious", "innocent",
	"quirky", "patient", "majestic", "adventurous", "sturdy",
	"perceptive", "sensitive", "eager", "elusive", "bright",
	"jovial", "mysterious", "resourceful", "graceful", "vivid",
	"zesty", "daring", "radiant", "clever", "buoyant",
	"magnificent", "joyful", "ambitious", "grateful", "tenacious",
	"capable", "lively", "persistent", "effervescent", "exuberant",
	"observant", "unassuming", "harmonious", "meticulous", "colorful",
	"dynamic", "determined", "vibrant", "attentive", "fascinating",
	"intrepid", "sincere", "versatile", "intuitive", "adaptable",
	"invincible", "courageous", "resourceful", "resilient", "exotic",
}

type Joke struct {
	IndefiniteArticle string `json:"IndefiniteArticle"`
	Noun              string `json:"Noun"`
	Verb              string `json:"Verb"`
	Adjective         string `json:"Adjective"`
}

// Returns a randomly generated object of type joke
func GetRandomJoke() Joke {
	noun := nouns[rand.Intn(len(nouns))]
	verb := verbs[rand.Intn(len(verbs))]
	adjective := adjectives[rand.Intn(len(adjectives))]
	indefiniteArticle := "A"
	if strings.Contains("aeiou", strings.ToLower(string(noun[0]))) {
		indefiniteArticle = "An"
	}
	return Joke{
		IndefiniteArticle: indefiniteArticle,
		Noun:              noun,
		Verb:              verb,
		Adjective:         adjective,
	}
}
