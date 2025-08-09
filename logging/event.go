// Package logging defines the event-driven logging system for GoSpring framework.
// It provides a structured way to log internal events during container operations,
// dependency injection, and lifecycle management.
package logging

import (
	"fmt"
	"time"
)

// Event defines an event emitted by GoSpring framework.
// All events implement this interface to provide a unified logging mechanism.
type Event interface {
	// String returns a human-readable description of the event
	String() string
}

// ContainerCreated is emitted when a new container is created.
type ContainerCreated struct {
	Timestamp time.Time
}

func (e *ContainerCreated) String() string {
	return fmt.Sprintf("[%s] Container created", e.Timestamp.Format("15:04:05.000"))
}

// ComponentRegistered is emitted when a component is registered to the container.
type ComponentRegistered struct {
	Timestamp    time.Time
	ComponentID  string
	ComponentType string
	Scope        string
}

func (e *ComponentRegistered) String() string {
	return fmt.Sprintf("[%s] Component registered: %s (type: %s, scope: %s)", 
		e.Timestamp.Format("15:04:05.000"), e.ComponentID, e.ComponentType, e.Scope)
}

// ComponentScanned is emitted when a component is discovered during scanning.
type ComponentScanned struct {
	Timestamp    time.Time
	PackagePath  string
	ComponentType string
	Tags         map[string]string
}

func (e *ComponentScanned) String() string {
	return fmt.Sprintf("[%s] Component scanned: %s in package %s (tags: %v)", 
		e.Timestamp.Format("15:04:05.000"), e.ComponentType, e.PackagePath, e.Tags)
}

// DependencyInjected is emitted when a dependency is successfully injected.
type DependencyInjected struct {
	Timestamp      time.Time
	TargetType     string
	DependencyType string
	FieldName      string
	ByType         bool
	ByName         bool
}

func (e *DependencyInjected) String() string {
	injectionType := "by type"
	if e.ByName {
		injectionType = "by name"
	}
	return fmt.Sprintf("[%s] Dependency injected: %s.%s <- %s (%s)", 
		e.Timestamp.Format("15:04:05.000"), e.TargetType, e.FieldName, e.DependencyType, injectionType)
}

// DependencyInjectionFailed is emitted when dependency injection fails.
type DependencyInjectionFailed struct {
	Timestamp      time.Time
	TargetType     string
	DependencyType string
	FieldName      string
	Error          error
}

func (e *DependencyInjectionFailed) String() string {
	return fmt.Sprintf("[%s] Dependency injection failed: %s.%s <- %s (error: %v)", 
		e.Timestamp.Format("15:04:05.000"), e.TargetType, e.FieldName, e.DependencyType, e.Error)
}

// ComponentCreated is emitted when a component instance is created.
type ComponentCreated struct {
	Timestamp     time.Time
	ComponentID   string
	ComponentType string
	CreationTime  time.Duration
}

func (e *ComponentCreated) String() string {
	return fmt.Sprintf("[%s] Component created: %s (type: %s, time: %v)", 
		e.Timestamp.Format("15:04:05.000"), e.ComponentID, e.ComponentType, e.CreationTime)
}

// ComponentDestroyed is emitted when a component instance is destroyed.
type ComponentDestroyed struct {
	Timestamp     time.Time
	ComponentID   string
	ComponentType string
}

func (e *ComponentDestroyed) String() string {
	return fmt.Sprintf("[%s] Component destroyed: %s (type: %s)", 
		e.Timestamp.Format("15:04:05.000"), e.ComponentID, e.ComponentType)
}

// LifecycleStarting is emitted before a component's Init method is called.
type LifecycleStarting struct {
	Timestamp     time.Time
	ComponentID   string
	ComponentType string
	MethodName    string
}

func (e *LifecycleStarting) String() string {
	return fmt.Sprintf("[%s] Lifecycle starting: %s.%s (type: %s)", 
		e.Timestamp.Format("15:04:05.000"), e.ComponentID, e.MethodName, e.ComponentType)
}

// LifecycleStarted is emitted after a component's Init method is called.
type LifecycleStarted struct {
	Timestamp     time.Time
	ComponentID   string
	ComponentType string
	MethodName    string
	Duration      time.Duration
	Error         error
}

