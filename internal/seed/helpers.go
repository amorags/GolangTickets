package seed

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

// randomCity returns a random city from the cities list
func randomCity() string {
	return cities[rand.Intn(len(cities))]
}

// randomEventType returns a random event type based on distribution
func randomEventType() string {
	types := []string{}
	for eventType, count := range eventTypeDistribution {
		for i := 0; i < count; i++ {
			types = append(types, eventType)
		}
	}
	return types[rand.Intn(len(types))]
}

// randomImageForType returns a random image URL for a given event type
func randomImageForType(eventType string) string {
	images, ok := eventTypeImages[eventType]
	if !ok || len(images) == 0 {
		// Fallback to "other" category
		images = eventTypeImages["other"]
	}
	return images[rand.Intn(len(images))]
}

// randomDateInRange generates a random date between start and end
func randomDateInRange(start, end time.Time) time.Time {
	delta := end.Unix() - start.Unix()
	sec := rand.Int63n(delta)
	return time.Unix(start.Unix()+sec, 0)
}

// generateEventDate returns a date based on distribution (past, present, future)
func generateEventDate() time.Time {
	now := time.Now()

	roll := rand.Intn(100)

	if roll < 15 {
		// 15% past events (last 3 months)
		start := now.AddDate(0, -3, 0)
		return randomDateInRange(start, now)
	} else if roll < 65 {
		// 50% near-term events (next 2 weeks to 2 months)
		start := now.AddDate(0, 0, 14)
		end := now.AddDate(0, 2, 0)
		return randomDateInRange(start, end)
	} else {
		// 35% future events (2-12 months out)
		start := now.AddDate(0, 2, 0)
		end := now.AddDate(0, 12, 0)
		return randomDateInRange(start, end)
	}
}

// randomPrice generates a price based on event type
// 20% chance of free events, rest within type's price range
func randomPrice(eventType string) float64 {
	// 20% chance of free event
	if rand.Intn(100) < 20 {
		return 0.0
	}

	priceRange, ok := priceRanges[eventType]
	if !ok {
		priceRange = priceRanges["other"]
	}

	min := priceRange[0]
	max := priceRange[1]

	// Generate random price in range
	price := min + rand.Float64()*(max-min)

	// Round to nearest $5
	price = math.Round(price/5) * 5

	return price
}

// randomCapacity generates a capacity based on distribution
func randomCapacity() int {
	roll := rand.Intn(100)

	if roll < 40 {
		// 40% small venues: 50-200
		return 50 + rand.Intn(151)
	} else if roll < 80 {
		// 40% medium venues: 200-1000
		return 200 + rand.Intn(801)
	} else if roll < 95 {
		// 15% large venues: 1000-5000
		return 1000 + rand.Intn(4001)
	} else {
		// 5% mega venues: 5000-20000
		return 5000 + rand.Intn(15001)
	}
}

// randomVenueName generates a venue name based on event type
func randomVenueName(eventType string) string {
	venues, ok := venueTypes[eventType]
	if !ok || len(venues) == 0 {
		venues = venueTypes["other"]
	}

	// Pick random venue type
	venueType := venues[rand.Intn(len(venues))]

	// Sometimes add a prefix (city name, person name, or adjective)
	roll := rand.Intn(3)

	switch roll {
	case 0:
		// City prefix: "Boston Arena"
		return fmt.Sprintf("%s %s", randomCity(), venueType)
	case 1:
		// Person's name: "Madison Square Garden" style
		return fmt.Sprintf("%s %s", gofakeit.LastName(), venueType)
	default:
		// Descriptive: "Grand Theater"
		adjectives := []string{
			"Grand", "Royal", "Historic", "Modern", "Classic",
			"Premier", "Central", "Downtown", "Uptown", "Riverside",
			"Metropolitan", "City", "National", "State", "Imperial",
		}
		return fmt.Sprintf("%s %s", adjectives[rand.Intn(len(adjectives))], venueType)
	}
}

// generateConcertName creates a realistic concert event name
func generateConcertName() string {
	roll := rand.Intn(3)

	switch roll {
	case 0:
		// "Artist Name + Tour/Concert"
		artist := gofakeit.Name()
		descriptor := concertDescriptors[rand.Intn(len(concertDescriptors))]
		return fmt.Sprintf("%s %s", artist, descriptor)
	case 1:
		// "Adjective + Descriptor"
		adj := concertAdjectives[rand.Intn(len(concertAdjectives))]
		descriptor := concertDescriptors[rand.Intn(len(concertDescriptors))]
		return fmt.Sprintf("%s %s", adj, descriptor)
	default:
		// "The Adjective Descriptor"
		adj := concertAdjectives[rand.Intn(len(concertAdjectives))]
		descriptor := concertDescriptors[rand.Intn(len(concertDescriptors))]
		return fmt.Sprintf("The %s %s", adj, descriptor)
	}
}

