use walkdir::WalkDir;
use std::fs::metadata;
mod utils;

/**
 * Process the file at the provided path
 */
fn process_path(current_path: String) {
    let metadata = metadata(&current_path).unwrap();
    if metadata.is_file() {
        println!("File: {}", current_path);
    } else if metadata.is_dir() {
        println!("Directory: {}", current_path);
    }
}

fn walk_file_path(
    src_build_path: String,
    _directory_build_path: String,
    _directive_type: String,
    ignored_paths: Vec<String>,
) {
    for entry in WalkDir::new(src_build_path) {
        let current_path = entry.unwrap().path().display().to_string();
        // TODO: Dont follow not needed paths
        if !utils::path_is_ignored(&current_path, &ignored_paths) {
            process_path(current_path);
        }
    }
}

fn main() {
    let src_build_path = String::from("/home/sanner/Coding/RAN/ran-app-native/");   
    let directory_build_path = String::from("build-target"); 
    let directive_type = String::from("mobile");
    let ignored_paths = vec![String::from("node_modules"), String::from("build-target")];
    walk_file_path(src_build_path, directory_build_path, directive_type, ignored_paths);
}