func (e *LifecycleStarted) String() string {
	if e.Error != nil {
		return fmt.Sprintf("[%s] Lifecycle started with error: %s.%s (type: %s, duration: %v, error: %v)", 
			e.Timestamp.Format("15:04:05.000"), e.ComponentID, e.MethodName, e.ComponentType, e.Duration, e.Error)
	}
	return fmt.Sprintf("[%s] Lifecycle started: %s.%s (type: %s, duration: %v)", 
		e.Timestamp.Format("15:04:05.000"), e.ComponentID, e.MethodName, e.ComponentType, e.Duration)
}

// LifecycleStopping is emitted before a component's Destroy method is called.
type LifecycleStopping struct {
	Timestamp     time.Time
	ComponentID   string
	ComponentType string
	MethodName    string
}

func (e *LifecycleStopping) String() string {
	return fmt.Sprintf("[%s] Lifecycle stopping: %s.%s (type: %s)", 
		e.Timestamp.Format("15:04:05.000"), e.ComponentID, e.MethodName, e.ComponentType)
}

// LifecycleStopped is emitted after a component's Destroy method is called.
type LifecycleStopped struct {
	Timestamp     time.Time
	ComponentID   string
	ComponentType string
	MethodName    string
	Duration      time.Duration
	Error         error
}

func (e *LifecycleStopped) String() string {
	if e.Error != nil {
		return fmt.Sprintf("[%s] Lifecycle stopped with error: %s.%s (type: %s, duration: %v, error: %v)", 
			e.Timestamp.Format("15:04:05.000"), e.ComponentID, e.MethodName, e.ComponentType, e.Duration, e.Error)
	}
	return fmt.Sprintf("[%s] Lifecycle stopped: %s.%s (type: %s, duration: %v)", 
		e.Timestamp.Format("15:04:05.000"), e.ComponentID, e.MethodName, e.ComponentType, e.Duration)
}

// ContextStarting is emitted when application context starts.
type ContextStarting struct {
	Timestamp time.Time
}

func (e *ContextStarting) String() string {
	return fmt.Sprintf("[%s] Application context starting", e.Timestamp.Format("15:04:05.000"))
}

// ContextStarted is emitted when application context has started.
type ContextStarted struct {
	Timestamp    time.Time
	Duration     time.Duration
	ComponentCount int
}

func (e *ContextStarted) String() string {
	return fmt.Sprintf("[%s] Application context started (duration: %v, components: %d)", 
		e.Timestamp.Format("15:04:05.000"), e.Duration, e.ComponentCount)
}

// ContextStopping is emitted when application context is stopping.
type ContextStopping struct {
	Timestamp time.Time
}

func (e *ContextStopping) String() string {
	return fmt.Sprintf("[%s] Application context stopping", e.Timestamp.Format("15:04:05.000"))
}

// ContextStopped is emitted when application context has stopped.
type ContextStopped struct {
	Timestamp time.Time
	Duration  time.Duration
}

func (e *ContextStopped) String() string {
	return fmt.Sprintf("[%s] Application context stopped (duration: %v)", 
		e.Timestamp.Format("15:04:05.000"), e.Duration)
}

// ScanStarting is emitted when component scanning starts.
type ScanStarting struct {
	Timestamp     time.Time
	ComponentType string
	PackagePath   string
}

func (e *ScanStarting) String() string {
	return fmt.Sprintf("[%s] Component scan starting: %s in package %s", 
		e.Timestamp.Format("15:04:05.000"), e.ComponentType, e.PackagePath)
}

// ScanCompleted is emitted when component scanning completes.
type ScanCompleted struct {
	Timestamp     time.Time
	ComponentType string
	PackagePath   string
	ComponentName string
	Scope         string
	Duration      time.Duration
	Success       bool
	Error         error
}

func (e *ScanCompleted) String() string {
	if !e.Success {
		return fmt.Sprintf("[%s] Component scan failed: %s in package %s (duration: %v, error: %v)", 
			e.Timestamp.Format("15:04:05.000"), e.ComponentType, e.PackagePath, e.Duration, e.Error)
	}
	return fmt.Sprintf("[%s] Component scan completed: %s (%s) in package %s (scope: %s, duration: %v)", 
		e.Timestamp.Format("15:04:05.000"), e.ComponentName, e.ComponentType, e.PackagePath, e.Scope, e.Duration)
}