// generateStandupName creates a standup comedy event name
func generateStandupName() string {
	roll := rand.Intn(2)

	if roll == 0 {
		// Use predefined title
		return standupTitles[rand.Intn(len(standupTitles))]
	} else {
		// "Comedian Name Live"
		return fmt.Sprintf("%s Live", gofakeit.Name())
	}
}

// generateTourName creates a tour event name
func generateTourName() string {
	adj := tourAdjectives[rand.Intn(len(tourAdjectives))]
	tourType := tourTypes[rand.Intn(len(tourTypes))]

	roll := rand.Intn(2)
	if roll == 0 {
		return fmt.Sprintf("%s %s", adj, tourType)
	} else {
		city := randomCity()
		return fmt.Sprintf("%s %s %s", adj, city, tourType)
	}
}

// generateLectureName creates a lecture event name
func generateLectureName() string {
	topic := lectureTopics[rand.Intn(len(lectureTopics))]
	format := lectureFormats[rand.Intn(len(lectureFormats))]

	return fmt.Sprintf("%s %s", topic, format)
}

// generateMusicalName returns a musical name
func generateMusicalName() string {
	return musicalNames[rand.Intn(len(musicalNames))]
}

// generateOtherEventName creates a generic event name
func generateOtherEventName() string {
	return fmt.Sprintf("%s %s Event", gofakeit.Adjective(), gofakeit.Noun())
}

// generateEventName creates an event name based on type
func generateEventName(eventType string) string {
	switch eventType {
	case "concert":
		return generateConcertName()
	case "standup":
		return generateStandupName()
	case "tour":
		return generateTourName()
	case "lecture":
		return generateLectureName()
	case "musical":
		return generateMusicalName()
	default:
		return generateOtherEventName()
	}
}

// generateEventDescription creates a description based on event type
func generateEventDescription(eventType, eventName string) string {
	templates := map[string][]string{
		"concert": {
			fmt.Sprintf("Join us for an unforgettable night of live music at %s. Experience amazing performances and incredible energy!", eventName),
			fmt.Sprintf("%s brings you the best in live entertainment. Don't miss this spectacular show!", eventName),
			fmt.Sprintf("Get ready for an electrifying performance at %s. Limited tickets available!", eventName),
		},
		"standup": {
			fmt.Sprintf("Laugh until it hurts at %s! Featuring top comedians and rising stars.", eventName),
			fmt.Sprintf("Join us for a night of comedy at %s. Guaranteed laughs and good times!", eventName),
			fmt.Sprintf("%s presents an evening of hilarious stand-up comedy. Bring your friends!", eventName),
		},
		"tour": {
			fmt.Sprintf("Explore and discover on this %s. Perfect for visitors and locals alike!", eventName),
			fmt.Sprintf("Join our %s and see the sights from a unique perspective. Book your spot today!", eventName),
			fmt.Sprintf("Experience the best of the city on this %s. Expert guides included!", eventName),
		},
		"lecture": {
			fmt.Sprintf("Expand your knowledge at %s. Featuring industry experts and thought leaders.", eventName),
			fmt.Sprintf("Join us for %s with insights from leading professionals in the field.", eventName),
			fmt.Sprintf("%s brings together experts to discuss cutting-edge topics. Don't miss it!", eventName),
		},
		"musical": {
			fmt.Sprintf("Experience the magic of %s live on stage. A theatrical masterpiece!", eventName),
			fmt.Sprintf("Don't miss this spectacular production of %s. Award-winning performances!", eventName),
			fmt.Sprintf("See %s in all its glory. Book your tickets for this limited run!", eventName),
		},
		"other": {
			fmt.Sprintf("Join us for %s, a unique experience you won't want to miss!", eventName),
			fmt.Sprintf("%s is coming to town! Get your tickets before they sell out.", eventName),
			fmt.Sprintf("Be part of %s and create unforgettable memories.", eventName),
		},
	}

	options, ok := templates[eventType]
	if !ok {
		options = templates["other"]
	}

	return options[rand.Intn(len(options))]
}
