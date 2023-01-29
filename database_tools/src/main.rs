mod database_creation;
mod database_insertion;

fn main() {
    println!("Hello, world!");
    database_creation::database_initialization::initialize();
    database_insertion::database_insertion::read_file();
}