use serde::{Serialize, Deserialize};

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "PascalCase")]
pub struct Data {
	#[serde(rename = "Hoge")]
	pub hoge: Hoge,
	#[serde(rename = "PkgHoge")]
	pub pkg_hoge: PkgHoge,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "PascalCase")]
pub struct Hoge {
	#[serde(rename = "Data")]
	pub data: u128,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "PascalCase")]
pub struct PkgHoge {
	#[serde(rename = "Data")]
	pub data: u128,
}

