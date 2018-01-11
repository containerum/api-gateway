package datastore

import "errors"

var (
	ErrUnableFindListener     = errors.New("Unable to find Listener")
	ErrUnableGetListeners     = errors.New("Unable to get Listeners")
	ErrListenerIDIsEmpty      = errors.New("Listener id is empty")
	ErrUnableToUpdateListener = errors.New("Unable to update Listener")
	ErrUnableToCreateListener = errors.New("Unable to create Listener")
	ErrUnableToDeleteListener = errors.New("Unable to delete Listener")
)
