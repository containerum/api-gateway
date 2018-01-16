package datastore

var (
	//SQLGetListener get listener by ID
	SQLGetListener = "SELECT id, created_at, updated_at, name, o_auth," +
		"active, strip_path, listen_path, upstream_url, method, group_refer FROM listeners " +
		" WHERE id = $1;"
	//SQLGetListeners get all listeners
	SQLGetListeners = "SELECT id, created_at, updated_at, name, o_auth," +
		"active, strip_path, listen_path, upstream_url, method, group_refer FROM listeners;"
	//SQLGetListenersActive get all listeners by active
	SQLGetListenersActive = "SELECT id, created_at, updated_at, name, o_auth," +
		"active, strip_path, listen_path, upstream_url, method, group_refer FROM listeners " +
		" WHERE active = $1;"
	//SQLCreateListener create listener
	SQLCreateListener = "INSERT INTO listeners (name, o_auth, active, strip_path, listen_path, upstream_url, method, group_refer) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at;"
	//SQLUpdateListener update listener
	SQLUpdateListener = "UPDATE listeners SET name = $1, o_auth = $2, active = $3, " +
		"strip_path = $4, listen_path = $5, upstream_url = $6, method = $7, group_refer = $8, updated_at = now() WHERE id = $9 RETURNING updated_at;"
	//SQLDeleteListener delete listener
	SQLDeleteListener = "DELETE FROM listeners WHERE id = $1;"

	//SQLGetGroups get all groups
	SQLGetGroups = "SELECT id, created_at, updated_at, name, active FROM groups;"
	//SQLCreateGroup create group
	SQLCreateGroup = "INSERT INTO groups (name, active) VALUES ($1, $2) RETURNING id, created_at;"
)
