package resource

import (
	"time"

	"git.containerum.net/ch/json-types/misc"
)

type Kind string // constants KindNamespace, KindVolume, ... It`s recommended to use strings.ToLower before comparsion

const (
	KindNamespace  Kind = "namespace"
	KindVolume          = "volume"
	KindExtService      = "extservice"
	KindIntService      = "intservice"
	KindDomain          = "domain"
)

type PermissionStatus string // constants PermissionStatusOwner, PermissionStatusRead

const (
	PermissionStatusOwner      PermissionStatus = "owner"
	PermissionStatusRead                        = "read"
	PermissionStatusWrite                       = "write"
	PermissionStatusReadDelete                  = "readdelete"
	PermissionStatusNone                        = "none"
)

type Resource struct {
	ID         string        `json:"id,omitempty" db:"id"`
	CreateTime time.Time     `json:"create_time,omitempty" db:"create_time"`
	Deleted    bool          `json:"deleted,omitempty" db:"deleted"` // not optional because we actually don`t need it if it`s false
	DeleteTime misc.NullTime `json:"delete_time,omitempty" db:"delete_time"`
	TariffID   string        `json:"tariff_id,omitempty" db:"tariff_id"`
}

func (r *Resource) Mask() {
	r.ID = ""
	r.CreateTime = time.Time{}
	r.Deleted = false
	r.DeleteTime.Valid = false
}

type Namespace struct {
	Resource

	RAM                 int `json:"ram" db:"ram"` // megabytes
	CPU                 int `json:"cpu" db:"cpu"`
	MaxExternalServices int `json:"max_external_services" db:"max_ext_services"`
	MaxIntServices      int `json:"max_internal_services" db:"max_int_services"`
	MaxTraffic          int `json:"max_traffic" db:"max_traffic"` // megabytes per month
}

type Volume struct {
	Resource

	Active     misc.NullBool `json:"active,omitempty" db:"active"`
	Capacity   int           `json:"capacity" db:"capacity"` // gigabytes
	Replicas   int           `json:"replicas,omitempty" db:"replicas"`
	Persistent bool          `json:"is_persistent" db:"is_persistent"`
}

func (v *Volume) Mask() {
	v.Resource.Mask()
	v.Active.Valid = false
	v.Replicas = 0
}

type Deployment struct {
	ID          string        `json:"id,omitempty" db:"id"`
	NamespaceID string        `json:"namespace_id,omitempty" db:"ns_id"`
	Name        string        `json:"name" db:"name"`
	RAM         int           `json:"ram" db:"ram"`
	CPU         int           `json:"cpu" db:"cpu"`
	CreateTime  time.Time     `json:"create_time,omitempty" db:"create_time"`
	Deleted     bool          `json:"deleted,omitempty" db:"deleted"`
	DeleteTime  misc.NullTime `json:"delete_time,omitempty" db:"delete_time"`
	Replicas    int           `json:"replicas" db:"replicas"`
}

func (d *Deployment) Mask() {
	d.ID = ""
	d.NamespaceID = ""
	d.CreateTime = time.Time{}
	d.Deleted = false
	d.DeleteTime.Valid = false
}

type PermissionRecord struct {
	PermID                string           `json:"perm_id,omitempty" db:"perm_id"`
	Kind                  Kind             `json:"kind,omitempty" db:"kind"`
	ResourceID            misc.NullString  `json:"resource_id,omitempty" db:"resource_id"` // it can be null for resources without tables
	ResourceLabel         string           `json:"label,omitempty" db:"resource_label"`
	OwnerUserID           string           `json:"owner_user_id,omitempty" db:"owner_user_id"`
	CreateTime            time.Time        `json:"create_time,omitempty" db:"create_time"`
	UserID                string           `json:"user_id" db:"user_id"`
	AccessLevel           PermissionStatus `json:"access" db:"access_level"`
	Limited               bool             `json:"limited,omitempty" db:"limited"`
	AccessLevelChangeTime time.Time        `json:"access_level_change_time" db:"access_level_change_time"`
	NewAccessLevel        PermissionStatus `json:"new_access_level,omitempty" db:"new_access_level"`
}

func (p *PermissionRecord) Mask() {
	p.PermID = ""
	p.Kind = "" // will be already known though
	p.ResourceID.Valid = false
	p.OwnerUserID = ""
	p.CreateTime = time.Time{}
	p.UserID = ""
	p.AccessLevel = p.NewAccessLevel
	p.Limited = false
	p.AccessLevelChangeTime = time.Time{}
	p.NewAccessLevel = ""
}

type Container struct {
	ID       string `json:"id,omitempty" db:"id"`
	DeployID string `json:"depl_id,omitempty" db:"depl_id"`
	Name     string `json:"name" db:"name"`
	Image    string `json:"image" db:"image"`
	RAM      int    `json:"ram" db:"ram"`
	CPU      int    `json:"cpu" db:"cpu"`
}

func (c *Container) Mask() {
	c.ID = ""
	c.DeployID = ""
}

type EnvironmentVariable struct {
	EnvID       string `json:"id,omitempty" db:"env_id"`
	ContainerID string `json:"container_id,omitempty" db:"container_id"`
	Name        string `json:"name" db:"name"`
	Value       string `json:"value" db:"value"`
}

func (e *EnvironmentVariable) Mask() {
	e.EnvID = ""
	e.ContainerID = ""
}

type VolumeMount struct {
	MountID     string          `json:"id,omitempty" db:"mount_id"`
	ContainerID string          `json:"container_id,omitempty" db:"container_id"`
	VolumeID    string          `json:"volume_id,omitempty" db:"volume_id"`
	MountPath   string          `json:"mount_path" db:"mount_path"`
	SubPath     misc.NullString `json:"sub_path,omitempty" db:"sub_path"`
}

func (vm *VolumeMount) Mask() {
	vm.MountID = ""
	vm.ContainerID = ""
	vm.VolumeID = ""
}

// Types below is not for storing in db

type NamespaceWithPermission struct {
	Namespace
	PermissionRecord
}

func (np *NamespaceWithPermission) Mask() {
	np.Namespace.Mask()
	np.PermissionRecord.Mask()
}

type VolumeWithPermission struct {
	Volume
	PermissionRecord
}

func (vp *VolumeWithPermission) Mask() {
	vp.Volume.Mask()
	vp.PermissionRecord.Mask()
}

type NamespaceWithVolumes struct {
	NamespaceWithPermission
	Volume []VolumeWithPermission `json:"volumes"`
}

func (nv *NamespaceWithVolumes) Mask() {
	nv.NamespaceWithPermission.Mask()
	for i := range nv.Volume {
		nv.Volume[i].Mask()
	}
}

type NamespaceWithUserPermissions struct {
	NamespaceWithPermission
	Users []PermissionRecord `json:"users,omitempty"`
}

func (nu *NamespaceWithUserPermissions) Mask() {
	borrowed := nu.UserID != nu.OwnerUserID
	nu.NamespaceWithPermission.Mask()
	if borrowed {
		nu.Users = nil
	} else {
		for i := range nu.Users {
			nu.Users[i].Mask()
		}
	}
}

type VolumeWithUserPermissions struct {
	VolumeWithPermission
	Users []PermissionRecord `json:"users,omitempty"`
}

func (vp *VolumeWithUserPermissions) Mask() {
	borrowed := vp.UserID != vp.OwnerUserID
	vp.VolumeWithPermission.Mask()
	if borrowed {
		vp.Users = nil
	} else {
		for i := range vp.Users {
			vp.Users[i].Mask()
		}
	}
}
