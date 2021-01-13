package tasks

// TaskStatus represents the status of a single task
type TaskStatus struct {
	Err error
}

func (st *TaskStatus) String() string {
	switch st {
	case DoneOk:
		return "OK"
	case Cancelled:
		return "Cancelled"
	case Running:
		return "Running"
	case Scheduled:
		return "Scheduled"
	default:
		if st.IsFailure() {
			return "Err: " + st.Err.Error()
		}
		return "UNKNOWN"
	}
}

// Scheduled is the status of a task scheduled but not yet executed
var Scheduled = &TaskStatus{}

// DoneOk is the status of a task successfully executed
var DoneOk = &TaskStatus{}

// Failed returns the status of a task that failed with an error
func Failed(err error) *TaskStatus {
	return &TaskStatus{err}
}

// IsFailure returns whether the task status represents a failure
func (status *TaskStatus) IsFailure() bool {
	return status.Err != nil
}

// Running is the status of a running task
var Running = &TaskStatus{}

// Cancelled is the status of a task that won't run, because
// one of it's prerequisites has failed
var Cancelled = &TaskStatus{}

/*
// SimulationTaskStatus is the status of a single WRF run
type SimulationTaskStatus struct {
	FinalDewetraDelivery        *TaskStatus
	VdADelivery                 []*TaskStatus
	ArpalDelivery               []*TaskStatus
	ArpaPiemonteDelivery        []*TaskStatus
	ContinuumDelivery           []*TaskStatus
	ArpaPiemonteIndexesDelivery []*TaskStatus

	AUXDownloadDomain1       []*TaskStatus
	AUXPostProcessDomain1    []*TaskStatus
	AUXDownloadDomain3       []*TaskStatus
	AUXPostProcessDomain3    []*TaskStatus
	OUTPostProcess           []*TaskStatus
	OUTPostProcessedDownload []*TaskStatus
	FinalZTDScript           *TaskStatus
}

// NewSimulationTaskStatus returns a new SimulationTaskStatus instance
// initialized for a run of totHours hours
func NewSimulationTaskStatus(totHours int) SimulationTaskStatus {
	arpaPiemontePhaseCount := int(math.Ceil(float64(totHours-1) / 12))
	status := SimulationTaskStatus{
		VdADelivery:                 make([]*TaskStatus, totHours),
		ArpalDelivery:               make([]*TaskStatus, totHours),
		ArpaPiemonteDelivery:        make([]*TaskStatus, totHours),
		ContinuumDelivery:           make([]*TaskStatus, totHours),
		AUXDownloadDomain1:          make([]*TaskStatus, totHours),
		AUXPostProcessDomain1:       make([]*TaskStatus, totHours),
		AUXDownloadDomain3:          make([]*TaskStatus, totHours),
		AUXPostProcessDomain3:       make([]*TaskStatus, totHours),
		OUTPostProcess:              make([]*TaskStatus, totHours),
		OUTPostProcessedDownload:    make([]*TaskStatus, totHours),
		ArpaPiemonteIndexesDelivery: make([]*TaskStatus, arpaPiemontePhaseCount),
		FinalZTDScript:              Scheduled,
		FinalDewetraDelivery:        Scheduled,
	}

	for i := 0; i < arpaPiemontePhaseCount; i++ {
		status.ArpaPiemonteIndexesDelivery[i] = Scheduled
	}

	for i := 0; i < totHours; i++ {
		status.VdADelivery[i] = Scheduled
		status.ArpalDelivery[i] = Scheduled
		status.ArpaPiemonteDelivery[i] = Scheduled
		status.ContinuumDelivery[i] = Scheduled
		status.AUXDownloadDomain1[i] = Scheduled
		status.AUXPostProcessDomain1[i] = Scheduled
		status.AUXDownloadDomain3[i] = Scheduled
		status.AUXPostProcessDomain3[i] = Scheduled
		status.OUTPostProcess[i] = Scheduled
		status.OUTPostProcessedDownload[i] = Scheduled
	}

	return status
}
*/
