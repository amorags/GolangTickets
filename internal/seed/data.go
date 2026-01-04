package seed

// cities contains a list of US cities for event locations
var cities = []string{
	"New York", "Los Angeles", "Chicago", "Houston", "Phoenix",
	"Philadelphia", "San Antonio", "San Diego", "Dallas", "San Jose",
	"Austin", "Jacksonville", "Fort Worth", "Columbus", "Charlotte",
	"San Francisco", "Indianapolis", "Seattle", "Denver", "Boston",
	"Portland", "Nashville", "Detroit", "Memphis", "Louisville",
	"Baltimore", "Milwaukee", "Albuquerque", "Tucson", "Fresno",
}

// eventTypeImages maps event types to curated Unsplash image URLs
var eventTypeImages = map[string][]string{
	"concert": {
		"https://images.unsplash.com/photo-1470229722913-7c0e2dbbafd3?w=800&h=600&fit=crop",
		"https://images.unsplash.com/photo-1501612780327-45045538702b?w=800&h=600&fit=crop",
		"https://images.unsplash.com/photo-1540039155733-5bb30b53aa14?w=800&h=600&fit=crop",
		"https://images.unsplash.com/photo-1514320291840-2e0a9bf2a9ae?w=800&h=600&fit=crop",
		"https://images.unsplash.com/photo-1506157786151-b8491531f063?w=800&h=600&fit=crop",
	},
	"standup": {
		"https://images.unsplash.com/photo-1585699324551-f6c309eedeca?w=800&h=600&fit=crop",
		"https://images.unsplash.com/photo-1516450360452-9312f5e86fc7?w=800&h=600&fit=crop",
		"https://images.unsplash.com/photo-1527224857830-43a7acc85260?w=800&h=600&fit=crop",
	},
	"tour": {
		"https://images.unsplash.com/photo-1476514525535-07fb3b4ae5f1?w=800&h=600&fit=crop",
		"https://images.unsplash.com/photo-1488646953014-85cb44e25828?w=800&h=600&fit=crop",
		"https://images.unsplash.com/photo-1503220317375-aaad61436b1b?w=800&h=600&fit=crop",
	},
	"lecture": {
		"https://images.unsplash.com/photo-1524178232363-1fb2b075b655?w=800&h=600&fit=crop",
		"https://images.unsplash.com/photo-1523580494863-6f3031224c94?w=800&h=600&fit=crop",
		"https://images.unsplash.com/photo-1517245386807-bb43f82c33c4?w=800&h=600&fit=crop",
	},
	"musical": {
		"https://images.unsplash.com/photo-1503095396549-807759245b35?w=800&h=600&fit=crop",
		"https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=800&h=600&fit=crop",
		"https://images.unsplash.com/photo-1478737270239-2f02b77fc618?w=800&h=600&fit=crop",
	},
	"other": {
		"https://images.unsplash.com/photo-1492684223066-81342ee5ff30?w=800&h=600&fit=crop",
		"https://images.unsplash.com/photo-1505236858219-8359eb29e329?w=800&h=600&fit=crop",
		"https://images.unsplash.com/photo-1511578314322-379afb476865?w=800&h=600&fit=crop",
	},
}

// venueTypes maps event types to venue name components
var venueTypes = map[string][]string{
	"concert": {
		"Arena", "Amphitheater", "Music Hall", "Stadium", "Concert Hall",
		"Club", "Theatre", "Pavilion", "Garden", "Dome",
	},
	"standup": {
		"Comedy Club", "Theater", "Improv", "Lounge", "Playhouse",
		"Comedy House", "Laugh Factory", "Stage", "Comedy Bar",
	},
	"tour": {
		"Visitor Center", "Historic Site", "Museum", "Park", "Gardens",
		"Harbor", "Plaza", "Square", "Pier", "Waterfront",
	},
	"lecture": {
		"Convention Center", "University Hall", "Library", "Auditorium",
		"Conference Center", "Civic Center", "Academic Center", "Hall",
	},
	"musical": {
		"Theater", "Playhouse", "Opera House", "Performing Arts Center",
		"Broadway Theater", "Civic Theatre", "Music Theatre", "Repertory",
	},
	"other": {
		"Community Center", "Hall", "Venue", "Space", "Center",
		"Arts Center", "Cultural Center", "Event Space", "Pavilion",
	},
}

