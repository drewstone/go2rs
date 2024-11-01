use serde::{Serialize, Deserialize};

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "PascalCase")]
pub struct Recursive {
	#[serde(rename = "Re")]
	pub re: Option<Recursive>,
}

