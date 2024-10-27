use walkdir::WalkDir;
use std::fs::metadata;
mod utils;

/**
 * Process the file at the provided path
 */
fn process_path(current_path: &str, source_path: &str, output_path: &str) {
    let metadata = metadata(&current_path).unwrap();
    let new_path = current_path.replace(source_path, output_path);
    if metadata.is_file() {
        println!("File: {} New Path: {}", current_path, new_path);

    } else if metadata.is_dir() {
        println!("Directory: {} New Directory: {}", current_path, new_path);
    }
}

/**
 * Walk the file path and process each file
 */
fn walk_file_path(
    src_path: String,
    output_path: String,
    _directive_type: String,
    ignored_paths: Vec<String>,
) {
    for entry in WalkDir::new(src_path.clone()) {
        let current_path = entry.unwrap().path().display().to_string();
        // TODO: Dont follow not needed paths
        if !utils::path_is_ignored(&current_path, &ignored_paths) {
            process_path(&current_path, &src_path, &output_path);
        }
    }
}

fn main() {
    let src_build_path = String::from("/home/sanner/Coding/RAN/ran-app-native/");   
    let directory_build_path = String::from("/home/sanner/Coding/RAN/ran-app-native/build-target/"); 
    let directive_types = vec![String::from("web"), String::from("mobile")];
    let ignored_paths = vec![String::from("node_modules"), String::from("build-target")];
    for device_type in directive_types {
        walk_file_path( src_build_path.clone(), directory_build_path.clone(), device_type, ignored_paths.clone());
    }
}
