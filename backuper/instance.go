package backuper

import "io"
import "github.com/hashicorp/go-multierror"

type InstanceIdentifer interface {
	Name() string
	ID() string
}

//go:generate counterfeiter -o fakes/fake_instance.go . Instance
type Instance interface {
	InstanceIdentifer
	IsBackupable() (bool, error)
	IsPostBackupUnlockable() (bool, error)
	IsRestorable() (bool, error)
	PreBackupLock() error
	Backup() error
	PostBackupUnlock() error
	Restore() error
	Cleanup() error
	StreamBackupFromRemote(io.Writer) error
	StreamBackupToRemote(io.Reader) error
	BackupSize() (string, error)
	BackupChecksum() (BackupChecksum, error)
}

type instances []Instance

func (is instances) IsEmpty() bool {
	return len(is) == 0
}
func (is instances) AllBackupable() (instances, error) {
	var backupableInstances []Instance

	for _, instance := range is {
		if backupable, err := instance.IsBackupable(); err != nil {
			return backupableInstances, err
		} else if backupable {
			backupableInstances = append(backupableInstances, instance)
		}
	}
	return backupableInstances, nil
}

func (is instances) AllPostBackupUnlockable() (instances, error) {
	var unlockableInstances []Instance
	var findUnlockableErrors error = nil

	for _, instance := range is {
		if unlockable, err := instance.IsPostBackupUnlockable(); err != nil {
			findUnlockableErrors = multierror.Append(err)
		} else if unlockable {
			unlockableInstances = append(unlockableInstances, instance)
		}
	}

	return unlockableInstances, findUnlockableErrors
}

func (is instances) AllRestoreable() (instances, error) {
	var backupableInstances []Instance

	for _, instance := range is {
		if backupable, err := instance.IsRestorable(); err != nil {
			return backupableInstances, err
		} else if backupable {
			backupableInstances = append(backupableInstances, instance)
		}
	}
	return backupableInstances, nil
}

func (is instances) Cleanup() error {
	var cleanupErrors error = nil
	for _, instance := range is {
		if err := instance.Cleanup(); err != nil {
			cleanupErrors = multierror.Append(cleanupErrors, err)
		}
	}
	return cleanupErrors
}

func (is instances) PreBackupLock() error {
	var lockErrors error = nil
	for _, instance := range is {
		if err := instance.PreBackupLock(); err != nil {
			lockErrors = multierror.Append(lockErrors, err)
		}
	}

	return lockErrors
}

func (is instances) Backup() error {
	for _, instance := range is {
		err := instance.Backup()
		if err != nil {
			return err
		}
	}
	return nil
}

func (is instances) PostBackupUnlock() error {
	var unlockErrors error = nil
	for _, instance := range is {
		if err := instance.PostBackupUnlock(); err != nil {
			unlockErrors = multierror.Append(unlockErrors, err)
		}
	}
	return unlockErrors
}

func (is instances) Restore() error {
	for _, instance := range is {
		err := instance.Restore()
		if err != nil {
			return err
		}
	}
	return nil
}
