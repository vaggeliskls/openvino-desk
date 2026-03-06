# OpenVINO Desktop

A Windows desktop application for managing and serving OpenVINO models via [OpenVINO Model Server (OVMS)](https://github.com/openvinotoolkit/model_server).

## Features

- Download (pull) or export models from Hugging Face
- Manage installed models (view, delete)
- Start/stop the OVMS inference server
- Detect available OpenVINO devices (CPU, GPU, NPU)
- Auto-start with Windows
- Local REST API for external integrations

## Requirements

- Windows 10/11 x64
- Internet connection (for first-time setup and model downloads)

## Getting Started

1. Launch the application. On first run it downloads OVMS and sets up the Python export environment automatically.
2. Go to the **Models** tab, search for a model on Hugging Face, select a target device and click **Pull** or **Export**.
3. The OVMS server starts automatically once models are available.

## Settings

| Setting | Description |
|---------|-------------|
| Setup Folder | Directory where OVMS and models are installed (default: `~/openvino-desktop`) |
| OVMS Download URL | URL to the OVMS Windows zip archive |
| uv Download URL | URL to `uv.exe` used for the export environment |
| REST API Port | Port for the local REST API (default: `3333`) |
| Search Limit | Max number of models returned by Hugging Face search |

> Changing the REST API port requires a restart to take effect.

## REST API

The app exposes a local HTTP API on port `3333` (configurable in Settings).

### Endpoints

#### `GET /status`
Returns server state and available devices.
```json
{
  "running": true,
  "deps_ready": true,
  "ovms_ready": true,
  "version": "2026.0",
  "available_devices": ["CPU"]
}
```

#### `GET /models`
Returns the list of installed models.
```json
[
  {
    "name": "OpenVINO/Qwen3-Embedding-0.6B-int8-ov",
    "base_path": "C:\\Users\\user\\openvino-desktop\\models\\OpenVINO\\Qwen3-Embedding-0.6B-int8-ov",
    "target_device": "CPU"
  }
]
```

#### `POST /models/pull`
Pulls a pre-converted OpenVINO model from Hugging Face. Returns `202 Accepted` immediately — operation runs in the background.

```bash
curl -X POST http://localhost:3333/models/pull \
  -H "Content-Type: application/json" \
  -d '{"model_id": "OpenVINO/Qwen3-Embedding-0.6B-int8-ov", "target_device": "CPU", "pipeline_tag": "feature-extraction"}'
```

| Field | Values |
|-------|--------|
| `pipeline_tag` | `"text-generation"` or `"feature-extraction"` |
| `target_device` | `"CPU"`, `"GPU"`, `"NPU"`, `"AUTO"` |

#### `POST /models/export`
Exports and converts a model from Hugging Face using `export_model.py`. Returns `202 Accepted` immediately.

```bash
curl -X POST http://localhost:3333/models/export \
  -H "Content-Type: application/json" \
  -d '{"model_id": "Qwen/Qwen3-Embedding-0.6B", "target_device": "CPU", "task": "feature-extraction", "extra_opts": {"weight-format": "int8"}}'
```

| Field | Values |
|-------|--------|
| `task` | `"text-generation"` or `"feature-extraction"` |
| `target_device` | `"CPU"`, `"GPU"`, `"NPU"`, `"AUTO"` |
| `extra_opts` | Optional key/value pairs passed as CLI flags to `export_model.py` |

#### `GET /job`
Poll the status of the current or last pull/export job.
```json
{ "busy": false, "last_ok": true, "last_error": "" }
```

Only one pull/export can run at a time. A second request while busy returns `409 Conflict`.

## Development

```bash
# Install dependencies
make deps

# Run in dev mode (hot-reload for frontend)
make ui-dev

# Build production binary
make ui-build
```
