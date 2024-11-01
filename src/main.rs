use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "PascalCase")]
pub struct Data {
    #[serde(rename = "Hoge")]
    pub hoge: Hoge,
    #[serde(rename = "PkgHoge")]
    pub pkg_hoge: Hoge,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "PascalCase")]
pub struct Hoge {
    #[serde(rename = "Data")]
    pub data: u128,
}

fn main() {
    println!("Hello, world!");
}
