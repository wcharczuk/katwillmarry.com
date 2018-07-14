package logger

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestLoggerSubContext(t *testing.T) {
	assert := assert.New(t)

	l := New().WithHeading("test-logger")
	sc := l.SubContext("sub-context")
	assert.NotNil(sc.Logger())
	assert.Equal([]string{"sub-context"}, sc.Headings())
}

func TestSubContextOutput(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	l := New().WithFlags(AllFlags()).WithWriter(NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	sc := l.SubContext("sub-context")
	sc.SyncInfof("this is only a test")
	assert.Equal("[sub-context] [info] this is only a test\n", output.String())
}

func TestSubContextOutputWithLoggerHeading(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	l := New().WithFlags(AllFlags()).WithHeading("test-logger").WithWriter(NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	sc := l.SubContext("sub-context")
	sc.SyncInfof("this is only a test")
	assert.Equal("[test-logger > sub-context] [info] this is only a test\n", output.String())
}

func TestSubContextInfof(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	l := New().WithFlags(AllFlags()).WithWriter(NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	defer l.Close()
	sc := l.SubContext("sub-context")
	sc.Infof("this is only a test")
	l.Drain()
	assert.Equal("[sub-context] [info] this is only a test\n", output.String())
}

func TestSubContextSyncInfof(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	l := New().WithFlags(AllFlags()).WithWriter(NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	defer l.Close()
	sc := l.SubContext("sub-context")
	sc.SyncInfof("this is only a test")
	assert.Equal("[sub-context] [info] this is only a test\n", output.String())
}

func TestSubContextSillyf(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	l := New().WithFlags(AllFlags()).WithWriter(NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	defer l.Close()
	sc := l.SubContext("sub-context")
	sc.Sillyf("this is only a test")
	l.Drain()
	assert.Equal("[sub-context] [silly] this is only a test\n", output.String())
}

func TestSubContextSyncSillyf(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	l := New().WithFlags(AllFlags()).WithWriter(NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	defer l.Close()
	sc := l.SubContext("sub-context")
	sc.SyncSillyf("this is only a test")
	assert.Equal("[sub-context] [silly] this is only a test\n", output.String())
}

func TestSubContextDebugf(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	l := New().WithFlags(AllFlags()).WithWriter(NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	defer l.Close()
	sc := l.SubContext("sub-context")
	sc.Debugf("this is only a test")
	l.Drain()
	assert.Equal("[sub-context] [debug] this is only a test\n", output.String())
}

func TestSubContextSyncDebugf(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	l := New().WithFlags(AllFlags()).WithWriter(NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	defer l.Close()
	sc := l.SubContext("sub-context")
	sc.SyncDebugf("this is only a test")
	assert.Equal("[sub-context] [debug] this is only a test\n", output.String())
}

func TestSubContextWarningf(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	l := New().WithFlags(AllFlags()).WithWriter(NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	defer l.Close()
	sc := l.SubContext("sub-context")
	sc.Warningf("this is only a test")
	l.Drain()
	assert.Equal("[sub-context] [warning] this is only a test\n", output.String())
}

func TestSubContextWarning(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	l := New().WithFlags(AllFlags()).WithWriter(NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	defer l.Close()
	sc := l.SubContext("sub-context")
	sc.Warning(fmt.Errorf("this is only a test"))
	l.Drain()
	assert.Equal("[sub-context] [warning] this is only a test\n", output.String())
}

func TestSubContextSyncWarningf(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	l := New().WithFlags(AllFlags()).WithWriter(NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	defer l.Close()
	sc := l.SubContext("sub-context")
	sc.SyncWarningf("this is only a test")
	assert.Equal("[sub-context] [warning] this is only a test\n", output.String())
}

func TestSubContextSyncWarning(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	l := New().WithFlags(AllFlags()).WithWriter(NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	defer l.Close()
	sc := l.SubContext("sub-context")
	sc.SyncWarning(fmt.Errorf("this is only a test"))
	assert.Equal("[sub-context] [warning] this is only a test\n", output.String())
}

func TestSubContextErrorf(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	l := New().WithFlags(AllFlags()).WithWriter(NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	defer l.Close()
	sc := l.SubContext("sub-context")
	sc.Errorf("this is only a test")
	l.Drain()
	assert.Equal("[sub-context] [error] this is only a test\n", output.String())
}

func TestSubContextError(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	l := New().WithFlags(AllFlags()).WithWriter(NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	defer l.Close()
	sc := l.SubContext("sub-context")
	sc.Error(fmt.Errorf("this is only a test"))
	l.Drain()
	assert.Equal("[sub-context] [error] this is only a test\n", output.String())
}

func TestSubContextSyncErrorf(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	l := New().WithFlags(AllFlags()).WithWriter(NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	defer l.Close()
	sc := l.SubContext("sub-context")
	sc.SyncErrorf("this is only a test")
	assert.Equal("[sub-context] [error] this is only a test\n", output.String())
}

func TestSubContextSyncError(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	l := New().WithFlags(AllFlags()).WithWriter(NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	defer l.Close()
	sc := l.SubContext("sub-context")
	sc.SyncError(fmt.Errorf("this is only a test"))
	assert.Equal("[sub-context] [error] this is only a test\n", output.String())
}

func TestSubContextFatalf(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	l := New().WithFlags(AllFlags()).WithWriter(NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	defer l.Close()
	sc := l.SubContext("sub-context")
	sc.Fatalf("this is only a test")
	l.Drain()
	assert.Equal("[sub-context] [fatal] this is only a test\n", output.String())
}

func TestSubContextFatal(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	l := New().WithFlags(AllFlags()).WithWriter(NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	defer l.Close()
	sc := l.SubContext("sub-context")
	sc.Fatal(fmt.Errorf("this is only a test"))
	l.Drain()
	assert.Equal("[sub-context] [fatal] this is only a test\n", output.String())
}

func TestSubContextSyncFatalf(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	l := New().WithFlags(AllFlags()).WithWriter(NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	defer l.Close()
	sc := l.SubContext("sub-context")
	sc.SyncFatalf("this is only a test")
	assert.Equal("[sub-context] [fatal] this is only a test\n", output.String())
}

func TestSubContextSyncFatal(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	l := New().WithFlags(AllFlags()).WithWriter(NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	defer l.Close()
	sc := l.SubContext("sub-context")
	sc.SyncFatal(fmt.Errorf("this is only a test"))
	assert.Equal("[sub-context] [fatal] this is only a test\n", output.String())
}

func TestSubContextTrigger(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	l := New().WithFlags(AllFlags()).WithWriter(NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	defer l.Close()
	sc := l.SubContext("sub-context")
	sc.Trigger(Messagef(Info, "this is only a test").WithHeadings("message"))
	l.Drain()
	assert.Equal("[sub-context > message] [info] this is only a test\n", output.String())
}

func TestSubContextSyncTrigger(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	l := New().WithFlags(AllFlags()).WithWriter(NewTextWriter(output).WithShowTimestamp(false).WithUseColor(false))
	defer l.Close()
	sc := l.SubContext("sub-context")
	sc.SyncTrigger(Messagef(Info, "this is only a test").WithHeadings("message"))
	assert.Equal("[sub-context > message] [info] this is only a test\n", output.String())
}
