use serde::{Serialize, Deserialize};
use std::collections::HashMap;

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "PascalCase")]
pub struct Recursive {
	#[serde(rename = "Children")]
	pub children: Vec<Recursive>,
	#[serde(rename = "Re")]
	pub re: Option<Box<Recursive>>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "PascalCase")]
pub struct RecursiveMap {
	#[serde(rename = "Map")]
	pub map: HashMap<String, RecursiveMap>,
}

