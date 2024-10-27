use walkdir::WalkDir;

fn walkFilePath(
    src_build_path: String,
    _directory_build_path: String,
    _directive_type: String,
    _ignored_paths: Vec<String>,
) {
    let current_directory = WalkDir::new(src_build_path);
    for entry in current_directory {
        println!("{:?}", entry);
    }
}

fn main() {
    let src_build_path = String::from("/home/sanner/Coding/RAN/ran-app-native/");   
    let directory_build_path = String::from("build-target"); 
    let directive_type = String::from("mobile");
    let ignored_paths = vec![String::from("node_modules"), String::from("build-target")];
    walkFilePath(src_build_path, directory_build_path, directive_type, ignored_paths);
    println!("Hello, world!");
}
