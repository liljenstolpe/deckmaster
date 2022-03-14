package main

import (
	"fmt"
	"image"
	"os"
)

// SinkMuteWidget is a widget displaying if the default Sink is muted.
type SinkMuteWidget struct {
	*ButtonWidget
	pulse      *PulseAudioClient
	mute       bool
	iconUnmute image.Image
	iconMute   image.Image
}

// NewSinkMuteWidget returns a new SinkMuteWidget.
func NewSinkMuteWidget(bw *BaseWidget, opts WidgetConfig) (*SinkMuteWidget, error) {
	widget, err := NewButtonWidget(bw, opts)
	if err != nil {
		return nil, err
	}

	var iconUnmutePath, iconMutePath string
	_ = ConfigValue(opts.Config["icon"], &iconUnmutePath)
	_ = ConfigValue(opts.Config["iconMute"], &iconMutePath)
	iconUnmute, err := preloadImage(widget.base, iconUnmutePath)
	if err != nil {
		return nil, err
	}
	iconMute, err := preloadImage(widget.base, iconMutePath)
	if err != nil {
		return nil, err
	}

	pulse, err := NewPulseAudioClient()
	if err != nil {
		return nil, err
	}

	Sink, err := pulse.DefaultSink()
	if err != nil {
		return nil, err
	}

	return &SinkMuteWidget{
		ButtonWidget: widget,
		pulse:        pulse,
		mute:         Sink.Mute,
		iconUnmute:   iconUnmute,
		iconMute:     iconMute,
	}, nil
}

// RequiresUpdate returns true when the widget wants to be repainted.
func (w *SinkMuteWidget) RequiresUpdate() bool {
	Sink, err := w.pulse.DefaultSink()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't set pulseaudio default Sink mute: %s\n", err)
		return false
	}

	return w.mute != Sink.Mute || w.BaseWidget.RequiresUpdate()
}

// Update renders the widget.
func (w *SinkMuteWidget) Update() error {
	Sink, err := w.pulse.DefaultSink()
	if err != nil {
		return err
	}

	if w.mute != Sink.Mute {
		w.mute = Sink.Mute

		if w.mute {
			w.SetImage(w.iconMute)
		} else {
			w.SetImage(w.iconUnmute)
		}
	}

	return w.ButtonWidget.Update()
}

// TriggerAction gets called when a button is pressed.
func (w *SinkMuteWidget) TriggerAction(hold bool) {
	Sink, err := w.pulse.DefaultSink()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't get pulseaudio default Sink: %s\n", err)
		return
	}
	err = w.pulse.SetSinkMute(Sink, !Sink.Mute)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't set pulseaudio default Sink mute: %s\n", err)
		return
	}
}

// Close gets called when a button is unloaded.
func (w *SinkMuteWidget) Close() error {
	return w.pulse.Close()
}
