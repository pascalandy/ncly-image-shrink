---
name: image-shrink
description: Use to make one local JPEG, PNG or WebP smaller for sharing, without its location, author, camera or capture date.
---

# Shrink an image for sharing

Use `ncly image shrink` for one still image. Dry-run the call first, then run it when the predicted destination is suitable. Read the flags and result fields with `ncly describe image shrink --json`.

The copy stays beside its input. Success means fewer bytes. The input stays unchanged. Orientation is applied, color data is retained, and descriptive metadata is removed.

On `IMAGE_SHRINK_NOT_SMALLER`, the input may already be optimized. Choose a smaller dimension or a lower lossy quality if acceptable. Never claim that a dry run proves byte saving.

On `IMAGE_SHRINK_EXISTS`, inspect the copy before choosing force. On `PREREQ_MISSING`, relay the install hint to the human. For animation, multiple pictures or formats outside JPEG, PNG and WebP, use another tool.
