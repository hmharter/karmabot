package karmabot

// MaxKarma the maximum karma that can be given before tripping buzzkill mode
const (
	Development                  = "development"
	Production                   = "production"
	MaxKarma                     = 5
	Positive                     = 1
	Negative                     = -1
	Disconnect                   = "disconnect"
	EventsAPI                    = "events_api"
	Message                      = "message"
	SlashCommands                = "slash_commands"
	Name                         = "name"
	SlackConnURL                 = "https://slack.com/api/apps.connections.open"
	SlackPostMessageURL          = "https://slack.com/api/chat.postMessage"
	SlackPostEphemeralMessageURL = "https://slack.com/api/chat.postEphemeral"
)

var funPositiveMessage = [24]string{
	"is the bees knees",
	"makes the sun shine on a cloudy day",
	"cloned 10 times is my squad goal",
	"puts the 'i' in impressive",
	"is as awesome as a baby goat hopping sideways",
	"started at the bottom now they here",
	"makes the wowest of moments",
	"is the opposite of an asshole",
	"is the champion, my friends",
	"is cooler than the other side of the pillow",
	"is cooler than a cucumber",
	"is the cat's pajamas",
	"returns their grocery cart",
	"says nice things to strangers on the internet",
	"is the fox's socks",
	"is the jewel in the crown",
	"is all that and a bag of chips",
	"is second to none",
	"is the best thing since sliced bread",
	"is the epitome of awesomeness",
	"earned a new badge",
	"is awarded 5 points to Gryffindor",
}
var funNegativeMessage = [13]string{
	"for shame",
	"better luck next time",
	"hope it was worth it",
	"but we still like them",
	"no goats for you",
	"5 points to Slytherin",
	"bless your heart",
	"but after all, tomorrow is another day",
	"you make Cersei Lannister proud",
	"404 - Karma not found",
	"it's a good time for self-reflection",
}
