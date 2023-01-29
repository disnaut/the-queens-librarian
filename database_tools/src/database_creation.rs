pub mod database_initialization { // This is a module that is public, which means that other things outside of file can access
    use std::{path::Path, fs::File};
    
    const FILE_NAME: &str = "database.sqlite3";

    pub fn initialize() {
        let exists = check_for_database(FILE_NAME); //Checks if the file exists.
    
        if exists == false {
            println!("The library is missing! Time to fix that!");
            create_database();
        }
        else {
            println!("The library is here! Ho0ray!");
        }
    }
    
    fn check_for_database(path: &str) -> bool { // functions are private by default
        return Path::new(&path).exists();
    }
    
    fn create_database() {      
        let path = Path::new(FILE_NAME); //creates file at root of project
    
        let display = path.display();
    
        match File::create(&path) {
            Err(why) => panic!("Couldn't create {}: {}", display, why), //Need to specify the error and okay message to do things.
            Ok(file) => file
        };
    }
}
