# 📦 Zip Repackager

A simple Go program to repackage a zip archive by flattening its file structure, removing symlinks/directories, and keeping only the largest files when name clashes occur.

---

## 🧠 Problem Statement

Given an input `.zip` file, repackage it into a new `.zip` file such that:

- All nested files (e.g., `foo/bar/baz.jpg`) are flattened to just `baz.jpg`.
- Directories and symbolic links are ignored.
- If multiple files have the same name, only the **largest one** is kept.
- File contents must **remain unchanged**.
- File data should be stored with the `STORE` method (no compression).
- SHA-256 hashes before and after repackaging must match.
- The program should gracefully handle errors with a single-line message and exit code `1`.

---

## 🚀 Usage
// run commend //
go run main.go <input.zip> <output.zip>
