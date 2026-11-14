package bundleinstall

import (
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/kitecloud/kite/kite-service/internal/api/handler"
)

var templateVersionLabel = regexp.MustCompile(`^[1-9][0-9]{0,8}(\.[0-9]{1,10})?$`)

// Labels are decimal strings so the creator's formatting survives publication.
// The integer revision still identifies the release for installs and reviews.
func nextTemplateVersionLabel(revision int, previous, requested string) (string, error) {
	if previous == "" {
		if strings.TrimSpace(requested) == "" {
			return "", nil
		}
		previous = strconv.Itoa(revision)
	}

	label := strings.TrimSpace(requested)
	if label == "" {
		major, _, _ := strings.Cut(previous, ".")
		n, _ := strconv.Atoi(major)
		label = strconv.Itoa(n + 1)
	}

	if !templateVersionLabel.MatchString(label) {
		return "", handler.ErrBadRequest("invalid_version_label", "Enter a version like 1.01, 2.5, or 3.")
	}

	current, _ := new(big.Rat).SetString(previous)
	next, _ := new(big.Rat).SetString(label)
	if next.Cmp(current) <= 0 {
		return "", handler.ErrBadRequest("invalid_version_label", fmt.Sprintf("Enter a version higher than %s.", previous))
	}
	return label, nil
}
