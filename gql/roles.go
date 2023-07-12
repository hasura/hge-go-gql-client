package gql

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hasura/hge-go-gql-client/util"
)

const (
	RoleAdmin  = "admin"
	RoleUser   = "user"
	RolePublic = "public"
)

// Actor denotes who is making a request to our APIs
type Actor struct {
	Role     string
	UserID   *uuid.UUID
	Email    string
	UsesSAML bool
}

func NewAdminActor() *Actor {
	return &Actor{
		Role: RoleAdmin,
	}
}

func NewUserActor(userID *uuid.UUID, email string) *Actor {
	return &Actor{
		Role:   RoleUser,
		UserID: userID,
		Email:  email,
	}
}

func (a *Actor) HasRole(role string) bool {
	return strings.EqualFold(a.Role, role)
}

func (a *Actor) IsAdmin() bool {
	return a.HasRole(RoleAdmin)
}

func (a *Actor) AsAdmin() *Actor {
	return &Actor{
		Role:     RoleAdmin,
		UserID:   a.UserID,
		Email:    a.Email,
		UsesSAML: a.UsesSAML,
	}
}

// Access holds information about what a given actor can access.
// code using Access values will determine whether a user can access a given resource based on the actor's role, identity or allowed lists
type Access struct {
	*Actor
	AllowedProjectIDs        []uuid.UUID
	MetricsAllowedProjectIDs []uuid.UUID
	AdminProjectIDs          []uuid.UUID
}

func NewAccess(actor *Actor) *Access {
	return &Access{Actor: actor}
}

// NewAccessFromSessionVariables parses de headers in the given StringMap
// And builds an Access object based on them
func NewAccessFromSessionVariables(sessionVariables util.StringMap) *Access {
	// get allowed project IDs
	getArrayUUID := func(name string) []uuid.UUID {
		rawValues, err := util.PostgresArrayToStrings(sessionVariables.Get(name))
		if err != nil {
			return []uuid.UUID{}
		}

		results, err := util.MapStringToUUID(rawValues)
		if err != nil {
			return []uuid.UUID{}
		}

		return results
	}

	usesSAML, _ := strconv.ParseBool(sessionVariables.Get(util.XHasuraIsSAMLUser))

	access := Access{
		Actor: &Actor{
			Role:     sessionVariables.Get(util.XHasuraRole),
			Email:    sessionVariables.Get(util.XHasuraUserEmail),
			UsesSAML: usesSAML,
		},
		AllowedProjectIDs:        getArrayUUID(util.XHasuraAllowedProjectIDs),
		MetricsAllowedProjectIDs: getArrayUUID(util.XHasuraAllowedMetricsProjectIDs),
		AdminProjectIDs:          getArrayUUID(util.XHasuraAdminProjectIDs),
	}

	userID, err := uuid.Parse(sessionVariables.Get(util.XHasuraUserID))
	if err != nil {
		access.UserID = nil
	} else {
		access.UserID = &userID
	}

	return &access
}
