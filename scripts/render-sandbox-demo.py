#!/usr/bin/env python3
"""Encode verified sandbox interactions without any audio or voice-over.

Requires ffmpeg/ffprobe. Never starts or retries checkout runs. Source timings
and real interaction frames are preserved; no captions are added.
"""
import argparse
import json
from pathlib import Path
import subprocess

ROOT = Path(__file__).resolve().parents[1]


def run(args):
    result = subprocess.run(args, text=True, capture_output=True)
    if result.returncode:
        raise RuntimeError(f"{args[0]} failed: {result.stderr[-1600:]}")
    return result.stdout


def probe(path):
    return json.loads(run(["ffprobe", "-v", "error", "-show_streams", "-show_format", "-of", "json", str(path)]))


def validate(source):
    if source.get("mode", "live") not in ("live", "saved_run_review"):
        raise ValueError("Unknown recording mode")
    if source["status"] != "verified" or len(source["runs"]) != 2:
        raise ValueError("Stopped/unverified recordings cannot be exported as a completed demo")
    baseline, retest = source["runs"]
    if baseline["status"] != "findings" or retest["status"] != "passed" or retest["parent_id"] != baseline["id"]:
        raise ValueError("Invalid baseline/retest evidence")
    for item in source["runs"]:
        if item["environment"] != "sandbox" or item["order_count"] != 1 or not item["total_verified"]:
            raise ValueError("Unverified sandbox result")
        if any(step.get("error") for step in item["steps"]):
            raise ValueError("A recorded step contains an error; export refused")
    if not source["frames"] or source["duration_seconds"] >= 180:
        raise ValueError("Empty/oversized recording")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("directory", type=Path)
    args = parser.parse_args()
    directory = args.directory.resolve()
    directory.relative_to(ROOT / "tmp" / "video")
    source = json.loads((directory / "recording.json").read_text())
    validate(source)
    frames = source["frames"]
    duration = source["duration_seconds"]
    lines = ["ffconcat version 1.0"]
    for i, frame in enumerate(frames):
        path = (directory / frame["file"]).resolve()
        path.relative_to(directory / "frames")
        if not path.is_file() or "'" in str(path):
            raise ValueError("Invalid frame path")
        start = frame["seconds"] if i else 0
        end = frames[i + 1]["seconds"] if i + 1 < len(frames) else duration
        if end <= start:
            raise ValueError("Frame timings are not increasing")
        lines += [f"file '{path}'", "option framerate 1000", f"duration {end-start:.9f}"]
    lines += [f"file '{(directory / frames[-1]['file']).resolve()}'", "option framerate 1000"]
    concat = directory / "frames.ffconcat"
    concat.write_text("\n".join(lines) + "\n")
    mode = source.get("mode", "live")
    review = mode == "saved_run_review"
    output = directory / ("CheckoutGuard-sandbox-review-silent.mp4" if review else "CheckoutGuard-sandbox-silent.mp4")
    title = "CheckoutGuard - Saved sandbox results review" if review else "CheckoutGuard - Silent sandbox demo"
    comment = "Review of verified saved runs; no new checkout or inference; no audio; original review interaction timings" if review else "Live app interaction recording; no audio; original timings preserved; sandbox segment"
    run(["ffmpeg", "-hide_banner", "-loglevel", "error", "-y", "-f", "concat", "-safe", "0", "-i", str(concat), "-map", "0:v:0", "-an", "-vf", "fps=30,format=yuv420p", "-c:v", "libx264", "-preset", "medium", "-crf", "18", "-t", f"{duration:.6f}", "-movflags", "+faststart", "-metadata", f"title={title}", "-metadata", f"comment={comment}", str(output)])
    report = probe(output)
    video = next(s for s in report["streams"] if s["codec_type"] == "video")
    actual = float(report["format"]["duration"])
    if any(s["codec_type"] == "audio" for s in report["streams"]):
        raise ValueError("Audio stream found: silent export refused")
    if (video["width"],video["height"]) != (1920,1080) or actual >= 180 or abs(actual-duration) > 1:
        raise ValueError("Export format or duration verification failed")
    run(["ffmpeg", "-hide_banner", "-v", "error", "-i", str(output), "-f", "null", "-"])
    result = {"status":"verified", "mode":mode, "video":str(output.relative_to(ROOT)), "duration_seconds":actual, "resolution":"1920x1080", "fps":30, "audio_streams":0, "voice_over":False, "captions":False, "capture_style":source.get("style"), "frames_captured":len(frames), "timing":"Original review interaction timings preserved" if review else "Original live recording preserved; no waiting time removed", "sandbox_runs":[r["id"] for r in source["runs"]], "step_errors":0, "decode_check":"passed", "scope":"Review of verified saved sandbox results; no new checkout or inference; no Shopee/Tokopedia test; partner evidence segments still needed" if review else "Sandbox segment only; no Shopee/Tokopedia test; partner evidence segments still needed"}
    (directory / "export-validation.json").write_text(json.dumps(result,indent=2)+"\n")
    print(json.dumps(result,indent=2))


if __name__ == "__main__":
    main()
