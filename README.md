# Process Loader
A concurrent terminal progress display package for Go that shows multiple processes with their progress percentages, names, and current steps - each on dedicated lines.

## Overview
The Process Loader package provides real-time progress visualization for concurrent operations. Each process is displayed on its own terminal line with:

Process Name: The file or task name

Progress Percentage: Visual bar (0-100%)

Current Step: Detailed status of what's happening

## Installation

```shell
go get github.com/mindwingx/go-procel

```
### Core Concepts

#### Process Structure
Each process tracks:

- **Name**: Identifier (e.g., "photo.jpg", "data-export")

- **Percentage**: Completion progress (0-100)

- **State/Step**: Current operation stage (e.g., "uploading", "validating", "processing")

#### Terminal Behavior
- Each process occupies a fixed terminal line

- Progress updates happen in-place without scrolling

- Concurrent processes don't interfere with each other's display

- Finished processes can be cleared or left visible

### Basic Usage

#### Single Process Example

```go
package main

import (
    "time"
    pl "github.com/mindwingx/go-procel"
)

func main() {
    // Initialize a new process
    upload := pl.NewProcess()
    upload.SetName("customer-data.csv")
    
    // Simulate upload progress
    upload.Load("initializing", 0).Process()
    time.Sleep(1 * time.Second)
    
    upload.Load("reading file", 25).Process()
    time.Sleep(1 * time.Second)
    
    upload.Load("uploading to server", 75).Process()
    time.Sleep(2 * time.Second)
    
    upload.Load("completed successfully", 100).Process()
    upload.Finish()
}
```

**Output:**

```text
customer-data.csv [0% >............................ ~ initializing]
customer-data.csv [25% =======>................... ~ reading file]  
customer-data.csv [75% ====================>...... ~ uploading to server]
customer-data.csv [100% ==========================> ~ completed successfully]
```

### Advanced Usage

#### Multiple Concurrent Processes

```go
func processFile(filename string, steps []string) {
    p := pl.NewProcess()
    p.SetName(filename)
    defer p.Finish()
    
    for i, step := range steps {
        percent := (i + 1) * 100 / len(steps)
        p.Load(step, percent).Process()
        time.Sleep(2 * time.Second)
    }
}

// Run multiple file processes concurrently
func main() {
    files := map[string][]string{
        "photo1.png": {"validating format", "compressing image", "uploading to CDN", "updating database"},
        "photo2.jpg": {"validating format", "applying filters", "generating thumbnails", "storage optimization"},
        "data.json":  {"parsing structure", "validation checks", "encrypting data", "backup creation"},
    }
    
    for filename, steps := range files {
        go processFile(filename, steps)
    }
    
    // Wait for all processes to complete
    time.Sleep(10 * time.Second)
}
```

**Concurrent Output:**

```text
photo1.png [25% =======>................... ~ validating format]
photo2.jpg [50% ==============>............ ~ applying filters] 
data.json [75% ====================>...... ~ encrypting data]
```
### Process Lifecycle

```go
// Finish marks process complete (auto-removes from tracking)
process.Finish()

// Cleanup removes the process line from terminal display  
process.Cleanup()
```

### Progress Bar Format
The visual progress display shows:

```text
filename.ext [XX% visual-bar ~ current-step-description]
```

Where:

XX%: Numeric percentage (0-100)

visual-bar: `=======>......` growing with progress

current-step-description: What the process is currently doing

### Best Practices
- Always call Finish() - Prevents memory leaks in long-running applications

- Use descriptive state names - Helps users understand what's happening

- Progress percentages should be meaningful - Reflect actual completion

- Handle errors gracefully - Update state to show failure reasons

- Use defer for cleanup - Ensure processes are properly terminated
