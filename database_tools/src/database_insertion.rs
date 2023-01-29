pub mod database_insertion {
    use serde::{Deserialize,Serialize};
    use serde_json::Value;
    use std::{path::Path,fs::File, default::Default};

    
    pub type Root = Vec<Card>;
    
    #[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
    #[serde(rename_all = "camelCase")]
    pub struct Card {
        pub object: String,
        pub id: String,
        #[serde(rename = "oracle_id")]
        pub oracle_id: String,
        #[serde(rename = "multiverse_ids")]
        pub multiverse_ids: Vec<i64>,
        #[serde(rename = "mtgo_id")]
        pub mtgo_id: Option<i64>,
        #[serde(rename = "mtgo_foil_id")]
        pub mtgo_foil_id: Option<i64>,
        #[serde(rename = "tcgplayer_id")]
        pub tcgplayer_id: Option<i64>,
        #[serde(rename = "cardmarket_id")]
        pub cardmarket_id: Option<i64>,
        pub name: String,
        pub lang: String,
        #[serde(rename = "released_at")]
        pub released_at: String,
        pub uri: String,
        #[serde(rename = "scryfall_uri")]
        pub scryfall_uri: String,
        pub layout: String,
        #[serde(rename = "highres_image")]
        pub highres_image: bool,
        #[serde(rename = "image_status")]
        pub image_status: String,
        #[serde(rename = "image_uris")]
        pub image_uris: Option<ImageUris>,
        #[serde(rename = "mana_cost")]
        pub mana_cost: Option<String>,
        pub cmc: f64,
        #[serde(rename = "type_line")]
        pub type_line: String,
        #[serde(rename = "oracle_text")]
        pub oracle_text: Option<String>,
        pub colors: Option<Vec<Value>>,
        #[serde(rename = "color_identity")]
        pub color_identity: Vec<Value>,
        pub keywords: Vec<Value>,
        pub legalities: Legalities,
        pub games: Vec<String>,
        pub reserved: bool,
        pub foil: bool,
        pub nonfoil: bool,
        pub finishes: Vec<String>,
        pub oversized: bool,
        pub promo: bool,
        pub reprint: bool,
        pub variation: bool,
        #[serde(rename = "set_id")]
        pub set_id: String,
        pub set: String,
        #[serde(rename = "set_name")]
        pub set_name: String,
        #[serde(rename = "set_type")]
        pub set_type: String,
        #[serde(rename = "set_uri")]
        pub set_uri: String,
        #[serde(rename = "set_search_uri")]
        pub set_search_uri: String,
        #[serde(rename = "scryfall_set_uri")]
        pub scryfall_set_uri: String,
        #[serde(rename = "rulings_uri")]
        pub rulings_uri: String,
        #[serde(rename = "prints_search_uri")]
        pub prints_search_uri: String,
        #[serde(rename = "collector_number")]
        pub collector_number: String,
        pub digital: bool,
        pub rarity: String,
        #[serde(rename = "flavor_text")]
        pub flavor_text: Option<String>,
        #[serde(rename = "card_back_id")]
        pub card_back_id: Option<String>,
        pub artist: String,
        #[serde(rename = "artist_ids")]
        pub artist_ids: Option<Vec<String>>,
        #[serde(rename = "illustration_id")]
        pub illustration_id: Option<String>,
        #[serde(rename = "border_color")]
        pub border_color: String,
        pub frame: String,
        #[serde(rename = "full_art")]
        pub full_art: bool,
        pub textless: bool,
        pub booster: bool,
        #[serde(rename = "story_spotlight")]
        pub story_spotlight: bool,
        #[serde(rename = "edhrec_rank")]
        pub edhrec_rank: Option<i64>,
        pub prices: Prices,
        #[serde(rename = "related_uris")]
        pub related_uris: RelatedUris,
    }
    
    #[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
    #[serde(rename_all = "camelCase")]
    pub struct ImageUris {
        pub small: String,
        pub normal: String,
        pub large: String,
        pub png: String,
        #[serde(rename = "art_crop")]
        pub art_crop: String,
        #[serde(rename = "border_crop")]
        pub border_crop: String,
    }
    
    #[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
    #[serde(rename_all = "camelCase")]
    pub struct Legalities {
        pub standard: String,
        pub future: String,
        pub historic: String,
        pub gladiator: String,
        pub pioneer: String,
        pub explorer: String,
        pub modern: String,
        pub legacy: String,
        pub pauper: String,
        pub vintage: String,
        pub penny: String,
        pub commander: String,
        pub brawl: String,
        pub historicbrawl: String,
        pub alchemy: String,
        pub paupercommander: String,
        pub duel: String,
        pub oldschool: String,
        pub premodern: String,
    }
    
    #[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
    #[serde(rename_all = "camelCase")]
    pub struct Prices {
        pub usd: Option<Value>,
        #[serde(rename = "usd_foil")]
        pub usd_foil: Option<Value>,
        #[serde(rename = "usd_etched")]
        pub usd_etched: Option<Value>,
        pub eur: Option<Value>,
        #[serde(rename = "eur_foil")]
        pub eur_foil: Option<Value>,
        pub tix: Option<Value>,
    }
    
    #[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
    #[serde(rename_all = "camelCase")]
    pub struct RelatedUris {
        //pub gatherer: String,
        #[serde(rename = "tcgplayer_infinite_articles")]
        pub tcgplayer_infinite_articles: Option<String>,
        #[serde(rename = "tcgplayer_infinite_decks")]
        pub tcgplayer_infinite_decks: Option<String>,
        pub edhrec: Option<String>,
    }

    const REAL_FILE: &str = "oracle-cards-20230119220451.json";

    pub fn read_file() {
        let json_file_path = Path::new(REAL_FILE);
        let file = File::open(json_file_path).expect("Issue with file reading.");
        // as responses go down, we need to have something to throw in case of an operational error.
        // in this example, we have to deal with IO, which means something can go wrong.
        // so we need to have an error message in order for this to fully work.

        let root: Root = serde_json::from_reader(file).expect("There was an issue here.");
    }
}