package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

const (
	GenshinProcess  = "GenshinImpact.exe"
	StarRailProcess = "StarRail.exe"
	ZZZProcess      = "ZenlessZoneZero.exe"

	StartEvent = "started"
	StopEvent  = "stopped"
)

type MonitorEvent struct {
	Name string
	Pid  string
	Type string
}

type Monitor struct {
	ctx    context.Context
	cancel context.CancelFunc
	events chan MonitorEvent
}

func NewMonitor(ctx context.Context) *Monitor {

	ctx, cancel := context.WithCancel(ctx)

	return &Monitor{
		ctx:    ctx,
		cancel: cancel,
		events: make(chan MonitorEvent),
	}
}

func (m *Monitor) Stop() {
	m.cancel()
}

func (m *Monitor) Run() {
	fmt.Println("starting monitor")
	ctx := m.ctx

	startQuery := `
	SELECT * FROM __InstanceCreationEvent WITHIN 2 
	WHERE TargetInstance ISA 'Win32_Process' AND 
	(TargetInstance.Name = 'notepad.exe' OR TargetInstance.Name = 'GenshinImpact.exe' OR TargetInstance.Name = 'StarRail.exe')`

	stopQuery := strings.Replace(startQuery, "__InstanceCreationEvent", "__InstanceDeletionEvent", 1)

	go watchEvent(ctx, startQuery, "started", m.events)
	go watchEvent(ctx, stopQuery, "stopped", m.events)

	<-ctx.Done()
	fmt.Println("Shutting down monitor")
}

func watchEvent(ctx context.Context, query string, eventType string, events chan<- MonitorEvent) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		log.Printf("[%s] CreateObject error: %v", eventType, err)
		return
	}
	defer unknown.Release()

	locator, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		log.Printf("[%s] QueryInterface error: %v", eventType, err)
		return
	}
	defer locator.Release()

	serviceRaw, err := oleutil.CallMethod(locator, "ConnectServer", nil, "root\\CIMV2")
	if err != nil {
		log.Printf("[%s] ConnectServer error: %v", eventType, err)
		return
	}
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	eventSourceRaw, err := oleutil.CallMethod(service, "ExecNotificationQuery", query)
	if err != nil {
		log.Printf("[%s] ExecNotificationQuery error: %v", eventType, err)
		return
	}
	eventSource := eventSourceRaw.ToIDispatch()
	defer eventSource.Release()

	const timeout = int32(1000)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			eventRaw, err := oleutil.CallMethod(eventSource, "NextEvent", timeout)
			if err != nil {
				continue
			}
			event := eventRaw.ToIDispatch()

			instanceVal, _ := oleutil.GetProperty(event, "TargetInstance")
			instance := instanceVal.ToIDispatch()
			nameVal, _ := oleutil.GetProperty(instance, "Name")
			pidVal, _ := oleutil.GetProperty(instance, "ProcessId")

			events <- MonitorEvent{
				Name: nameVal.ToString(),
				Pid:  pidVal.ToString(),
				Type: eventType,
			}

			instance.Release()
			event.Release()
		}
	}
}
