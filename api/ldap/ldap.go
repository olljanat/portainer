package ldap

import (
	"fmt"
	"strings"

	"github.com/portainer/portainer"
	"github.com/portainer/portainer/crypto"

	"gopkg.in/ldap.v2"
)

const (
	// ErrUserNotFound defines an error raised when the user is not found via LDAP search
	// or that too many entries (> 1) are returned.
	ErrUserNotFound = portainer.Error("User not found or too many entries returned")
)

// Service represents a service used to authenticate users against a LDAP/AD.
type Service struct{
	UserService           portainer.UserService
	LDAPService           portainer.LDAPService
	TeamService           portainer.TeamService
	TeamMembershipService portainer.TeamMembershipService
}

func searchUser(username string, conn *ldap.Conn, settings []portainer.LDAPSearchSettings) (string, error) {
	var userDN string
	found := false
	for _, searchSettings := range settings {
		searchRequest := ldap.NewSearchRequest(
			searchSettings.BaseDN,
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			fmt.Sprintf("(&%s(%s=%s))", searchSettings.Filter, searchSettings.UserNameAttribute, username),
			[]string{"dn"},
			nil,
		)

		// Deliberately skip errors on the search request so that we can jump to other search settings
		// if any issue arise with the current one.
		sr, err := conn.Search(searchRequest)
		if err != nil {
			continue
		}

		if len(sr.Entries) == 1 {
			found = true
			userDN = sr.Entries[0].DN
			break
		}
	}

	if !found {
		return "", ErrUserNotFound
	}

	return userDN, nil
}

func createConnection(settings *portainer.LDAPSettings) (*ldap.Conn, error) {

	if settings.TLSConfig.TLS || settings.StartTLS {
		config, err := crypto.CreateTLSConfigurationFromDisk(settings.TLSConfig.TLSCACertPath, settings.TLSConfig.TLSCertPath, settings.TLSConfig.TLSKeyPath, settings.TLSConfig.TLSSkipVerify)
		if err != nil {
			return nil, err
		}
		config.ServerName = strings.Split(settings.URL, ":")[0]

		if settings.TLSConfig.TLS {
			return ldap.DialTLS("tcp", settings.URL, config)
		}

		conn, err := ldap.Dial("tcp", settings.URL)
		if err != nil {
			return nil, err
		}

		err = conn.StartTLS(config)
		if err != nil {
			return nil, err
		}

		return conn, nil
	}

	return ldap.Dial("tcp", settings.URL)
}

// AuthenticateUser is used to authenticate a user against a LDAP/AD.
func (service *Service) AuthenticateUser(username, password string, settings *portainer.LDAPSettings) error {

	connection, err := createConnection(settings)
	if err != nil {
		return err
	}
	defer connection.Close()

	err = connection.Bind(settings.ReaderDN, settings.Password)
	if err != nil {
		return err
	}

	userDN, err := searchUser(username, connection, settings.SearchSettings)
	if err != nil {
		if err == ErrUserNotFound {			
			user := &portainer.User{
				Username: username,
				Role:     portainer.StandardUserRole,
			}

			if err := service.UserService.CreateUser(user); err != nil {
				return err
			}

			if err := service.addLdapUserIntoTeams(user, settings); err != nil {
				return err
			}

		} else {
			return err
		}
	}

	err = connection.Bind(userDN, password)
	if err != nil {
		return err
	}

	return nil
}

/*
func addLdapUser(username string, settings *portainer.LDAPSettings) error {
// user *portainer.User, settings *portainer.LDAPSettings) error {
 // func (*Service) AuthenticateUser(username, password string, settings *portainer.LDAPSettings) error {
 
	user := &portainer.User{
		Username: username,
		Role:     portainer.StandardUserRole,
	}

	if err := handler.UserService.CreateUser(user); err != nil {
		return err
	}

	if err := handler.addLdapUserIntoTeams(user, settings); err != nil {
		return err
	}

	return nil
 
}*/

// GetUserGroups is used to retrieve user groups from LDAP/AD.
func (*Service) GetUserGroups(username string, settings *portainer.LDAPSettings) ([]string, error) {
	connection, err := createConnection(settings)
	if err != nil {
		return nil, err
	}
	defer connection.Close()

	err = connection.Bind(settings.ReaderDN, settings.Password)
	if err != nil {
		return nil, err
	}

	userDN, err := searchUser(username, connection, settings.SearchSettings)
	if err != nil {
		return nil, err
	}

	userGroups := getGroups(userDN, connection, settings.SearchSettings)

	return userGroups, nil
}

// Get a list of group names for specified user from LDAP/AD
func getGroups(userDN string, conn *ldap.Conn, settings []portainer.LDAPSearchSettings) []string {
	groups := make([]string, 0)

	for _, searchSettings := range settings {
		searchRequest := ldap.NewSearchRequest(
			searchSettings.GroupBaseDN,
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			fmt.Sprintf("(&%s(%s=%s))", searchSettings.GroupFilter, searchSettings.GroupAttribute, userDN),
			[]string{"cn"},
			nil,
		)

		// Deliberately skip errors on the search request so that we can jump to other search settings
		// if any issue arise with the current one.
		sr, err := conn.Search(searchRequest)
		if err != nil {
			continue
		}

		// Collect groups
		for _, entry := range sr.Entries {
			for _, attr := range entry.Attributes {
				groups = append(groups, attr.Values[0])
			}
		}
	}

	return groups
}

// TestConnectivity is used to test a connection against the LDAP server using the credentials
// specified in the LDAPSettings.
func (*Service) TestConnectivity(settings *portainer.LDAPSettings) error {

	connection, err := createConnection(settings)
	if err != nil {
		return err
	}
	defer connection.Close()

	err = connection.Bind(settings.ReaderDN, settings.Password)
	if err != nil {
		return err
	}
	return nil
}

func (service *Service) addLdapUserIntoTeams(user *portainer.User, settings *portainer.LDAPSettings) error {
	teams, err := service.TeamService.Teams()
	if err != nil {
		return err
	}

	userLdapGroups, err := service.LDAPService.GetUserGroups(user.Username, settings)
	if err != nil {
		return err
	}

	userMemberships, err := service.TeamMembershipService.TeamMembershipsByUserID(user.ID)
	if err != nil {
		return err
	}

	for _, team := range teams {
		if teamExists(team.Name, userLdapGroups) {

			if teamMembershipExists(team.ID, userMemberships) {
				continue
			}

			membership := &portainer.TeamMembership{
				UserID: user.ID,
				TeamID: team.ID,
				Role:   portainer.TeamMember,
			}

			service.TeamMembershipService.CreateTeamMembership(membership)
		}
	}
	return nil
}

func teamExists(teamName string, ldapGroups []string) bool {
	for _, group := range ldapGroups {
		if group == teamName {
			return true
		}
	}
	return false
}

func teamMembershipExists(teamID portainer.TeamID, memberships []portainer.TeamMembership) bool {
	for _, membership := range memberships {
		if membership.TeamID == teamID {
			return true
		}
	}
	return false
}
