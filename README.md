NOTE: probably don't use this, use `modernc.org/tk9.0` and GetOpenFile; it's multiplatform (Windows, BSD, Linux; ARM and x86), cgo-less; and will allow you to do more GUI. The negative of tk9.0 that it's a big bazooka, that it includes a big pre-built binary DLL so it's a bit harder to debug. I will leave this open as it can still be useful to someone.

---

Cross-platform quick-and-dirty file-picker that allows more files. Just for Windows and MacOs, I don't need it for Linux.

It's a fork of this repo https://github.com/sqweek/dialog for macos and fork of this repo https://github.com/harry1453/go-common-file-dialog for Windows

Copyright (c) 2018, the dialog authors.
Copyright (c) 2019 Harry Phillips
Copyright (c) 2022 Karel Bilek

See example/simple

Originally I wanted to move all functionality to https://github.com/sqweek/dialog through PRs, but it was stuck on something. You are welcome to add PR for linux here.
