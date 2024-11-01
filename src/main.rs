use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

pub mod example_01 {
    use super::*;

    #[derive(Debug, Clone, Copy, PartialEq, Serialize, Deserialize)]
    #[serde(rename_all = "lowercase")]
    pub enum EnumArrayValues {
        A,
        B,
        C,
    }

    #[derive(Debug, Clone, Copy, PartialEq, Serialize, Deserialize)]
    #[serde(rename_all = "PascalCase")]
    pub enum Status {
        Failure,
        OK,
    }

    #[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
    #[serde(rename_all = "PascalCase")]
    pub struct Data {
        #[serde(rename = "A")]
        pub a: u128,
        #[serde(rename = "Array")]
        pub array: Option<Vec<u128>>,
        #[serde(rename = "C")]
        pub c: String,
        #[serde(rename = "D")]
        pub d: Option<u128>,
        #[serde(rename = "EnumArray")]
        pub enum_array: Vec<EnumArrayValues>,
        #[serde(skip_serializing_if = "Option::is_none")]
        pub Foo: Option<Foo>,
        #[serde(rename = "Map")]
        pub map: HashMap<String, Status>,
        #[serde(rename = "OptionalArray")]
        pub optional_array: Vec<Option<String>>,
        #[serde(rename = "Package")]
        pub package: Option<Package>,
        #[serde(rename = "Status")]
        pub status: Status,
        #[serde(rename = "Time")]
        pub time: DateTime<Utc>,
        #[serde(rename = "U")]
        pub u: U,
        #[serde(skip_serializing_if = "Option::is_none")]
        pub b: Option<u128>,
        #[serde(skip_serializing_if = "Option::is_none")]
        pub foo: Option<u128>,
    }

    #[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
    #[serde(rename_all = "PascalCase")]
    pub struct Embedded {
        #[serde(skip_serializing_if = "Option::is_none")]
        pub foo: Option<u128>,
    }

    #[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
    #[serde(rename_all = "PascalCase")]
    pub struct Foo {
        #[serde(rename = "V")]
        pub v: u128,
    }

    #[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
    #[serde(rename_all = "PascalCase")]
    pub struct Package {
        pub data: u128,
    }

    #[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
    #[serde(rename_all = "PascalCase")]
    pub struct U {
        #[serde(rename = "Data")]
        pub data: u128,
    }
}

pub mod example_02 {
    use super::*;

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
}

pub mod example_03 {
    use super::*;

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
}

pub mod example_04 {
    use super::*;

    #[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
    #[serde(rename_all = "PascalCase")]
    pub struct Data {}

    #[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
    #[serde(rename_all = "PascalCase")]
    pub struct TestdataData {}

    #[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
    #[serde(rename_all = "PascalCase")]
    pub struct TestdataDataE5e4 {}
}

pub mod example_05 {
    use super::*;

    #[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
    #[serde(rename_all = "PascalCase")]
    pub struct CustomTest {
        #[serde(rename = "C")]
        pub c: CustomTestC,
    }

    #[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
    #[serde(rename_all = "PascalCase")]
    pub struct CustomTestC {}
}

fn main() {
    println!("Hello, world!");
}
