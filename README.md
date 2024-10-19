# gocomploy

Cross-Platform Repo Compiler. Write once, deploy everywhere

On compilation:

- if web, remove tauri dependencies from package.json as well as all references
  to it in imports
- remove all `EXCLUSIVE`s pointing to devices which are not part of the target
- Fill `src/config/platform.ts` with the correct information
