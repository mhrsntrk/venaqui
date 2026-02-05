package realdebrid

import (
	"fmt"
)

// UnrestrictOptions provides options for link unrestriction
type UnrestrictOptions struct {
	Password string
}

// UnrestrictLinkWithOptions converts a hoster link to an unrestricted download link with options
func (c *Client) UnrestrictLinkWithOptions(link string, options *UnrestrictOptions) (*UnrestrictedLink, error) {
	// For now, we'll use the basic UnrestrictLink method
	// Future enhancement: add support for password-protected links
	if options != nil && options.Password != "" {
		return nil, fmt.Errorf("password-protected links not yet supported")
	}

	return c.UnrestrictLink(link)
}

// IsLinkSupported checks if a link is supported by Real-Debrid
// This is a placeholder - actual implementation would require checking supported hosts
func IsLinkSupported(link string) bool {
	// Basic check - Real-Debrid supports many hosters
	// In a real implementation, you might want to check against a list of supported hosts
	return link != ""
}