// priceRanges defines min and max price ranges for each event type
var priceRanges = map[string][2]float64{
	"concert": {25.0, 200.0},
	"tour":    {40.0, 150.0},
	"standup": {15.0, 75.0},
	"lecture": {0.0, 50.0},
	"musical": {30.0, 180.0},
	"other":   {0.0, 100.0},
}

// concertDescriptors for generating concert names
var concertDescriptors = []string{
	"Tour", "Festival", "Show", "Concert", "Night", "Sessions",
	"Live", "Experience", "Fest", "Jam", "Rocks", "Unplugged",
}

// concertAdjectives for concert event names
var concertAdjectives = []string{
	"Summer", "Winter", "Midnight", "Electric", "Acoustic", "Jazz",
	"Rock", "Soul", "Blues", "Folk", "Indie", "Alternative",
	"Metal", "Pop", "Classical", "Symphony", "Urban", "Underground",
}

// standupTitles for standup comedy event names
var standupTitles = []string{
	"Comedy Night: Laugh Out Loud",
	"Stand-Up Spectacular",
	"Friday Funnies",
	"Comedy Club Live",
	"The Laugh Factory",
	"Open Mic Night",
	"Comedy Central Presents",
	"Joke Fest",
	"Giggle Fest",
	"The Comedy Show",
}

// tourTypes for tour event names
var tourTypes = []string{
	"Walking Tour", "Bus Tour", "Harbor Cruise", "City Tour",
	"Historical Tour", "Food Tour", "Art Tour", "Architecture Tour",
	"Bike Tour", "Night Tour", "Sunset Tour", "Wine Tour",
}

// tourAdjectives for tour descriptions
var tourAdjectives = []string{
	"Historic", "Scenic", "Guided", "Private", "Cultural",
	"Sunset", "Morning", "Evening", "Downtown", "Waterfront",
}

// lectureTopics for lecture event names
var lectureTopics = []string{
	"AI & Machine Learning",
	"Climate Change Solutions",
	"Future of Technology",
	"Blockchain & Cryptocurrency",
	"Space Exploration",
	"Sustainable Energy",
	"Data Science",
	"Cybersecurity",
	"Digital Marketing",
	"Entrepreneurship",
	"Mental Health",
	"Nutrition & Wellness",
	"Education Reform",
	"Urban Planning",
	"Global Economics",
}

// lectureFormats for lecture event types
var lectureFormats = []string{
	"Summit", "Forum", "Conference", "Workshop", "Seminar",
	"Symposium", "Lecture Series", "Panel Discussion", "Talk",
	"Masterclass", "Deep Dive", "Roundtable",
}

// musicalNames - popular musical titles
var musicalNames = []string{
	"Les Misérables",
	"The Phantom of the Opera",
	"Hamilton",
	"The Lion King",
	"Wicked",
	"Chicago",
	"Cats",
	"Mamma Mia!",
	"The Book of Mormon",
	"Rent",
	"West Side Story",
	"Grease",
	"Jersey Boys",
	"Matilda the Musical",
	"Aladdin",
	"Frozen",
	"Dear Evan Hansen",
	"Come From Away",
	"Kinky Boots",
	"Avenue Q",
}

// eventTypeDistribution defines how many events of each type to create
var eventTypeDistribution = map[string]int{
	"concert": 60, // 30%
	"standup": 40, // 20%
	"tour":    30, // 15%
	"musical": 25, // 12.5%
	"lecture": 25, // 12.5%
	"other":   20, // 10%
}
