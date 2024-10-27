pub fn path_is_ignored(path: &str, ignored_paths: &Vec<String>) -> bool {
    for ignored_path in ignored_paths {
        if path.contains(ignored_path) {
            return true;
        }
    }
    return false;
}