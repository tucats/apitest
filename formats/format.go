package formats

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Duration formats a duration into a human-readable format. This is similar to the
// Go time.Duration.String() function, but with additional formatting controls
// for small time intervals.
func Duration(d time.Duration, extendedFormat bool) string {
	if !extendedFormat {
		return d.String()
	}

	if d == 0 {
		return "0s"
	}

	// If this is a very small duration that is less than a second, use the
	// default formatter.
	if v := math.Abs(float64(d)); v < float64(time.Second) {
		if v > float64(time.Second) {
			return fmt.Sprintf("%6.2fs", v/float64(time.Second))
		}

		if v > float64(time.Millisecond) {
			return fmt.Sprintf("%6.2fms", v/float64(time.Millisecond))
		}

		if v > float64(time.Microsecond) {
			return fmt.Sprintf("%6.2fµs", v/(float64(time.Microsecond)))
		}

		return d.String()
	}

	// Otherwise, lets build a formatted duration string.
	var result strings.Builder

	// if the duration is negative, add a sign and make the remaining duration
	// a positive value.
	if d < 0 {
		result.WriteRune('-')

		d = -d
	}

	// If the number of hours is greater than a day, extract the number of days
	// as a separate part of the duration string.
	if hours := int(d.Hours()); hours > 0 {
		if hours > 23 {
			result.WriteString(fmt.Sprintf("%dd", hours/24))
			hours = hours % 24
		}

		// If the remaining number of hours is greater than 0, add the hours
		if hours > 0 {
			if result.Len() > 1 {
				result.WriteRune(' ')
			}

			result.WriteString(fmt.Sprintf("%dh", hours))
		}
	}

	// If there are more than 0 minutes, add the minutes.
	minutes := int64(d.Minutes()) % 60
	if minutes >= 1 {
		if result.Len() > 1 {
			result.WriteRune(' ')
		}

		result.WriteString(fmt.Sprintf("%dm", minutes))
	}

	// If there are more than 1 second, add the seconds.
	seconds := int64(d.Seconds()) % 60
	if seconds >= 1 {
		if result.Len() > 1 {
			result.WriteRune(' ')
		}

		result.WriteString(fmt.Sprintf("%ds", seconds))
	}

	// Adjust the duration string to be right-justified if
	// it is six or fewer characters into a six-character
	// field to match the duration size for small intervals.

	text := result.String()
	for len(text) < 6 {
		text = " " + text
	}

	return text
}

// Gibberish generates a random string of gibberish characters from the given UUID.
// The value is created by converting the 128-bit UUID value into two 64-bit integers.
// Each is then formatted using base-32, with a character set of 32 numbers and letters
// comprising the arabic digits and letters, but omitting the values "0", "1", "o", and "l"
// for increased human readability.
func Gibberish(u uuid.UUID) string {
	var result strings.Builder

	digits := []byte("abcdefghjkmnpqrstuvwxyz23456789")
	radix := uint64(len(digits))

	// Make two 64-bit integers from the UUID value
	var hi, low uint64

	for i := 0; i < 8; i++ {
		hi = hi<<8 + uint64(u[i])
	}

	for i := 9; i < 16; i++ {
		low = low<<8 + uint64(u[i])
	}

	for low > 0 {
		result.WriteByte(digits[low%radix])
		low = low / radix
	}

	for hi > 0 {
		result.WriteByte(digits[hi%radix])
		hi = hi / radix
	}

	text := result.String()
	if len(text) == 0 {
		return "-empty-"
	}

	return text
}
