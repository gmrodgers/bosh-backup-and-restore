package orchestrator

type InstanceIdentifer interface {
	Name() string
	Index() string
	ID() string
}

//go:generate counterfeiter -o fakes/fake_instance.go . Instance
type Instance interface {
	InstanceIdentifer
	IsBackupable() bool
	Backup() error
	Jobs() []Job
}

//go:generate counterfeiter -o fakes/fake_job.go . Job
type Job interface {
	HasBackup() bool
	PreBackupLock() error
	PostBackupUnlock() error
	BackupShouldBeLockedBefore() []JobSpecifier
	Name() string
	Release() string
	InstanceIdentifier() string
}

type JobSpecifier struct {
	Name    string
	Release string
}

type ArtifactIdentifier interface {
	InstanceName() string
	InstanceIndex() string
	InstanceID() string
	Name() string
	HasCustomName() bool
}

type instances []Instance

func (is instances) IsEmpty() bool {
	return len(is) == 0
}

func (is instances) Jobs() []Job {
	var jobs []Job
	for _, instance := range is {
		jobs = append(jobs, instance.Jobs()...)
	}

	return jobs
}

func (is instances) AllBackupable() instances {
	var backupableInstances []Instance

	for _, instance := range is {
		if instance.IsBackupable() {
			backupableInstances = append(backupableInstances, instance)
		}
	}
	return backupableInstances
}


func (is instances) AllBackupableOrRestorable() instances {
	var instances []Instance

	for _, instance := range is {
		if instance.IsBackupable() {
			instances = append(instances, instance)
		}
	}
	return instances
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

