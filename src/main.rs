use clap::{arg, command};
use directories::ProjectDirs;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let dirs = ProjectDirs::from("com", "dislogical", "home-builder").unwrap();

    let matches = command!()
        .arg(arg!(--one <VALUE>).required(true))
        .get_matches();

    println!("config: {}", dirs.config_dir().display());
    println!("one: {}", matches.get_one::<String>("one").unwrap());

    Ok(())
}
