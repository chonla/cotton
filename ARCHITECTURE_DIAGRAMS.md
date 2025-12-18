# Architecture Diagram

## Dependency Graph

```
┌─────────────────────────────────────────────────────┐
│                    User Executable                  │
│                    (./cotton)                        │
└────────────────────┬────────────────────────────────┘
                     │
         ┌───────────▼──────────┐
         │   main.go            │
         │  (Entry Point)       │
         │  - Creates App       │
         │  - Calls Run()       │
         │  - Exits with code   │
         └───────────┬──────────┘
                     │
         ┌───────────▼────────────────────┐
         │   cli Package (CLI Layer)      │
         │                                │
         │  ┌──────────────────────────┐ │
         │  │ app.go (App)             │ │
         │  │ ├─ NewApp()              │ │
         │  │ ├─ Run(args)             │ │
         │  │ │  ├─ Parse flags        │ │
         │  │ │  │   - debug/quiet/monochrome/fail-fast/version
         │  │ │  │   - --timeout <ms> (validated, capped ≤ 60000)
         │  │ │  ├─ Resolve path       │ │
         │  │ │  └─ Call runner        │ │
         │  │ └─ Exit code handling    │ │
         │  └──────────────────────────┘ │
         │                                │
         │  ┌──────────────────────────┐ │
         │  │ formatter.go             │ │
         │  │ (ResultFormatter)        │ │
         │  │ └─ PrintResult()         │ │
         │  └──────────────────────────┘ │
         │                                │
         └────────────┬───────────────────┘
                      │
                      │ imports
                      │
         ┌────────────▼──────────────────────────────┐
         │   engine Package (Core Engine)            │
         │                                           │
         │  ┌──────────────────────────────────────┐ │
         │  │ runner.go (TestRunner)               │ │
         │  │ ├─ RunSpec()                         │ │
         │  │ │  ├─ Setup phase                    │ │
         │  │ │  ├─ Variable resolution            │ │
         │  │ │  ├─ Request execution (http timeout)│ │
         │  │ │  ├─ Expectation evaluation         │ │
         │  │ │  ├─ Teardown phase                 │ │
         │  │ │  └─ Cycle detection (visited)      │ │
         │  │ └─ Manages TestContext               │ │
         │  └──────────────────────────────────────┘ │
         │                                           │
         │  ┌──────────────────────────────────────┐ │
         │  │ parser.go (MarkdownParser)           │ │
         │  │ ├─ Parse()                           │ │
         │  │ │  ├─ Extract title                  │ │
         │  │ │  ├─ Extract setup links            │ │
         │  │ │  ├─ Extract variables (precompiled regex)
         │  │ │  ├─ Extract request body           │ │
         │  │ │  ├─ Extract expectations           │ │
         │  │ │  └─ Extract teardown links         │ │
         │  │ └─ Returns ParsedSpec                │ │
         │  └──────────────────────────────────────┘ │
         │                                           │
         │  ┌──────────────────────────────────────┐ │
         │  │ executor.go (HTTPExecutor)           │ │
         │  │ ├─ ExecuteRequest()                  │ │
         │  │ ├─ substituteVariables() (strict: missing → ParseError)
         │  │ ├─ parseRequest()                    │ │
         │  │ └─ Manages last response body        │ │
         │  └──────────────────────────────────────┘ │
         │                                           │
         │  ┌──────────────────────────────────────┐ │
         │  │ evaluator.go (Evaluator)             │ │
         │  │ ├─ EvaluateExpectations()            │ │
         │  │ ├─ evaluateExpectation()             │ │
         │  │ ├─ compare()                         │ │
         │  │ └─ Comparison operators              │ │
         │  └──────────────────────────────────────┘ │
         │                                           │
         │  ┌──────────────────────────────────────┐ │
         │  │ errors.go (Typed Errors)             │ │
         │  │ ├─ ParseError                        │ │
         │  │ └─ RunError                          │ │
         │  └──────────────────────────────────────┘ │
         │                                           │
         └───────────────────────────────────────────┘
```

## Data Flow

```
┌──────────────────┐
│  User runs cmd   │
│  cotton <spec|dir|glob> --timeout 30000 │
└────────┬─────────┘
         │
         ▼
┌─────────────────────┐
│  main.go            │ ◄── Creates App instance
└────────┬────────────┘
         │
         ▼
┌─────────────────────┐
│  cli.App.Run()      │ ◄── Parses flags & paths, validates timeout
└────────┬────────────┘
         │
         ▼
┌──────────────────────────────┐
│  engine.TestRunner.RunSpec() │ ◄── Starts test execution (fail-fast aware)
└────────┬─────────────────────┘
         │
         ├─► engine.MarkdownParser   ◄── Reads & parses spec
         │
         ├─► engine.HTTPExecutor     ◄── Executes requests (http.Client.Timeout)
         │   │
         │   └─► Variable substitution
         │
         ├─► engine.Evaluator        ◄── Evaluates expectations
         │
         └─► engine.TestResult       ◄── Collects results; errors typed
             │
             ▼
         cli.ResultFormatter  ◄── Formats output
             │
             ▼
         Console output
             │
             ▼
         Exit code (non-zero if any failure)
```

## Module Isolation

```
                    ┌─────────────────┐
                    │   main.go       │
                    └────────┬────────┘
                             │
                    ┌────────▼────────┐
                    │  cli package    │ ◄── User Interface Layer
                    │  - app.go       │
                    │  - formatter.go │
                    └────────┬────────┘
                             │
                    ┌────────▼────────────────┐
                    │  engine package        │ ◄── Business Logic Layer
                    │  - runner.go           │    (No CLI dependencies)
                    │  - parser.go           │
                    │  - executor.go         │
                    │  - evaluator.go        │
                    │  - errors.go           │
                    └────────────────────────┘


No coupling ──┐
between      │  engine can be imported by:
layers ◄─────┤  - Web UI
             │  - REST API
             │  - IDE plugin
             │  - Other tools
             └──►
```

## Extensibility Points

```
engine Package (Stable, Core Logic)
│
├─ Reuse in Web UI
│  └─ New package: webui/
│     ├─ Import engine
│     └─ Add HTTP handlers
│
├─ Reuse in REST API
│  └─ New package: api/
│     ├─ Import engine
│     └─ Add REST endpoints
│
├─ Alternative CLI
│  └─ Extend cli/
│     ├─ Add JSON formatter
│     ├─ Add XML formatter
│     └─ Add verbose logging
│
└─ Programmatic API
   └─ External tools
      ├─ Import "github.com/chonla/cotton/engine"
      └─ Use TestRunner directly
```
